package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Text        string             `bson:"text" validate:"required"`
	Format      string             `bson:"format" validate:"required,oneof=textbox multiple_choice likert"`
	Options     []string           `bson:"options,omitempty"`      // For multiple choice
	LikertScale []string           `bson:"likert_scale,omitempty"` // For likert scale
	IsRequired  bool               `bson:"is_required"`
}
