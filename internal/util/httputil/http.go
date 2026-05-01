package httputil

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func BindJSON(r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func Chain(handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
