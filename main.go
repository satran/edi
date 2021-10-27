package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

//go:embed s templates
var contents embed.FS

func main() {
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewStore(root)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	if err := migrateDB(s.DB); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{Addr: "127.0.0.1:8080"}
	http.Handle("/s/", http.FileServer(http.FS(contents)))
	http.Handle("/", Handler(s))

	go func() {
		log.Println("Starting server at :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("shutting down")
	s.Close()
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
