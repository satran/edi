package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/satran/edi/http"
	"github.com/webview/webview"
)

func main() {
	debug := flag.Bool("debug", false, "Enables the debug console on the webview")
	docDir := flag.String("dir", "", "directory that stores the data, defaults to $HOME/notes")
	start := flag.String("start-file", "Start", "file to render on index")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	username := randSeq(3)
	password := randSeq(10)
	addr := "localhost:3012"

	srv, err := http.Server(
		http.WithRootDir(*docDir),
		http.WithStartFile(*start),
		http.WithServerAddr(addr),
		http.WithBasicAuth(username, password),
	)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s:%s@%s", username, password, addr)
	go func() {
		if *debug {
			log.Println("Starting server:", url)
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	defer func() {
		log.Println("shutting down")
		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}()

	w := webview.New(*debug)
	defer w.Destroy()
	w.SetTitle("EDI")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(url)
	w.Run()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
