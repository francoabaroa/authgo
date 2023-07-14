package main

import (
	"encoding/json"
	"net/http"
	"github.com/francoabaroa/authgo/pkg/handlers"
)

type HelloResponse struct {
	Hello string `json:"hello"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		response := SuccessResponse{
			Success: "ful",
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/reset_password", handlers.ResetPasswordHandler)

	http.ListenAndServe(":8080", nil)
}
