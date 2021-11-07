package main

import (
	"context"
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

//go:embed s templates
var contents embed.FS

func main() {
	staticDir := flag.String("static", "", "Render static files from folder")
	templateDir := flag.String("template", "", "Render template files from folder")
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

	var tmpls *template.Template
	if *templateDir != "" {
		tmpls = template.Must(template.ParseFiles(
			filepath.Join(*templateDir, "edit.html"),
			filepath.Join(*templateDir, "list.html"),
		))
	} else {
		tmpls = template.Must(template.ParseFS(contents, "templates/*"))
	}

	srv := &http.Server{Addr: *addr}
	if *staticDir != "" {
		http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(*staticDir))))
	} else {
		http.Handle("/s/", http.FileServer(http.FS(contents)))
	}
	http.Handle("/", Handler(s, tmpls))

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
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
