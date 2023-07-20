package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"
	"golang.org/x/crypto/bcrypt"

	"github.com/joho/godotenv"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if os.Getenv("FLY_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			http.Error(w, "Warning: No .env file found. Falling back to system environment variables.", http.StatusInternalServerError)
			return
		}
	}

	// Get secret key from .env file
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get hashed password from the database
	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		// Render the error page
		if err := t.ExecuteTemplate(w, "error.html", "Username not found"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare the hashed password with the password provided by the user
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Render the error page
		if err := t.ExecuteTemplate(w, "error.html", "Incorrect password"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set the user cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "user",
		Value:   username,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // Cookie expires after 24 hours
	})

	// Redirect the user to the root route
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check for a "user" cookie
	_, err := r.Cookie("user")
	if err == nil {
		// If there is no error, a cookie was found -> user is logged in
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the login page
	if err := t.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
