package handlers

import (
	"database/sql"
	//"encoding/json"
	"html/template"
	"net/http"
	"time"

	"real-time-forum/backend/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"` // nickname or email
	Password   string `json:"password"`
}

// Handles client login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("./frontend/template/login.html")
		if err != nil {
			// ErrorHandler(w, r, "Failed to load login page", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	creds := LoginRequest{
		Identifier: r.FormValue("identifier"),
		Password:   r.FormValue("password"),
	}
	var hashedPassword string
	var userID string
	err = database.Db.QueryRow(`SELECT id,password FROM users WHERE nickname = ? OR email = ?`, creds.Identifier, creds.Identifier).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}
	// Create a session
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionID := u.String()

	expiresAt := time.Now().Add(24 * time.Hour) // Sessions expire after 24 hours

	_, err = database.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
		userID, sessionID, expiresAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured while logging you in"))
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
