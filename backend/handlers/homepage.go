package handlers

import (
	"html/template"
	"net/http"
)

type HomePageData struct {
	IsLoggedIn bool
	UserName   string
}
func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./frontend/template/index.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Example: simulate login check using a cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		// Not logged in
		tmpl.Execute(w, HomePageData{
			IsLoggedIn: false,
		})
		return
	}

	// TODO: Query your database to get user info by session token
	// For now, we assume it's valid for testing
	tmpl.Execute(w, HomePageData{
		IsLoggedIn: true,
		UserName:   "JohnDoe", // Replace with actual user from DB
	})
}
