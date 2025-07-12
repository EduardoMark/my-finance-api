package httputils

import "net/http"

func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func Error(w http.ResponseWriter, r *http.Request, status int, message any) {
	EncodeJson(w, r, status, message)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
