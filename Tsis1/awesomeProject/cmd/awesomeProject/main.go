package main

import (
	"awesomeProject/internal/app/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/characters", handlers.GetCharacterList).Methods("GET")
	router.HandleFunc("/characters/{name}", handlers.GetCharacterDetails).Methods("GET")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
