package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Character struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Power string `json:"house"`
}

var characters = []Character{
	{"Rook", 40, "Armor"},
	{"Doc", 45, "Heal"},
	{"Sledge", 35, "Hammer"},
}

func GetCharacterList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

func GetCharacterDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	name := params["name"]

	for _, character := range characters {
		if character.Name == name {
			json.NewEncoder(w).Encode(character)
			return
		}
	}

	http.Error(w, "Character not found", http.StatusNotFound)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("App is healthy!"))
}
