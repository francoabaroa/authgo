package handlers

import (
	"net/http"
)

func WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
    // Try to get the user cookie
    cookie, err := r.Cookie("user")
    var data interface{} // Data to pass to the template
    if err == nil {
        // If the cookie is found, pass the username to the template
        data = map[string]string{
            "Username": cookie.Value,
        }
    }

    if err := t.ExecuteTemplate(w, "welcome.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
