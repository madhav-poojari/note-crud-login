package main

import (
	"fmt"
	"log"
	"net/http"

	"notetaking/database"
	"notetaking/router"
)

func main() {
	database.InitializeMongoDB()
	fmt.Println("Connected to MongoDB Atlas successfully")

	r := router.Router()

	err := http.ListenAndServe(":4068", r)
	if err != nil {
		log.Fatal(err)
	}
}
