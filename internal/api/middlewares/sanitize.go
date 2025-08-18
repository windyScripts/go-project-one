package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"restapi/pkg/utils"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	fmt.Println("Init XSSMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("XSSMiddleware Ran")

		// sanitize the url path
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// sanitize query params
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var sanitizedValues []string
			for _, value := range values {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}

		r.URL.Path = sanitizedPath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					utils.ErrorHandler(err, "")
					http.Error(w, utils.ErrorHandler(err, "Error reading request body").Error(), http.StatusBadRequest)
					return
				}
				bodyString := strings.TrimSpace(string(bodyBytes))
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))
				if len(bodyString) > 0 {
					// reset the request body
					var inputData interface{}
					err = json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, utils.ErrorHandler(err, "Error reading request body").Error(), http.StatusBadRequest)
						return
					}
					sanitizedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					// marshall sanitized data back to the body
					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						http.Error(w, utils.ErrorHandler(err, "Error sanitizing body").Error(), http.StatusBadRequest)
						return
					}
					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))

				} else {
					log.Println("Request body is empty.")
				}

			} else {
				log.Println("No body in the request")
			}
		} else if r.Header.Get("Content-Type") != "" {
			log.Printf("Received request with unsupported Content-Type: %s. Expected application/json. \n", r.Header.Get("Content-Type"))
			http.Error(w, "Unsupported Content-Type. Please use application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("XSSMiddleware ran")
	})
}

// clean sanitizes input data to prevent XSS attacks
func clean(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = sanitizeValue(val)
		}
		return v, nil
	case []interface{}:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		return nil, utils.ErrorHandler(fmt.Errorf("unsupported type %T", data), fmt.Sprintf("unsupported type %T", data))
	}
}

func sanitizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		for key, val := range v {
			v[key] = sanitizeValue(val)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v
	default:
		return v // return v as it is unsupported
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
