package routes

import (
	"net/http"

	"real-time-forum/backend/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	//  Serve static files
	fileServer := http.FileServer(http.Dir("frontend/css"))
	mux.Handle("/css/", http.StripPrefix("/css/", fileServer))
	scriptServer := http.FileServer(http.Dir("frontend/js/"))
	mux.Handle("/js/", http.StripPrefix("/js/", scriptServer))
	imageServer := http.FileServer(http.Dir("frontend/images"))
	mux.Handle("/static/images/", http.StripPrefix("/static/images/", imageServer))
    
	mux.HandleFunc("/", handlers.HandleHomepage)
	mux.HandleFunc("/register", handlers.RegisterUser)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/logout", handlers.LogoutUser)
	mux.HandleFunc("/api/reaction", handlers.HandleComments)
	mux.HandleFunc("/api/posts", handlers.HandlePosts)
	mux.HandleFunc("/api/categories", handlers.HandleCategories)
	return mux, nil
}
