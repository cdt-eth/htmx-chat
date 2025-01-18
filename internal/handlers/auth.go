package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/auth"
	"github.com/cdt-eth/htmx-chat/internal/models"
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

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    user, err := models.CreateUser(username, password)
    if err != nil {
        w.Write([]byte(fmt.Sprintf(`<div class="error">%s</div>`, err.Error())))
        return
    }

    // Generate JWT
    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        w.Write([]byte(`<div class="error">Error creating account</div>`))
        return
    }

    // Set cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
    })

    // Return success and redirect
    w.Header().Set("HX-Redirect", "/")
} 