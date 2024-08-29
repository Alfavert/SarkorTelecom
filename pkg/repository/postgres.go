package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectPostgres(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func CreateProductTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS product(
    	id BIGSERIAL PRIMARY KEY,
    	name TEXT NOT NULL,
    	description TEXT NOT NULL,
    	price NUMERIC(6, 2) NOT NULL,
    	quantity INT NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
