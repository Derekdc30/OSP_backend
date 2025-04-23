package handlers

import (
	"fmt"
	"net/http"
)

// HandlePost handles POST /api/post
func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from form data
	token := r.FormValue("token")
	fmt.Printf("[SERVER] Received token: %s\n", token) // Display in terminal

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success: Valid token received!"))
}

// HandleCheckToken handles POST /api/check-token
func HandleCheckToken(w http.ResponseWriter, r *http.Request) {
	// Your token checking logic here
	w.Write([]byte("Token check endpoint"))
}
