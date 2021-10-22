package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const migrations = []string{`
create table if not exists files (
    id integer primary key
  , object_id text not null
  , created_at integer not null
  , updated_at integer not null
)`, `
create table if not exists log (
    file_id integer not null references files(id) 
  , object_id text not null
  , updated_at integer not null
)`, `
create table if not exists tags (
    file_id integer not null references files(id)
  , tag text not null
)`,
}

func openDB(path string) (*sql.DB, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	f.Close()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("could open db: %w", err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func migrateDB(db *sql.DB) error {
	println("migrateDB")
	r, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("can't open migration file: %s, %w", migrationPath, err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	//defer tx.Rollback()
	for _, stmt := range migrations {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("can't run migration: %q: %w", stmt, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migrations: %w", err)
	}
	return nil
}
