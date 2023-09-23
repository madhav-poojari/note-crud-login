package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"notetaking/database"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// var MongoClient *mongo.Client = database.MongoClient
// var DbName = database.DbName
// var CollectionName = database.CollectionName

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the request body into the User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get a handle to the MongoDB collection
	collection := database.MongoClient.Database(database.DbName).Collection(database.UsersCollectionName)

	// checking for duplicate email
	existingUser := User{}
	err = collection.FindOne(context.TODO(), bson.M{"email": strings.ToLower(user.Email)}).Decode(&existingUser)
	if err == nil {
		// User with the same email already exists
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		// An error occurred during the query
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the user data into MongoDB
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
