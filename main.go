package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/satran/edi/http"
)

func main() {
	docDir := flag.String("dir", "", "directory that stores the data, defaults to $HOME/notes")
	addr := flag.String("addr", "127.0.0.1:8080", "addr and port to serve from")
	start := flag.String("start-file", "Start", "file to render on index")
	flag.Parse()

	srv, err := http.Server(
		http.WithRootDir(*docDir),
		http.WithStartFile(*start),
		http.WithServerAddr(*addr),
	)
	if err != nil {
		log.Fatal(err)
	}
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
