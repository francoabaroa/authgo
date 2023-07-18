package handlers

import (
	"net/http"
	"time"
)

// LogoutHandler handles the logout process
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Set the cookie's expiration time to a past time (zero value)
	http.SetCookie(w, &http.Cookie{
		Name:    "user",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		Path:    "/",
	})

	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
