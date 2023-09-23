package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"notetaking/database"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email"`
	Note  string             `json:"note"`
}

type JustNote struct {
	ID   primitive.ObjectID `json:"id"`
	Note string             `json:"note"`
}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var noteReq struct {
		SID  string `json:"sid"`
		Note string `json:"note"`
	}

	// Decode the request body into the noteReq struct
	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the JWT token from the request JSON
	tokenString := noteReq.SID

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userEmail := claims["sub"].(string)

	// Get a handle to the MongoDB collection for notes
	notesCollection := database.MongoClient.Database(database.DbName).Collection(database.NotesCollectionName)

	// Create a new note document
	note := Note{
		Email: userEmail,
		Note:  noteReq.Note,
	}

	// Insert the note document into MongoDB
	result, err := notesCollection.InsertOne(context.TODO(), note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID of the newly created note
	resp := map[string]interface{}{
		"id": result.InsertedID.(primitive.ObjectID).Hex(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
