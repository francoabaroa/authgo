package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"

	"github.com/francoabaroa/authgo/pkg/handlers"
)

type HelloResponse struct {
	Hello string `json:"hello"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
    })
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

	http.Handle("/login", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
				handlers.LoginPageHandler(w, r)
		} else if r.Method == http.MethodPost {
				handlers.LoginHandler(w, r)
		}
	})))

	http.Handle("/register", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
				handlers.RegisterPageHandler(w, r)
		} else if r.Method == http.MethodPost {
				handlers.RegisterHandler(w, r)
		}
	})))

	http.Handle("/reset_password", loggingMiddleware(http.HandlerFunc(handlers.ResetPasswordHandler)))
	http.Handle("/show_users", loggingMiddleware(http.HandlerFunc(handlers.ShowUsersHandler)))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
