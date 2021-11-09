package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
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
	basic := flag.Bool("basic", false, "enable Basic Authentication, requires you to set USERNAME and PASSWORD environment variables")

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
	if *basic {
		username := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")
		if username == "" || password == "" {
			log.Fatal("expected USERNAME and PASSWORD enviornment variables to be set")
		}
		http.Handle("/", basicAuth(Handler(s, tmpls), username, password))
	} else {
		http.Handle("/", Handler(s, tmpls))
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

func basicAuth(next http.Handler, username string, password string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUsername, reqPassword, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		usernameHash := sha256.Sum256([]byte(reqUsername))
		passwordHash := sha256.Sum256([]byte(reqPassword))
		expectedUsernameHash := sha256.Sum256([]byte(username))
		expectedPasswordHash := sha256.Sum256([]byte(password))

		// ConstantTimeCompare is use to avoid leaking information using timing attacks
		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)
		if usernameMatch && passwordMatch {
			next.ServeHTTP(w, r)
			return
		}

	})
}
