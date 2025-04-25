package handlers

import (
	"backend/db"
	"backend/db/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandleSubmitResponse handles survey response submission
func HandleSubmitResponse(w http.ResponseWriter, r *http.Request) {
	var response struct {
		SurveyToken string `json:"surveyToken"`
		Answers     []struct {
			QuestionID string      `json:"questionId"`
			Value      interface{} `json:"value"`
		} `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get survey by token
	surveyCollection := db.GetCollection("surveys")
	var survey models.Survey
	err := surveyCollection.FindOne(
		context.Background(),
		bson.M{"token": response.SurveyToken},
	).Decode(&survey)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	// Convert question IDs to ObjectIDs
	var answers []models.Answer
	for _, a := range response.Answers {
		qid, err := primitive.ObjectIDFromHex(a.QuestionID)
		if err != nil {
			http.Error(w, "Invalid question ID format", http.StatusBadRequest)
			return
		}
		answers = append(answers, models.Answer{
			QuestionID: qid,
			Value:      a.Value,
		})
	}

	// Create response document
	responseDoc := models.Response{
		SurveyID:    survey.ID,
		SurveyToken: response.SurveyToken,
		Answers:     answers,
		SubmittedAt: time.Now(),
	}

	// Insert into responses collection
	responseCollection := db.GetCollection("responses")
	_, err = responseCollection.InsertOne(context.Background(), responseDoc)
	if err != nil {
		http.Error(w, "Failed to save response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Response recorded"})
}
