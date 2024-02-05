package services

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"my.com/models"
)

func (c *CommentService) SetContext(ctx *mongo.Client) {
	c.db_client = ctx
}

func create_comment(client *mongo.Client, dbName string, uid string, pid string, content string) error {
	collection := client.Database(dbName).Collection("comments")
	p_objectid, err := primitive.ObjectIDFromHex(pid)
	u_objectid, err := primitive.ObjectIDFromHex(uid)
	comment := models.Comment{User_id: u_objectid, Post_id: p_objectid, Content: content}
	_, err = collection.InsertOne(context.Background(), comment)
	return err
}

func (cs *CommentService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	uid := r.Form.Get("user_id")
	pid := r.Form.Get("post_id")
	content := r.Form.Get("content")

	_, err = get_user_byId(cs.db_client, "myBlog", uid)
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}

	_, err = get_post_byId(cs.db_client, "myBlog", pid)
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Invalid post_id",
		})
		return
	}
	if len(content) == 0 {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Content can not be empty",
		})
		return
	}
	err = create_comment(cs.db_client, "myBlog", uid, pid, "New comment content")
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Failed to create a new comment",
		})
	} else {
		write_response(w, Response{
			Type:    SUCCESS,
			Message: "",
		})
	}
}

type CommentService struct {
	CommentServiceInferface
	db_client *mongo.Client
}

type CommentServiceInferface interface {
	Create(post *models.Post)
	SetContext(ctx *mongo.Client)
}
