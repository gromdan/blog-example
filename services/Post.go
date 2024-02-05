package services

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"my.com/models"
)

func (s *PostService) SetContext(ctx *mongo.Client) {
	s.db_client = ctx
}

func get_post_byId(client *mongo.Client, dbName string, id string) (models.Post, error) {
	var post models.Post
	collection := client.Database(dbName).Collection("posts")
	filter := bson.D{{Key: "id", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(&post)
	return post, err
}

func creat_post(client *mongo.Client, dbName string, title string, content string, uid int32) error {
	collection := client.Database(dbName).Collection("posts")
	post := models.Post{
		Title:   title,
		Content: content,
		User_id: uid,
	}

	_, err := collection.InsertOne(context.Background(), post)
	return err
}

func (ps *PostService) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (ps *PostService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	uid := r.Form.Get("user_id")

	user, err := get_user_byId(ps.db_client, "myBlog", uid)
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}
	if !user.Author {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Has no Author permission",
		})
		return
	}
	if len(title) == 0 {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Title can not be empty",
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

	err = creat_post(ps.db_client, "myBlog", title, content, 123)
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Failed to create a new post",
		})
	} else {
		write_response(w, Response{
			Type:    SUCCESS,
			Message: "",
		})
	}
}

type PostServiceInferface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	SetContext(ctx *mongo.Client)
}

type PostService struct {
	PostServiceInferface
	db_client *mongo.Client
}
