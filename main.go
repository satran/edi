package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewStore(root)
	if err != nil {
		log.Fatal(err)
	}
	if err := migrateDB(s.DB); err != nil {
		log.Fatal(err)
	}

	http.Handle("/s/", http.StripPrefix("/s/",
		http.FileServer(http.Dir("./s"))))

	http.Handle("/", Handler(s))
	log.Println("Starting server at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
