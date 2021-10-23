package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")
}

func fileHandler(db *sql.DB, objects_path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
