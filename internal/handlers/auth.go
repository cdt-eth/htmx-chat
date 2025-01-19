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

func SignupFormHandler(w http.ResponseWriter, r *http.Request) {
    html := `
    <div>
        <form id="signup-form" hx-post="/auth/signup" hx-target="#signup-message" hx-swap="innerHTML">
            <h2>Create Account</h2>
            <div class="form-group">
                <label for="signup-username">Username</label>
                <input type="text" 
                       name="username" 
                       id="signup-username" 
                       required
                       pattern="[a-zA-Z0-9_]+"
                       title="Only letters, numbers, and underscores allowed"
                       minlength="3"
                       maxlength="30" />
            </div>
            <div class="form-group">
                <label for="signup-password">Password</label>
                <input type="password" 
                       name="password" 
                       id="signup-password"
                       required
                       minlength="8"
                       pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}"
                       title="Must contain at least one number, one uppercase and lowercase letter, and be at least 8 characters long" />
            </div>
            <div id="signup-message"></div>
            <button type="submit">Sign Up</button>
        </form>

        <a href="#" 
           class="auth-link"
           hx-get="/auth/login-form"
           hx-target="#auth-form"
           hx-swap="innerHTML">
            Already have an account? Login
        </a>
    </>`
    
    w.Write([]byte(html))
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
    html := `
    <div>
        <form id="login-form" hx-post="/auth/login" hx-target="#login-message" hx-swap="innerHTML">
            <h2>Login</h2>
            <div class="form-group">
                <label for="login-username">Username</label>
                <input type="text" name="username" id="login-username" required />
            </div>
            <div class="form-group">
                <label for="login-password">Password</label>
                <input type="password" name="password" id="login-password" required />
            </div>
            <div id="login-message"></div>
            <button type="submit">Login</button>
        </form>

        <a href="#" 
           class="auth-link"
           hx-get="/auth/signup-form"
           hx-target="#auth-form"
           hx-swap="innerHTML">
            Need an account? Sign up
        </a>
    </div>`
    
    w.Write([]byte(html))
} 