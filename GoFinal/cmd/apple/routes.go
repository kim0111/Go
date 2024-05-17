package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	prod1 := r.PathPrefix("/api/v1").Subrouter()
	store := r.PathPrefix("/api/v1").Subrouter()

	// Product Singleton
	// localhost:8081/api/v1/products
	prod1.HandleFunc("/products", app.getProductsList).Methods("GET")
	// Create a new prod
	prod1.HandleFunc("/products", app.createProductsHandler).Methods("POST")
	// Get a specific prod
	prod1.HandleFunc("/products/{id:[0-9]+}", app.getProductHandler).Methods("GET")
	// Update a specific prod
	prod1.HandleFunc("/products/{id:[0-9]+}", app.updateProductHandler).Methods("PUT")
	prod1.HandleFunc("/products/nopermission/{id:[0-9]+}", app.deleteProductHandler).Methods("DELETE")

	// Delete a specific prod
	prod1.HandleFunc("/products/{id:[0-9]+}", app.requirePermissions("products:write", app.deleteProductHandler)).Methods("DELETE")

	//Stores
	store.HandleFunc("/stores", app.getStoresList).Methods("GET")
	store.HandleFunc("/stores", app.createStoresHandler).Methods("POST")
	store.HandleFunc("/stores/{id:[0-9]+}", app.getStoreHandler).Methods("GET")
	store.HandleFunc("/stores/{id:[0-9]+}", app.updateStoreHandler).Methods("PUT")
	store.HandleFunc("/stores/{id:[0-9]+}", app.requirePermissions("products:write", app.deleteStoreHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
