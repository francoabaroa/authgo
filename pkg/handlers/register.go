package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"time"
)

type User struct {
	Username string
	Email    string
	Password string
}

var mockDB = make(map[string]User)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Check if the username or email already exists in the mockDB
    for _, user := range mockDB {
        if user.Username == username || user.Email == email {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(RegisterResponse{
                Message: "Username or email already exists",
            })
            return
        }
    }

    // Hash the password
    hash := sha256.Sum256([]byte(password))
    hashedPassword := fmt.Sprintf("%x", hash)

    // Insert into the mockDB
    mockDB[username] = User{
        Username: username,
        Email:    email,
        Password: hashedPassword,
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

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
    // Check for a "user" cookie
    _, err := r.Cookie("user")

    // If the cookie exists (no error), the user is already logged in
    if err == nil {
        // Redirect the user to the main page
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // The user is not logged in, serve them the registration page
    tmpl, err := template.ParseFiles("templates/register.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Render the template
    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ShowUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User

	// Add each user from the mockDB
	for _, user := range mockDB {
		users = append(users, user)
	}

	// Create response with users
	response := UsersResponse{
		Users: users,
	}

	// Convert to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
