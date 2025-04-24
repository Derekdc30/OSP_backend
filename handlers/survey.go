package handlers

import (
	"backend/db"
	"backend/db/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateSurvey() {
	// Fixed token for testing
	token := "ABC12"

	// Sample questions
	questions := []models.Question{
		{
			Text:       "What is your name?",
			Format:     "textbox",
			IsRequired: true,
		},
		{
			Text:       "Choose your favorite color",
			Format:     "multiple_choice",
			Options:    []string{"Red", "Blue", "Green"},
			IsRequired: true,
		},
		{
			Text:        "Rate your satisfaction (1-5)",
			Format:      "likert",
			LikertScale: []string{"1", "2", "3", "4", "5"},
			IsRequired:  false,
		},
	}

	// Create survey object
	survey := models.Survey{
		Title:     "Sample Survey",
		Token:     token,
		Questions: questions,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into database
	collection := db.GetCollection("surveys")
	_, err := collection.InsertOne(context.Background(), survey)
	if err != nil {
		log.Printf("Failed to insert sample survey: %v", err)
		return
	}
	fmt.Printf("Successfully inserted sample survey with token: %s\n", token)
}

func TestGetSurvey(token string) {
	collection := db.GetCollection("surveys")
	filter := bson.M{"token": token}

	var survey models.Survey
	err := collection.FindOne(context.Background(), filter).Decode(&survey)
	if err != nil {
		log.Printf("Failed to retrieve survey with token %s: %v", token, err)
		return
	}

	// Display survey details
	fmt.Println("\nRetrieved Survey Details:")
	fmt.Printf("Title: %s\n", survey.Title)
	fmt.Printf("Token: %s\n", survey.Token)
	fmt.Println("Questions:")
	for i, q := range survey.Questions {
		fmt.Printf("  %d. %s\n", i+1, q.Text)
		fmt.Printf("     Format: %s\n", q.Format)
		if q.Format == "multiple_choice" {
			fmt.Printf("     Options: %v\n", q.Options)
		} else if q.Format == "likert" {
			fmt.Printf("     Likert Scale: %v\n", q.LikertScale)
		}
		fmt.Printf("     Required: %t\n", q.IsRequired)
	}
	fmt.Printf("Created At: %s\n", survey.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated At: %s\n", survey.UpdatedAt.Format(time.RFC3339))
}
