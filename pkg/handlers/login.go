package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
	"os"
	"html/template"

	"github.com/joho/godotenv"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		fmt.Println("No secret key set in .env file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, ok := mockDB[username]
	if !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hash)

	if user.Password != hashedPassword {
		http.Error(w, "Invalid password", http.StatusBadRequest)
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
	_, err := r.Cookie("user")
    if err == nil {
        // If there is no error, a cookie was found -> user is logged in
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
