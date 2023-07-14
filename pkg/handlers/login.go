package handlers

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement login logic
	w.Write([]byte("Login handler not implemented"))
}