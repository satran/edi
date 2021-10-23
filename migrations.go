package main

import (
	"database/sql"
	"fmt"
)

func migrateDB(db *sql.DB) error {
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

var migrations = []string{`
create table if not exists files (
    id integer primary key
  , object_id text not null
  , created_at integer not null
  , updated_at integer not null
  , content_type text not null
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
