package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Role     string             `json:"role"`
	Email    string             `json:"email"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}
