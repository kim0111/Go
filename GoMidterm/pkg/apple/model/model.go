package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Products ProductModel
	Stores   StoreModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Products: ProductModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Stores: StoreModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
