package main

import (
	"awesomeProject/internal/app/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/characters", handlers.GetCharacterList).Methods("GET")
	router.HandleFunc("/characters/{name}", handlers.GetCharacterDetails).Methods("GET")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
