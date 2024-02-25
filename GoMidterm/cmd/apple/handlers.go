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

func (app *application) createProductsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		ForWhatCountry string `json:"forWhatCountry"`
		Price          uint   `json:"price"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product := &model.Products{
		Title:          input.Title,
		Description:    input.Description,
		ForWhatCountry: input.ForWhatCountry,
		Price:          input.Price,
	}

	err = app.models.Products.Insert(product)
	if err != nil {
		log.Printf("%s", err)
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, product)
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, product)
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title          *string `json:"title"`
		Description    *string `json:"description"`
		ForWhatCountry *string `json:"forWhatCountry"`
		Price          *uint   `json:"price"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		product.Title = *input.Title
	}

	if &input.Description != nil {
		product.Description = *input.Description
	}

	if &input.ForWhatCountry != nil {
		product.ForWhatCountry = *input.ForWhatCountry
	}

	if &input.Price != nil {
		product.Price = *input.Price
	}

	err = app.models.Products.Update(product)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, product)
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = app.models.Products.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
