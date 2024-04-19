package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kim0111/GoMidterm/pkg/apple/model"
	"log"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createStoresHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title            string `json:"title"`
		Description      string `json:"description"`
		Address          string `json:"address"`
		Coordinates      string `json:"coordinates"`
		NumberOfBranches uint   `json:"numberOfBranches"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	store := &model.Store{
		Title:            input.Title,
		Description:      input.Description,
		Address:          input.Address,
		Coordinates:      input.Coordinates,
		NumberOfBranches: input.NumberOfBranches,
	}

	err = app.models.Stores.Insert(store)
	if err != nil {
		log.Printf("%s", err)
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, store)
}

func (app *application) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["storeId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := app.models.Stores.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, store)
}

func (app *application) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["storeId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := app.models.Stores.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title            string `json:"title"`
		Description      string `json:"description"`
		Address          string `json:"address"`
		Coordinates      string `json:"coordinates"`
		NumberOfBranches uint   `json:"numberOfBranches"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if &input.Title != nil {
		store.Title = input.Title
	}

	if &input.Description != nil {
		store.Description = input.Description
	}

	if &input.Address != nil {
		store.Address = input.Address
	}

	if &input.Coordinates != nil {
		store.Coordinates = input.Coordinates
	}

	if &input.NumberOfBranches != nil {
		store.NumberOfBranches = input.NumberOfBranches
	}

	err = app.models.Stores.Update(store)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, store)
}

func (app *application) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["storeId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	err = app.models.Stores.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
