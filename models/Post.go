package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID      primitive.ObjectID `bson:"id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
	User_id int32              `bson:"user_id"`
}
