package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/logger"
)

var sessionManager = auth.NewSessionManager()

// HandleLogin handles login API requests
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// Validate credentials
	if creds.Username != config.AppConfig.HTTPAuth.Username ||
		creds.Password != config.AppConfig.HTTPAuth.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Unauthorized",
			"message": "用户名或密码错误",
		})
		logger.Warn("Failed login attempt for user: %s from %s", creds.Username, r.RemoteAddr)
		return
	}

	// Create session
	session, err := sessionManager.CreateSession(creds.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create session",
		})
		logger.Error("Failed to create session: %v", err)
		return
	}

	// Set cookie
	cookieMaxAge := 0
	if creds.Remember {
		cookieMaxAge = 86400 * 30 // 30 days if remember me is checked
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Should be true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAge,
	})

	// Set username cookie if remember me is checked
	if creds.Remember {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_username",
			Value:    creds.Username,
			Path:     "/",
			HttpOnly: false,
			Secure:   false,
			MaxAge:   cookieMaxAge,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"token":    session.Token,
		"username": session.Username,
		"expires":  session.ExpiresAt.Unix(),
	})

	logger.Info("User logged in: %s from %s", creds.Username, r.RemoteAddr)
}

// HandleLogout handles logout requests
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		sessionManager.DeleteSession(cookie.Value)
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_username",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1,
	})

	// Check if it's an API request
	if r.Header.Get("Accept") == "application/json" || r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"success": true,
		})
		return
	}

	// Redirect to login page for web requests
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	logger.Info("User logged out")
}

// HandleLoginPage serves the login page
func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	// If authentication is not enabled, redirect to home
	if !config.AppConfig.HTTPAuth.Enabled {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if user is already logged in
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		if _, valid := sessionManager.ValidateSession(cookie.Value); valid {
			// Already logged in, redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Serve login page
	tmpl, err := template.ParseFS(templates, "templates/login.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("Failed to parse login template: %v", err)
		return
	}

	data := struct {
		Version    string
		ProjectURL string
	}{
		Version:    config.AppConfig.Version,
		ProjectURL: config.AppConfig.ProjectURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("Failed to execute login template: %v", err)
	}
}

// HandleCheckAuth checks if user is authenticated
func HandleCheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	session, valid := sessionManager.ValidateSession(cookie.Value)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"username":      session.Username,
		"expires":       session.ExpiresAt.Unix(),
	})
}

// GetSessionManager returns the session manager instance
func GetSessionManager() *auth.SessionManager {
	return sessionManager
}
