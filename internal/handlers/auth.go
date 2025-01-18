package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/cdt-eth/htmx-chat/internal/auth"
	"github.com/cdt-eth/htmx-chat/internal/models"
	"golang.org/x/crypto/bcrypt"
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

    // Get user from database
    user, err := models.GetUserByUsername(username)
    if err != nil {
        w.Write([]byte("<div class='error'>Invalid credentials</div>"))
        return
    }

    // Check password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        w.Write([]byte("<div class='error'>Invalid credentials</div>"))
        return
    }

    // Generate JWT
    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        w.Write([]byte("<div class='error'>Login error</div>"))
        return
    }

    // Set JWT cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
    })

    // Redirect to chat
    w.Header().Set("HX-Redirect", "/")
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

    // Show success message and redirect after 1 second
    w.Header().Set("HX-Trigger", "showMessage")
    w.Write([]byte(`
        <div class="success">
            Account created! Redirecting...
            <script>
                setTimeout(() => {
                    window.location.href = "/";
                }, 1000);
            </script>
        </div>
    `))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Clear the cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
    })

    // Both HTMX and regular redirect
    if r.Header.Get("HX-Request") == "true" {
        w.Header().Set("HX-Redirect", "/auth/login")
    } else {
        http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
    }
} 