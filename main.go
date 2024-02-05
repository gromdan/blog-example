package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my.com/services"
)

var g_client *mongo.Client

func get_post_service() *services.PostService {
	s := &services.PostService{}
	s.SetContext(g_client)
	return s
}

func get_user_service() *services.UserService {
	u := &services.UserService{}
	u.SetContext(g_client)
	return u
}

func get_comment_service() *services.CommentService {
	c := &services.CommentService{}
	c.SetContext(g_client)
	return c
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	err = http.ListenAndServe(":8080", Install())
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func Install() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handleIndex).Methods("PUT")

	router.HandleFunc("/user", handlePutUser).Methods("PUT")
	router.HandleFunc("/post", handlePutPost).Methods("PUT")
	router.HandleFunc("/post", handleGetPost).Methods("GET")
	router.HandleFunc("/comment", handlePutComment).Methods("PUT")
	return router
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Hello, World!")
}

func handlePutUser(w http.ResponseWriter, r *http.Request) {
	get_user_service().Create(w, r)
}

func handlePutPost(w http.ResponseWriter, r *http.Request) {
	get_post_service().Create(w, r)
}

func handlePutComment(w http.ResponseWriter, r *http.Request) {
	get_comment_service().Create(w, r)
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	get_post_service().GetAll(w, r)
}
