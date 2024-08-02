package db

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var tables = []string{
	`CREATE TABLE IF NOT EXISTS certifications (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS domains (
		certification_id INTEGER NOT NULL,
		number INTEGER NOT NULL,
		name TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS tasks (
		certification_id INTEGER NOT NULL,
		domain_number INTEGER NOT NULL,
		number INTEGER NOT NULL,
		name TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY,
		certification_id INTEGER NOT NULL,
		domain_number INTEGER NOT NULL,
		task_number INTEGER NOT NULL,
		question TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS answers (
		id INTEGER PRIMARY KEY,
		question_id INTEGER NOT NULL,
		answer TEXT NOT NULL,
		correct BOOLEAN NOT NULL,
		explanation TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE
	)`,
	`CREATE TABLE IF NOT EXISTS responses (
		user_id INTEGER NOT NULL,
		question_id INTEGER NOT NULL,
		answer_id INTEGER NOT NULL,
		date DATETIME NOT NULL
	)`,
}

func NewSQLiteDB(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}
	db := &DB{conn}
	return db, initialize_sqlite(conn)
}

func initialize_sqlite(conn *sql.DB) error {
	var err, errs error
	for _, table := range tables {
		_, err = conn.Exec(table)
		errs = errors.Join(errs, err)
	}
	return errs
}
