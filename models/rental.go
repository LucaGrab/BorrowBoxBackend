package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rental struct {
	UserEmail string             `json:"userEmail"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Start     time.Time          `json:"start"`
	End       time.Time          `json:"end"`
	ItemId    primitive.ObjectID `json:"itemId" bson:"itemId"`
	Active    bool               `json:"active"`
}

type RentalWithId struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	UserEmail string             `json:"userEmail"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Start     time.Time          `json:"start"`
	End       time.Time          `json:"end"`
	ItemId    primitive.ObjectID `json:"itemId" bson:"itemId"`
	Active    bool               `json:"active"`
}

type ReturnData struct {
	ItemId   primitive.ObjectID `json:"itemId" bson:"itemId"`
	Location string             `json:"location"`
}
