package handlers

import (
	"html/template"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl := template.Must(template.ParseFiles("templates/login.html"))
        tmpl.Execute(w, nil)
        return
    }

    // Parse form
    r.ParseForm()
    username := r.FormValue("username")
    password := r.FormValue("password")

    // TODO: Check credentials against database
    // For now, just demo:
    if username == "demo" && password == "password" {
        token, _ := auth.GenerateToken(1, username)
        
        // Set JWT as cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "token",
            Value:    token,
            Path:     "/",
            HttpOnly: true,
        })

        // Redirect to chat
        w.Header().Set("HX-Redirect", "/chat")
        return
    }

    // Return error message that HTMX will insert
    w.Write([]byte("<div class='error'>Invalid credentials</div>"))
} 