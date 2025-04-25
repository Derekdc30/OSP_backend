package handlers

import (
	"backend/db"
	"backend/db/models"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// HandleVerifyAdmin verifies admin token
func HandleVerifyAdmin(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Compare with environment variable
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		http.Error(w, "Admin not configured", http.StatusInternalServerError)
		return
	}

	if request.Token != adminToken {
		http.Error(w, "Invalid admin token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "authorized"})
}

// HandleGetAllSurveys retrieves all surveys
func HandleGetAllSurveys(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("surveys")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var surveys []models.Survey
	if err = cursor.All(context.Background(), &surveys); err != nil {
		http.Error(w, "Decode error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveys)
}

// HandleDeleteSurvey deletes a survey by token
func HandleDeleteSurvey(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("surveys")
	result, err := collection.DeleteOne(context.Background(), bson.M{"token": token})
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleGetSurveyResponses retrieves survey responses
func HandleGetSurveyResponses(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	// 1. Get survey by token
	surveyCollection := db.GetCollection("surveys")
	var survey models.Survey
	err := surveyCollection.FindOne(
		context.Background(),
		bson.M{"token": token},
	).Decode(&survey)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Survey not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 2. Get responses using survey ID
	responseCollection := db.GetCollection("responses")
	cursor, err := responseCollection.Find(
		context.Background(),
		bson.M{"surveyToken": survey.Token},
	)

	if err != nil {
		http.Error(w, "Failed to fetch responses", http.StatusInternalServerError)
		return
	}

	var responses []models.Response
	if err = cursor.All(context.Background(), &responses); err != nil {
		http.Error(w, "Failed to decode responses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// HandleCreateSurvey creates a new survey
func HandleCreateSurvey(w http.ResponseWriter, r *http.Request) {
	var survey struct {
		Title     string            `json:"title"`
		Questions []models.Question `json:"questions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate token and timestamps
	newSurvey := models.Survey{
		Title:     survey.Title,
		Token:     utils.GenerateToken(5),
		Questions: survey.Questions,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := db.GetCollection("surveys")
	result, err := collection.InsertOne(context.Background(), newSurvey)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    result.InsertedID,
		"token": newSurvey.Token,
	})
}

// HandleUpdateSurvey updates a survey
func HandleUpdateSurvey(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	var survey struct {
		Title     string            `json:"title"`
		Questions []models.Question `json:"questions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(survey.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	for _, q := range survey.Questions {
		if strings.TrimSpace(q.Text) == "" {
			http.Error(w, "Question text is required", http.StatusBadRequest)
			return
		}
		if q.Format == "multiple_choice" && len(q.Options) < 2 {
			http.Error(w, "Multiple choice questions require at least 2 options", http.StatusBadRequest)
			return
		}
		if q.Format == "likert" && len(q.LikertScale) < 2 {
			http.Error(w, "Likert scale questions require at least 2 items", http.StatusBadRequest)
			return
		}
	}

	collection := db.GetCollection("surveys")
	update := bson.M{
		"$set": bson.M{
			"title":     survey.Title,
			"questions": survey.Questions,
			"updatedAt": time.Now(),
		},
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"token": token},
		update,
	)
	if err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Survey updated successfully",
		"token":   token,
	})
}
