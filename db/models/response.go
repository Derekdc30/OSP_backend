package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	SurveyID    primitive.ObjectID `bson:"survey_id" validate:"required"`
	SurveyToken string             `json:"surveyToken" bson:"surveyToken" validate:"required"`
	Answers     []Answer           `bson:"answers" validate:"required"`
	SubmittedAt time.Time          `bson:"submitted_at"`
}

type Answer struct {
	QuestionID primitive.ObjectID `bson:"question_id" validate:"required"`
	Value      interface{}        `bson:"value" validate:"required"`
}
