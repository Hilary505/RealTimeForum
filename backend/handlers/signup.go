package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"


	"real-time-forum/backend/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Nickname        string `json:"nickname"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("./frontend/template/signup.html")
		if err != nil {
			log.Println("Error loading signup page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	req := SignupRequest{
		Nickname:        r.FormValue("nickname"),
		FirstName:       r.FormValue("first_name"),
		LastName:        r.FormValue("last_name"),
		Age:             parseIntOrDefault(r.FormValue("age"), 0),
		Gender:          r.FormValue("gender"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	if req.Nickname == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Required fields are missing", http.StatusBadRequest)
		return
	}

	// Check if nickname or email already exists
	var exists int
	err = database.Db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE nickname = ? OR email = ?
	`, req.Nickname, req.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "Nickname or email already in use", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Generate user ID
	userID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Error generating user ID", http.StatusInternalServerError)
		return
	}

	// Insert new user
	_, err = database.Db.Exec(`
		INSERT INTO users (uuid, nickname, firstname, lastname, age, gender, email, password)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, userID.String(), req.Nickname, req.FirstName, req.LastName, req.Age, strings.ToLower(req.Gender), req.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func parseIntOrDefault(value string, def int) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}
