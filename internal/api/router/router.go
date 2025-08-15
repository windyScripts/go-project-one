package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	eRouter := execsRouter()
	tRouter := teachersRouter()
	sRouter := studentRouter()

	sRouter.Handle("/", eRouter)
	tRouter.Handle("/", sRouter)
	return tRouter

	/* mux := http.NewServeMux()
	// single space between method and endpoint.
	mux.HandleFunc("GET /", handlers.RootHandler)
	mux.HandleFunc("GET /execs/", handlers.ExecsHandler) */
	//return mux

}
