package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	Description   string             `json:"description"`
	ItemId        primitive.ObjectID `json:"itemId" bson:"itemId"`
	UserId        primitive.ObjectID `json:"userId" bson:"userId"`
	Time          time.Time          `json:"time"`
	StateCritical bool               `json:"reportState"`
}
