package handlers

import (
	"net/http"
)

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement password reset logic
	w.Write([]byte("Reset Password handler not implemented"))
}