package model

import (
	"context"
	"database/sql"
	"log"
	"time"
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

func (p StoreModel) Insert(store *Store) error {
	query := `
		INSERT INTO stores (title, description, address, coordinates, number_of_branches) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{store.Title, store.Description, store.Address, store.Coordinates, store.NumberOfBranches}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&store.Id, &store.CreatedAt, &store.UpdatedAt)
}

func (p StoreModel) Get(id int) (*Store, error) {
	query := `
		SELECT id, created_at, updated_at, title, description, address, coordinates, number_of_branches
		FROM stores
		WHERE id = $1
		`
	var store Store
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&store.Id, &store.CreatedAt, &store.UpdatedAt, &store.Title, &store.Description, &store.Address, &store.Coordinates, store.NumberOfBranches)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (p StoreModel) Update(store *Store) error {
	query := `
		UPDATE stores
		SET title = $1, description = $2, address = $3, coordinates = &4, number_of_branches = &5
		WHERE id = $6
		RETURNING updated_at
		`
	args := []interface{}{store.Title, store.Description, store.Address, store.Coordinates, store.NumberOfBranches, store.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&store.UpdatedAt)
}

func (p StoreModel) Delete(id int) error {
	query := `
		DELETE FROM stores
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}
