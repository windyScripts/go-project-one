package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func execsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /execs/", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)
	mux.HandleFunc("POST /execs/updatepassword", handlers.AddExecHandler)

	// mux.HandleFunc("POST /execs/login", handlers.AddExecHandler)
	// mux.HandleFunc("POST /execs/logout", handlers.AddExecHandler)
	// mux.HandleFunc("POST /execs/forgotpassword", handlers.AddExecHandler)
	// mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.AddExecHandler)

	return mux

}
