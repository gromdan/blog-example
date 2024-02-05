package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"my.com/models"
)

type Response struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var SUCCESS string = "success"
var ERROR string = "error"

func write_response(w http.ResponseWriter, resp Response) {
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
}

func hashed_password(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed_password), nil
}

func (u *UserService) SetContext(ctx *mongo.Client) {
	u.db_client = ctx
}

func get_user_byId(client *mongo.Client, dbName string, id string) (models.User, error) {
	collection := client.Database(dbName).Collection("users")
	filter := bson.D{{Key: "id", Value: id}}
	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func create_user(client *mongo.Client, dbName string, user models.User) error {
	collection := client.Database(dbName).Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func count_user_byName(client *mongo.Client, dbName string, username string) (int64, error) {
	collection := client.Database(dbName).Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.D{{Key: "username", Value: username}})
	return count, err
}

func (ps *UserService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	count, err := count_user_byName(ps.db_client, "myBlog", "username")
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Conflict Username",
		})
		return
	}
	if count != 0 {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Conflict Username",
		})
		return
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Password can not be empty",
		})
		return
	}
	password, _ = hashed_password(password)

	user := models.User{Username: username, Password: password, Author: true}
	err = create_user(ps.db_client, "mydatabase", user)
	if err != nil {
		write_response(w, Response{
			Type:    ERROR,
			Message: "Failed to create account",
		})
	} else {
		write_response(w, Response{
			Type:    SUCCESS,
			Message: "",
		})
	}
}

type PostUserInferface interface {
	Create(w http.ResponseWriter, r *http.Request)
	CheckUserById(Id int32)
	SetContext(ctx *mongo.Client)
}

type UserService struct {
	PostUserInferface
	db_client *mongo.Client
}
