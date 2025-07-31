package handlers

import "net/http"

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
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
