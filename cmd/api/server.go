package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"strconv"
	"strings"
)

/* type user struct {
	// field names must be public so they can be accessed.
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
} */

type Teacher struct {
	ID        int
	FirstName string
	LastName  string
	Class     string
	Subject   string
}

var (
	teachers = make(map[int]Teacher)
	// mutex = &sync.Mutex{}
	nextID = 1
)

// Initialize some dummy data

func init() {
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "11C",
		Subject:   "English",
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Root Route"))
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		teacherList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}

		}
		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	// Handle path param
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)

}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		w.Write([]byte("Hello POST method on teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH method on teachers Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on teachers Route"))
	default:
		w.Write([]byte("Hello teachers Route"))
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello Students Route")
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on Students Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST method on Students Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH method on Students Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on Students Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on Students Route"))
	default:
		w.Write([]byte("Hello Students Route"))
		fmt.Println("Hello Students Route")
	}
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello execs Route")
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on Execs Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST method on Execs Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH method on Execs Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on Execs Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on Execs Route"))
	default:
		w.Write([]byte("Hello Execs Route"))
	}
}

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery: true,
	// 	CheckBody: true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist: []string{"sortBy","sortOrder","name","age","class"},
	// }

	//secureMux := mw.Cors(rl.Middleware(mw.ReponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)((mux)))))))
	//secureMux := applyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ReponseTimeMiddleware, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(mux) // sidestepping middlewares for testing

	// create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting server", err)
	}
}

// Middleware is a function that wraps an http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

/*
mux refers to request multiplexer, which is used to route requests to the appropriate handler based on the URL path and HTTP method.
used to organize api better
separate logic for different routes.
use http.HandleFunc when low number of routes. This implicitly uses mux, but doesn't require explicit syntax
*/

/*
middleware for logging, auth, data validation, error handling.
*/

/*
Wrapping multiple middleware functions insde one another is called chaininng handlers
*/

/*
query parameters are depolluted automatically, only the first key value pair is stored.
in body parameters, cleaning is not done automatically.
hpp middleware handles this situation. It normalizes by removing duplicates, reducing ambiguity.

*/
