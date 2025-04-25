package handlers

import (
	"backend/db"
	"backend/db/models"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HandleCheckToken validates survey token
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

// HandleGetSurvey retrieves survey by token
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
