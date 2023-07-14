package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest

	// Decode the JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the username already exists in the mockDB
	if _, exists := mockDB[request.Username]; exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RegisterResponse{
			Message: "Username or email already exists",
		})
		return
	}

	// Hash the password
	hash := sha256.Sum256([]byte(request.Password))
	hashedPassword := fmt.Sprintf("%x", hash)

	// Insert into the mockDB
	mockDB[request.Username] = User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		Message: "Registration successful",
	})
}
