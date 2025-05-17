package main

import (
	"log"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
)

func main() {
	database.Init()
	defer database.Db.Close()
	//http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/", handlers.HandleLogin)
	http.HandleFunc("/signup", handlers.HandleSignup)

	log.Println("server starting: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
