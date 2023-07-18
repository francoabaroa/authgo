package main

import (
	"log"
	"net/http"
	"time"

	"github.com/francoabaroa/authgo/pkg/handlers"
)

const (
	GET  = http.MethodGet
	POST = http.MethodPost
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
    })
}

func main() {
	http.Handle("/", loggingMiddleware(http.HandlerFunc(handlers.WelcomePageHandler)))

	http.Handle("/login", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == GET {
			handlers.LoginPageHandler(w, r)
		} else if r.Method == POST {
			handlers.LoginHandler(w, r)
		}
	})))

	http.Handle("/register", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == GET {
			handlers.RegisterPageHandler(w, r)
		} else if r.Method == POST {
			handlers.RegisterHandler(w, r)
		}
	})))

	http.Handle("/reset_password", loggingMiddleware(http.HandlerFunc(handlers.ResetPasswordHandler)))
	http.Handle("/logout", loggingMiddleware(http.HandlerFunc(handlers.LogoutHandler)))
	http.Handle("/show_users", loggingMiddleware(http.HandlerFunc(handlers.ShowUsersHandler)))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
