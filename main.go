package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func db_connection() (*sql.DB, error) {
	pg_string := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
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
	db, err := db_connection()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", signupHandler(db))

	fmt.Println("> Server on port 8080")
	http.ListenAndServe(":8080", mux)
}

func signupHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var user User
		if err := decoder.Decode(&user); err != nil {
			http.Error(w, "Invalid request", 400)
			return
		}

		password_hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			http.Error(w, "Internal error", 500)
		}

		_, err = db.Exec(
			`INSERT INTO users(email, password_hash) VALUES ($1, $2)`, 
			user.Email, 
			password_hash)
		
		if err != nil {
			http.Error(w, "Internal error", 500)
		}

		w.WriteHeader(201)

	}
}
