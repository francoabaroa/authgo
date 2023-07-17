package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"

	"github.com/golang-jwt/jwt"
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
	var request LoginRequest

	err := godotenv.Load(".env")
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

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := mockDB[request.Username]
	if !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(request.Password))
	hashedPassword := fmt.Sprintf("%x", hash)

	if user.Password != hashedPassword {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &jwt.StandardClaims{
		Subject:   request.Username,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey)) // TODO: Use a more secure key in production

	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// TODO: store JWTs in HTTP-only cookies in production
	json.NewEncoder(w).Encode(LoginResponse{
		Token: tokenString,
	})
}
