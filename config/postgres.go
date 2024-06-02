package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func ConfigureDB() (*PostgresStore, error){
	fmt.Println("Connecting to Postgres")
	connection_string := "user=postgres dbname=postgres password=godeploy sslmode=disable"
	db, err := sql.Open( "postgres", connection_string)

	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
} 