package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"notetaking/database"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SID string `json:"sid"`
		ID  string `json:"id"`
	}

	// Decode the request JSON into the req struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the JWT token from the request JSON
	tokenString := req.SID

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

	// Validate the ID provided in the request
	noteID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Get a handle to the MongoDB collection for notes
	notesCollection := database.MongoClient.Database(database.DbName).Collection(database.NotesCollectionName)

	// Define a filter to find the note to delete by its ID and user's email
	filter := bson.M{"_id": noteID, "email": userEmail}

	// Delete the note that matches the filter
	result, err := notesCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Note not found or unauthorized", http.StatusBadRequest)
		return
	}
	resp := map[string]string{
		"message": "succesfully deleted note",
	}
	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
