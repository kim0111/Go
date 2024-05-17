package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kim0111/GoMidterm/pkg/apple/validator"
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

func (s StoreModel) GetAll(title string, from, to int, filters Filters) ([]*Store, Metadata, error) {

	// Retrieve all stores items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, title, description, address, coordinates, number_of_branches
		FROM stores
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		AND (number_of_branches >= $2 OR $2 = 0)
		AND (number_of_branches <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{title, from, to, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			s.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var stores []*Store
	for rows.Next() {
		var store Store
		err := rows.Scan(&totalRecords, &store.Id, &store.CreatedAt, &store.UpdatedAt, &store.Title, &store.Description, &store.Address, &store.Coordinates, &store.NumberOfBranches)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		stores = append(stores, &store)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return stores, metadata, nil
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

func (s StoreModel) Get(id int) (*Store, error) {
	// Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, updated_at, title, description, address, coordinates, number_of_branches
		FROM stores
		WHERE id = $1
		`
	var store Store
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&store.Id, &store.CreatedAt, &store.UpdatedAt, &store.Title, &store.Description, &store.Address, &store.Coordinates, &store.NumberOfBranches)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive store with id: %v, %w", id, err)
	}
	return &store, nil
}

func (s StoreModel) Update(store *Store) error {
	query := `
		UPDATE stores
		SET title = $1, description = $2, address = $3, coordinates = &4, number_of_branches = &5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND updated_at = $7
		RETURNING updated_at
		`
	args := []interface{}{store.Title, store.Description, store.Address, store.Coordinates, store.NumberOfBranches, store.Id, store.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&store.UpdatedAt)
}

func (p StoreModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM stores
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateStore(v *validator.Validator, store *Store) {
	// Check if the title field is empty.
	v.Check(store.Title != "", "title", "must be provided")
	// Check if the title field is not more than 100 characters.
	v.Check(len(store.Title) <= 100, "title", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(store.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	// Check if the forWhatCountry is not more than 10 characters.
	v.Check(len(store.Address) <= 500, "Address", "must not be more than 500 bytes long")
	v.Check(len(store.Coordinates) <= 100, "Address", "must not be more than 100 bytes long")
	v.Check(store.NumberOfBranches <= 500, "numberOfBranches", "must not be more than 1000")
}
