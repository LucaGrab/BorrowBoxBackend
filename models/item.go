package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	TagNames      []string           `bson:"tagNames" json:"tagNames"`
	Name          string             `json:"name"`
	Location      string             `json:"location"`
	Description   string             `json:"description"`
	Available     bool               `json:"available"`
	CurrentRenter string             `json:"currentRenter"`
	RentedSince   string             `json:"rentedSince"`
}
type ItemMitReport struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	TagNames            []string           `bson:"tagNames" json:"tagNames"`
	Name                string             `json:"name"`
	Location            string             `json:"location"`
	Description         string             `json:"description"`
	Available           bool               `json:"available"`
	CurrentRenter       string             `json:"currentRenter"`
	RentedSince         string             `json:"rentedSince"`
	ReportDescription   string             `json:"reportDescription"`
	ReportTime          string             `json:"reportTime"`
	ReportUser          string             `json:"reportUser"`
	ReportStateCritical bool               `json:"reportStateCritical"`
}

type ItemMitTagIds struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	TagIds      []primitive.ObjectID `bson:"tagIds" json:"tagIds"`
	Name        string               `json:"name"`
	Location    string               `json:"location"`
	Description string               `json:"description"`
}

type AddItem struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Name        string               `json:"name"`
	Location    string               `json:"location"`
	Description string               `json:"description"`
	TagIds      []primitive.ObjectID `json:"tagIds"`
}

type ItemForInsert struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name"`
	Location    string             `json:"location"`
	Description string             `json:"description"`
}

type ItemForInsertMitDeleted struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name"`
	Location    string             `json:"location"`
	Description string             `json:"description"`
	Deleted     bool               `json:"deleted"`
}
