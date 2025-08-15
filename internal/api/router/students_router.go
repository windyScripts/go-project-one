package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func studentRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddStudentHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentHandler)

	return mux
}
