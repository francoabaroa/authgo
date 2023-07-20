package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/francoabaroa/authgo/pkg/handlers"
	"github.com/joho/godotenv"
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
	var psqlInfo string

	if os.Getenv("FLY_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: No .env file found. Falling back to system environment variables.")
		}

		host := os.Getenv("DB_HOST")
		port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // convert port to int
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	} else {
		psqlInfo = os.Getenv("DATABASE_URL")
	}

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Provide db connection to handlers
	handlers.Init(db)

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
