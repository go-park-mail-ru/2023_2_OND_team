package ramrepo

import (
	"database/sql"
	"fmt"

	_ "github.com/proullon/ramsql/driver"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("ramsql", "RamRepository")
	if err != nil {
		return nil, err
	}

	err = createUsersTable(db)
	if err != nil {
		return nil, err
	}

	err = createPinTable(db)
	if err != nil {
		return nil, err
	}

	err = createSessionTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE users(
		id bigserial PRIMARY KEY,
		username varchar(30) UNIQUE,
		password varchar(50),
		email varchar(50) UNIQUE,
		avatar varchar(50)
	);`)
	if err != nil {
		return fmt.Errorf("create table users: %w", err)
	}
	return nil
}

func createPinTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE pin(
		id bigserial PRIMARY KEY,
		author int,
		picture varchar(50)
	);`)
	if err != nil {
		return fmt.Errorf("create table pin: %w", err)
	}
	return nil
}

func createSessionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE session(
		session_key varchar(30) PRIMARY KEY,
		user_id int,
		expire timestamp
	);`)
	if err != nil {
		return fmt.Errorf("create table session: %w", err)
	}
	return nil
}
