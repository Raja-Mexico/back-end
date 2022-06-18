package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "raja-mexico.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	migrate(db)
}

func migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name varchar(255) not null,
    	email varchar(255) not null UNIQUE,
    	password varchar(255) not null
		);
	`)

	if err != nil {
		panic(err)
	}
}
