package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
)

// ==============Validate email format==========
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ===== Check if a username or email (credential) already exists in a database db =====
func CredentialExists(db *sql.DB, credential string) bool {
	query := `SELECT COUNT(*) FROM users WHERE nickname = ? OR email = ?`
	var count int
	err := db.QueryRow(query, credential, credential).Scan(&count)
	if err != nil {
		fmt.Printf("|credential exist| ---> {%v}", err)
		return false
	}
	return count > 0
}

/*
=== ValidateSession checks if a session token is valid. The function takes a pointer to the request
and returns a boolean value and a user_ID of type string based on the session_token found in the
cookie present in the header, within the request =====
*/
func ValidateSession(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Printf("|validate session| ---> no session cookie found")
		return false, ""
	}

	var (
		userID    string
		expiresAt time.Time
	)

	err = database.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		fmt.Printf("|validate session| ---> {%v}", err)
		return false, ""
	}

	if time.Now().After(expiresAt) {
		fmt.Printf("session expired for user %s", userID)
		return false, ""
	}

	fmt.Printf("[SUCCESS]: Session valid for user: %s", userID)
	return true, userID
}

// getUserFromSession gets a user from a session ID
func GetUserFromSession(sessionID string) (*models.User, error) {
	var userID int
	var expiresAt time.Time

	// Get the user ID from the session
	err := database.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return nil, err
	}

	// Check if the session is expired
	if time.Now().After(expiresAt) {
		database.Db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
		return nil, sql.ErrNoRows
	}

	// Get the user information
	var user models.User
	err = database.Db.QueryRow(
		"SELECT id, nickname, age, gender, first_name, last_name, email, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
