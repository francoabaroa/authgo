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

    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
