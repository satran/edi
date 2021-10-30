package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

//go:embed s templates
var contents embed.FS

func main() {
	static := flag.String("static", "", "Render static files from folder")
	dir := flag.String("dir", ".", "directory that stores the data")
	addr := flag.String("addr", "127.0.0.1:8080", "addr and port to serve from")
	flag.Parse()
	root, err := filepath.Abs(*dir)
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

	srv := &http.Server{Addr: *addr}
	if *static != "" {
		http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(*static))))
	} else {
		http.Handle("/s/", http.FileServer(http.FS(contents)))
	}
	http.Handle("/", Handler(s))

	go func() {
		log.Printf("Starting server %s", *addr)
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
