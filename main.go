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
	http.Handle("/app/", http.StripPrefix("/app/",
		http.FileServer(http.Dir("./app"))))
	http.Handle("/search", FileSearchHandler(s))
	http.Handle("/meta/", FileMetaHandler(s, "/meta/"))
	http.Handle("/files/", FilesHandler(s, "/files/"))
	log.Println("Starting server at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
