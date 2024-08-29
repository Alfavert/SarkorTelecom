package repository

import (
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewProductPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateProduct(products Products) (int, error) {
	var id int
	query := `INSERT INTO product (name, description, price, quantity) values ($1, $2, $3, $4) RETURNING id`

	row := r.db.QueryRow(query, products.Name, products.Description, products.Price, products.Quantity)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetById(productId int) (Products, error) {
	var item Products
	query := `SELECT id, name, description, price, quantity FROM product WHERE id = $1`
	if err := r.db.QueryRow(query, productId).Scan(&item.Id, &item.Name, &item.Description, &item.Price, &item.Quantity); err != nil {
		return item, err
	}

	return item, nil
}
func (r AuthPostgres) GetAll() ([]Products, error) {
	var items []Products
	query := `SELECT id, name, description, price, quantity FROM product`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Products
		err = rows.Scan(&item.Id, &item.Name, &item.Description, &item.Price, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *AuthPostgres) Update(productId int, input UpdateProducts) error {
	query := `UPDATE product SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5`

	_, err := r.db.Exec(query, input.Name, input.Description, input.Price, input.Quantity, productId)
	return err
}

func (r *AuthPostgres) Delete(productId int) error {
	query := `DELETE FROM product WHERE id = $1`
	_, err := r.db.Exec(query, productId)
	return err
}
