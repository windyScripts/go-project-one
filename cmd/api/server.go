package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
)

type user struct {
	// field names must be public so they can be accessed.
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Root Route")
		w.Write([]byte("Hello Root Route"))
		fmt.Println("Hello Root Route")
	}

func teachersHandler(w http.ResponseWriter, r *http.Request){
/* 	fmt.Println(r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	userID := strings.TrimSuffix(path, "/")

	// teachers/{id}
	// teachers/?key=value8query=value2&sortby=email&sortorder=ASC

	fmt.Println("ID IS: ",userID)

	fmt.Println("Query Params:", r.URL.Query())
	queryParams := r.URL.Query()
	
	key := queryParams.Get("key")
	sortby := queryParams.Get("sortby")
	sortorder := queryParams.Get("sortorder")

	if sortorder == "" {
		sortorder = "ASC"
	} */

	switch r.Method{
		case http.MethodGet:
			w.Write([]byte("Hello GET method on teachers Route"))
			fmt.Println("Hello GET method on teachers Route")
		case http.MethodPost:
			w.Write([]byte("Hello POST method on teachers Route"))
			fmt.Println("Hello POST method on teachers Route")
		case http.MethodPatch:
			w.Write([]byte("Hello PATCH method on teachers Route"))
			fmt.Println("Hello PATCH method on teachers Route")
		case http.MethodPut:
			w.Write([]byte("Hello PUT method on teachers Route"))
			fmt.Println("Hello PUT method on teachers Route")
		case http.MethodDelete:
			w.Write([]byte("Hello DELETE method on teachers Route"))
			fmt.Println("Hello DELETE method on teachers Route")
		default:
		w.Write([]byte("Hello teachers Route"))
		fmt.Println("Hello teachers Route")
		}
	}

	func studentsHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Students Route")
		switch r.Method{
		case http.MethodGet:
			w.Write([]byte("Hello GET method on Students Route"))
			fmt.Println("Hello GET method on Students Route")
		case http.MethodPost:
			w.Write([]byte("Hello POST method on Students Route"))
			fmt.Println("Hello POST method on Students Route")
		case http.MethodPatch:
			w.Write([]byte("Hello PATCH method on Students Route"))
			fmt.Println("Hello PATCH method on Students Route")
		case http.MethodPut:
			w.Write([]byte("Hello PUT method on Students Route"))
			fmt.Println("Hello PUT method on Students Route")
		case http.MethodDelete:
			w.Write([]byte("Hello DELETE method on Students Route"))
			fmt.Println("Hello DELETE method on Students Route")
		default:
		w.Write([]byte("Hello Students Route"))
		fmt.Println("Hello Students Route")
		}
	}

	func execsHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello execs Route")
		switch r.Method{
		case http.MethodGet:
			w.Write([]byte("Hello GET method on Execs Route"))
			fmt.Println("Hello GET method on Execs Route")
		case http.MethodPost:
			w.Write([]byte("Hello POST method on Execs Route"))
			fmt.Println("Hello POST method on Execs Route")
		case http.MethodPatch:
			w.Write([]byte("Hello PATCH method on Execs Route"))
			fmt.Println("Hello PATCH method on Execs Route")
		case http.MethodPut:
			w.Write([]byte("Hello PUT method on Execs Route"))
			fmt.Println("Hello PUT method on Execs Route")
		case http.MethodDelete:
			w.Write([]byte("Hello DELETE method on Execs Route"))
			fmt.Println("Hello DELETE method on Execs Route")
		default:
		w.Write([]byte("Hello teacExecshers Route"))
		fmt.Println("Hello Execs Route")
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

	// create custom server

	server := &http.Server{
		Addr:      port,
		Handler: mw.Compression(mw.ReponseTimeMiddleware((mw.SecurityHeaders(mw.Cors(mux))))),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ",port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting server", err)
	}
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

