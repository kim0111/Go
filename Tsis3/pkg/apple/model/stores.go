package model

import (
	"database/sql"
	"errors"
	"log"
)

type Store struct {
	Id               string `json:"id"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Address          string `json:"address"`
	Coordinates      string `json:"coordinates"`
	NumberOfBranches uint   `json:"numberOfBranches"`
}

type StoreModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var stores = []Store{
	{
		Id:               "1",
		Title:            "TechnoDom",
		Address:          "Some street in almaty",
		NumberOfBranches: 3,
	},
	{
		Id:               "2",
		Title:            "iPoint",
		Address:          "Some street in almaty 2",
		NumberOfBranches: 2,
	},
	{
		Id:               "3",
		Title:            "Sulpak",
		Address:          "Some street in almaty 3",
		NumberOfBranches: 7,
	},
}

func GetStores() []Store {
	return stores
}

func GetStore(id string) (*Store, error) {
	for _, r := range stores {
		if r.Id == id {
			return &r, nil
		}
	}
	return nil, errors.New("Store not found")
}
