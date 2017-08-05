package snapi

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type Response struct {
	Type    string
	Message string
	Data    interface{}
}

func JSONResponseHandler(w http.ResponseWriter, returnval interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnval)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Type:"error", Message:fmt.Sprintf("Route %s not found, " +
		"please check request and try again", r.URL.Path)})
}

func resourceNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Type:"error", Message:fmt.Sprintf("Resource %s not found, " +
		"please check request and try again", r.URL.Path)})
}
