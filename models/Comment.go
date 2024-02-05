package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID      primitive.ObjectID `bson:"id"`
	User_id primitive.ObjectID `bson:"user_id"`
	Post_id primitive.ObjectID `bson:"post_id"`
	Content string             `bson:"content"`
}
