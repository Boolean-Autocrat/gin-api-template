package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(user, password, dbname, host string) (*Postgres, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s port=5432 host=%s sslmode=disable", user, password, dbname, host))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
