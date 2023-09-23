package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"notetaking/database"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Session struct {
	UserID string `json:"user_id"`
	SID    string `json:"sid"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var JWTSecret = []byte("ntoheurcdboeunthonetu")
var sessionTimeout = 30 * time.Minute

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Decode the request body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get a handle to the MongoDB collection
	collection := database.MongoClient.Database(database.DbName).Collection(database.UsersCollectionName)

	// Find the user by email and password (you should hash and salt the password in production)
	var user User
	err = collection.FindOne(context.TODO(), bson.M{
		"email":    strings.ToLower(loginReq.Email),
		"password": loginReq.Password,
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a unique session ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email, // Subject (can be a user identifier)
		"exp": time.Now().Add(sessionTimeout).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"sid": tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
