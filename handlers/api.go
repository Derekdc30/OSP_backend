package handlers

import (
	"backend/db"
	"backend/db/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	var request struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("surveys")

	// Using FindOne to check if token exists
	var result struct {
		Token string `bson:"token"`
	}

	err := collection.FindOne(
		context.Background(),
		bson.M{"token": request.Token},
		options.FindOne().SetProjection(bson.M{"token": 1}),
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid token", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "valid",
		"token":  result.Token,
	})
}

func HandleGetSurvey(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	collection := db.GetCollection("surveys")

	var survey models.Survey
	err := collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&survey)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey)
}
