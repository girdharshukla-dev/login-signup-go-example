package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func db_connection() (*sql.DB, error){
	pg_string := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_HOST"),
	)

	db, err := sql.Open("pgx", pg_string)

	if err != nil {
		return nil, err
	}
	
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}

func main() {
	
}
