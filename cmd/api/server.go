package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

/* type user struct {
	// field names must be public so they can be accessed.
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
} */

func main() {

	err := godotenv.Load()
	if err != nil {
		return
	}


	_, err2 := sqlconnect.ConnectDb()
	if err2 != nil {
		fmt.Println("Error connecting DB: ", err2)
		return
	}

	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	router := router.Router()

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
	//secureMux := utils.ApplyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ReponseTimeMiddleware, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(router) // sidestepping middlewares for testing

	// create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)
	err1 := server.ListenAndServeTLS(cert, key)
	if err1 != nil {
		log.Fatalln("Error starting server", err1)
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

/*
query parameters are depolluted automatically, only the first key value pair is stored.
in body parameters, cleaning is not done automatically.
hpp middleware handles this situation. It normalizes by removing duplicates, reducing ambiguity.

*/

/*
Mariadb is used.
insert into table (v1, v2) values("v1","v2"),("v3","v4")
update table set value = condition where check = condition
delete from my_table where condition
rename table my_table my_new_table

renaming database in maria and mysql doesn't work by using rename database.
done by creating db, copying tables, delete old database.

drop table my_table
drop database test_database
*/
