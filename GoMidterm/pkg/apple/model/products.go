package model

import (
	"context"
	"database/sql"
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

func (p ProductModel) Insert(product *Products) error {
	query := `
		INSERT INTO products (title, description, forWhatCountry, price) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{product.Title, product.Description, product.ForWhatCountry, product.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt)
}

func (p ProductModel) Get(id int) (*Products, error) {
	query := `
		SELECT id, created_at, updated_at, title, description, forWhatCountry, price
		FROM products
		WHERE id = $1
		`
	var product Products
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.ForWhatCountry, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p ProductModel) Update(product *Products) error {
	query := `
		UPDATE products
		SET title = $1, description = $2, forWhatCountry = $3, price = &4
		WHERE id = $5
		RETURNING updated_at
		`
	args := []interface{}{product.Title, product.Description, product.ForWhatCountry, product.Price, product.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.UpdatedAt)
}

func (p ProductModel) Delete(id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}
