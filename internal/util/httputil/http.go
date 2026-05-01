package httputil

import (
	"bytes"
	"encoding/json"
	"io"
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

func GetJSONKeys(r *http.Request) ([]string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(jsonData))
	for key := range jsonData {
		keys = append(keys, key)
	}

	return keys, nil
}
