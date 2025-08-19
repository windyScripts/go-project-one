package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

/* type user struct {
	// field names must be public so they can be accessed.
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
} */

//go:embed .env
var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
	// create temp file for env
	tempfile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating .env file: %v", err)
	}

	defer os.Remove(tempfile.Name())

	// Write content to temp file
	_, err = tempfile.Write(content)
	if err != nil {
		log.Fatalf("Error writing to temp .env file: %v", err)
	}
	err = tempfile.Close()
	if err != nil {
		log.Fatalf("Error closing temp .env file: %v", err)
	}

	// load env vars from temp file:
	err = godotenv.Load(tempfile.Name())
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

}

func main() {
	// only in development for running source code
	// err := godotenv.Load()
	// if err != nil {
	// 	return
	// }

	// load env vars from embedded .env file
	loadEnvFromEmbeddedFile()

	fmt.Println("Env variable cert_file:", os.Getenv("CERT_FILE"))

	port := os.Getenv("API_PORT")

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	router := router.MainRouter()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	//secureMux := mw.Cors(rl.Middleware(mw.ReponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)((mux)))))))
	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/execs/login", "/execs/forgotpassword", "/execs/resetpassword/reset")

	// secureMux := jwtMiddleware(mw.SecurityHeaders(router)) // sidestepping middlewares for testing
	// secureMux := mw.XSSMiddleware(router)
	secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ReponseTimeMiddleware, rl.Middleware, mw.Cors)
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

/* 
go build -o binaries/rest_api cmd/api/server.go
Binary file only runs on same os and cpu architecture as the source computer

GOOS=darwin GOARCH=arm64 go build -o rest_api_macos_arm64 cmd/api/server.go
save binaries in binary folders
./rest_api
for benchmarking, remove rate limiter.

garble, gobfuscate for obfuscating code. mvdan.cc/garble github.com/unixpickle/gobfuscate

garble build -o file_name server.go_location

*/