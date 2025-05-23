package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
)

func HandleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.Db.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			http.Error(w, "Error parsing categories", http.StatusInternalServerError)
			return
		}
		categories = append(categories, map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
} 