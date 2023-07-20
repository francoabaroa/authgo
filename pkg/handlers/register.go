package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Email    string
	Password string
}

type UsersResponse struct {
	Users []User `json:"users"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate that username, email and password are not empty
	if username == "" || email == "" || password == "" {
		// Render the error page
		if err := t.ExecuteTemplate(w, "error.html", "Username, email or password cannot be empty"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if the username or email already exists in the DB
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = $1 OR email = $2", username, email).Scan(&userID)
	if err != sql.ErrNoRows {
		// Render the error page
		if err := t.ExecuteTemplate(w, "error.html", "Username or email already exists"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into the DB
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the user cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "user",
		Value:   username,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // Cookie expires after 24 hours
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterPageHandler serves the registration page
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check for a "user" cookie
	_, err := r.Cookie("user")

	// If the cookie exists (no error), the user is already logged in
	if err == nil {
		// Redirect the user to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the template
	if err := t.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ShowUsersHandler returns a list of users
func ShowUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Query the DB for users
	rows, err := db.Query("SELECT username, email, password FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		// Scan the retrieved row into the User struct
		err = rows.Scan(&user.Username, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Append the User struct to the users slice
		users = append(users, user)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
