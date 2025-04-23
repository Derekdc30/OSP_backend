package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Survey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title" validate:"required"`
	Token     string             `bson:"token" validate:"required,len=5"`
	Questions []Question         `bson:"questions" validate:"required,min=1"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
