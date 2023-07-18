package handlers

import (
	"net/http"
	"html/template"
)

func WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/welcome.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Try to get the user cookie
    cookie, err := r.Cookie("user")
    var data interface{} // Data to pass to the template
    if err == nil {
        // If the cookie is found, pass the username to the template
        data = map[string]string{
            "Username": cookie.Value,
        }
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}