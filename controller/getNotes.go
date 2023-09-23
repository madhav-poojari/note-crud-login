package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"notetaking/database"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserNotesHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SID string `json:"sid"`
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

	// Get a handle to the MongoDB collection for notes
	notesCollection := database.MongoClient.Database(database.DbName).Collection(database.NotesCollectionName)

	// Define a filter to find notes for the user's email
	filter := bson.M{"email": userEmail}

	// Find all notes that match the filter
	cur, err := notesCollection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var notes []JustNote

	// Iterate through the cursor and decode notes
	for cur.Next(context.Background()) {
		var note Note
		err := cur.Decode(&note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, JustNote{ID: note.ID, Note: note.Note})
	}

	// Check for any errors that may have occurred during iteration
	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response with the notes array
	var response struct {
		Notes []JustNote `json:"notes"`
	}

	response.Notes = notes

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
