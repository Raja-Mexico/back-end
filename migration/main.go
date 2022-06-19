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
	seed(db)
}

func migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) not null,
    	email VARCHAR(255) not null UNIQUE,
    	password VARCHAR(255) not null
		);

		CREATE TABLE IF NOT EXISTS financial_account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER not null,
			bank_id INTEGER not null,
			access_token VARCHAR(255) not null,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS team (	
			id VARCHAR(255) PRIMARY KEY,
			creator_id INTEGER not null,
			balance FLOAT not null default 0,
			prepaid_balance FLOAT not null default 0,
			name VARCHAR(255) null,
			FOREIGN KEY (creator_id) REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS membership (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER not null,
			team_id VARCHAR(255) not null,
			is_admin BOOLEAN not null default true,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (team_id) REFERENCES team(id)
		);

		CREATE TABLE IF NOT EXISTS user_balance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER not null,
			team_id INTEGER not null,
			balance FLOAT NOT NULL DEFAULT 0,
			no_virtual_account VARCHAR(255) NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(team_id) REFERENCES team(id)
		);

		CREATE TABLE IF NOT EXISTS top_up (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nominal BIGINT not null,
			date_time DATETIME not null,
			success_state SMALLINT not null,
			user_balance_id VARCHAR(255) not null,
			user_id INTEGER not null,
			FOREIGN KEY(user_balance_id) REFERENCES user_balance(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		
		CREATE TABLE IF NOT EXISTS service (
			id SMALLINT PRIMARY KEY,
			name VARCHAR(255) not null
		);

		CREATE TABLE IF NOT EXISTS status (
			id SMALLINT PRIMARY KEY,
			name VARCHAR(255) not null
		);

		CREATE TABLE IF NOT EXISTS prepaid_card (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			creator_id INTEGER not null,
			service_id SMALLINT not null,
			team_id varchar(255) not null,
			deadline_day INTEGER not null,
			is_automated BOOLEAN not null default false,
			nominal FLOAT not null,
			destination_number VARCHAR(255) not null,
			status_id SMALLINT not null default 1,
			FOREIGN KEY(creator_id) REFERENCES users(id),
			FOREIGN KEY(service_id) REFERENCES service(id),
			FOREIGN KEY(team_id) REFERENCES team(id),
			FOREIGN KEY(user_involved_id) REFERENCES users(id),
			FOREIGN KEY(status_id) REFERENCES status(id)
		);

		CREATE TABLE IF NOT EXISTS manage_prepaid_card (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			prepaid_card_id INTEGER not null,
			user_id INTEGER not null,
			is_confirmated BOOLEAN not null default false,
			pay_request float not null,
			FOREIGN KEY(prepaid_card_id) REFERENCES prepaid_card(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)

	if err != nil {
		panic(err)
	}
}

func seed(db *sql.DB) {
	_, err := db.Exec(`
		INSERT INTO service (id, name) VALUES 
		(1, 'PLN'),
		(2, 'Pulsa');

		INSERT INTO status (id, name) VALUES
		(1, 'Menunggu Pembayaran'),
		(2, 'Sudah Dibayar'),
		(3, 'Menunggu Persetujuan Pembayaran Bersama'),
		(4, 'Menunggu Konfirmasi (withdraw)');
	`)

	if err != nil {
		panic(err)
	}
}
