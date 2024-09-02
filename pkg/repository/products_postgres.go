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
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	query := `INSERT INTO product (name, description, price, quantity) values ($1, $2, $3, $4) RETURNING id`

	row := tx.QueryRow(query, products.Name, products.Description, products.Price, products.Quantity)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AuthPostgres) GetById(productId int) (Products, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return Products{}, err
	}

	var item Products
	query := `SELECT id, name, description, price, quantity FROM product WHERE id = $1`
	if err := tx.QueryRow(query, productId).Scan(&item.Id, &item.Name, &item.Description, &item.Price, &item.Quantity); err != nil {
		tx.Rollback()
		return item, err
	}

	if err := tx.Commit(); err != nil {
		return item, err
	}

	return item, nil
}
func (r *AuthPostgres) GetAll() ([]Products, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return []Products{}, err
	}

	var items []Products
	query := `SELECT id, name, description, price, quantity FROM product`
	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Products
		err = rows.Scan(&item.Id, &item.Name, &item.Description, &item.Price, &item.Quantity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		items = append(items, item)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return items, rows.Err()
}
func (r *AuthPostgres) Update(productId int, input UpdateProducts) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `UPDATE product SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5`

	_, err = r.db.Exec(query, input.Name, input.Description, input.Price, input.Quantity, productId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (r *AuthPostgres) Delete(productId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `DELETE FROM product WHERE id = $1`
	_, err = tx.Exec(query, productId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
