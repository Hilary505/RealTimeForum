package handlers

import (
	"html/template"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

// HandleHomepage serves the index.html page for the forum platform
func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get user session info
	valid, userID := utils.ValidateSession(r)
	
	// Parse and serve the template
	tmpl, err := template.ParseFiles("/home/docker/real-time-forum/frontend/template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create template data with login state
	data := struct {
		IsLoggedIn bool
		Nickname   string
	}{
		IsLoggedIn: valid,
		Nickname:   "",
	}

	// If user is logged in, fetch their nickname
	if valid {
		var nickname string
		err := database.Db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
		if err == nil {
			data.Nickname = nickname
		}
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}