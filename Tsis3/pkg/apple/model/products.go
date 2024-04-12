package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kim0111/GoMidterm/pkg/apple/validator"
	"log"
	"time"
)

type Products struct {
	Id             string `json:"id"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	ForWhatCountry string `json:"forWhatCountry"`
	Price          uint   `json:"price"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (p ProductModel) GetAll(title string, from, to int, filters Filters) ([]*Products, Metadata, error) {

	// Retrieve all products items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, title, description, for_what_country, price
		FROM products
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		AND (price >= $2 OR $2 = 0)
		AND (price <= $3 OR $3 = 0)
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
	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			p.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var products []*Products
	for rows.Next() {
		var prod Products
		err := rows.Scan(&totalRecords, &prod.Id, &prod.CreatedAt, &prod.UpdatedAt, &prod.Title, &prod.Description, &prod.ForWhatCountry, &prod.Price)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		products = append(products, &prod)
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
	return products, metadata, nil
}

func (p ProductModel) Insert(product *Products) error {
	query := `
		INSERT INTO products (title, description, for_what_country, price) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{product.Title, product.Description, product.ForWhatCountry, product.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt)
}

func (p ProductModel) Get(id int) (*Products, error) {
	// Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, updated_at, title, description, for_what_country, price
		FROM products
		WHERE id = $1
		`
	var product Products
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.ForWhatCountry, &product.Price)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive menu with id: %v, %w", id, err)
	}
	return &product, nil
}

func (p ProductModel) Update(product *Products) error {
	query := `
		UPDATE products
		SET title = $1, description = $2, for_what_country = $3, price = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND updated_at = $6
		RETURNING updated_at
		`
	args := []interface{}{product.Title, product.Description, product.ForWhatCountry, product.Price, product.Id, product.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.UpdatedAt)
}

func (p ProductModel) Delete(id int) error {
	// Return an error if the ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateProduct(v *validator.Validator, prod *Products) {
	// Check if the title field is empty.
	v.Check(prod.Title != "", "title", "must be provided")
	// Check if the title field is not more than 100 characters.
	v.Check(len(prod.Title) <= 100, "title", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(prod.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	// Check if the forWhatCountry is not more than 10 characters.
	v.Check(len(prod.ForWhatCountry) <= 10, "forWhatCountry", "must not be more than 10 bytes long")
	// Check if the nutrition value is not more than 10000.
	v.Check(prod.Price <= 1000, "price", "must not be more than 1000")
}
