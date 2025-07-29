package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ReponseTimeMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("Received Request in ResponseTime.")
		start:= time.Now()

		// Create a custom ResponseWriter to capture the status code
		wrappedWriter := &responseWriter{ResponseWriter : w, status:http.StatusOK }

		
		
		// Calculate the duration
		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter,r)
		//log request details
		fmt.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Sent Response from Response time middleware.")
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int){
	rw.status = code
	// passes on the code to the native writeHeader method.
	rw.ResponseWriter.WriteHeader(code)

}