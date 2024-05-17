package main

import (
	"errors"
	"github.com/kim0111/GoMidterm/pkg/apple/model"
	"github.com/kim0111/GoMidterm/pkg/apple/validator"
	"log"
	"net/http"
)

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
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
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
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"stores": store}, nil)
}

func (app *application) getStoresList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string
		BranchesFrom int
		BranchesTo   int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the title and nutrition value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Title = app.readStrings(qs, "title", "")
	input.BranchesFrom = app.readInt(qs, "branchesFrom", 0, v)
	input.BranchesTo = app.readInt(qs, "branchesTo", 0, v)

	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on store ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "title", "numberOfBranches",
		// descending sort values
		"-id", "-title", "-numberOfBranches",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	stores, metadata, err := app.models.Stores.GetAll(input.Title, input.BranchesFrom, input.BranchesTo, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"stores": stores, "metadata": metadata}, nil)
}

func (app *application) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	store, err := app.models.Stores.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"stores": store}, nil)
}

func (app *application) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	store, err := app.models.Stores.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
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
		app.badRequestResponse(w, r, err)
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

	v := validator.New()

	if model.ValidateStore(v, store); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Stores.Update(store)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"stores": store}, nil)
}

func (app *application) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Stores.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
