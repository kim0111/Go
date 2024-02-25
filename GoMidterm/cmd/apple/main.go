package main

import (
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/kim0111/GoMidterm/pkg/apple/model"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:8956@localhost:5432/lab2", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Println("Successfully connected to the database")
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/products", app.createProductsHandler).Methods("POST")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.getProductHandler).Methods("GET")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.updateProductHandler).Methods("PUT")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.deleteProductHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
