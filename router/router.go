package router

import (
	"notetaking/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", controller.Signup).Methods("POST")
	router.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	router.HandleFunc("/notes", controller.GetUserNotesHandler).Methods("GET")
	router.HandleFunc("/notes", controller.PostNoteHandler).Methods("POST")
	router.HandleFunc("/notes", controller.DeleteUserNoteHandler).Methods("DELETE")

	return router
}
