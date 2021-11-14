package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func Handler(s *Store, tmpls *template.Template) http.HandlerFunc {
	get := FileGetHandler(s, tmpls)
	new_ := NewHandler(s, tmpls)
	edit := EditHandler(s, tmpls, "/edit/")
	update := FileWriteHandler(s, "/edit/")
	write := FileWriteHandler(s, "/_new")
	shell := ShellHandler(s)
	menu := MenuHandler(s, tmpls)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &ResponseWriter{ResponseWriter: w}
		switch {
		case r.URL.Path == "/":
			if shouldReject(w, r, http.MethodGet) {
				return
			}
			http.Redirect(w, r, "/"+s.Index(), http.StatusTemporaryRedirect)

		case r.URL.Path == "/_menu":
			if shouldReject(w, r, http.MethodGet) {
				return
			}
			menu(wr, r)

		case r.URL.Path == "/_sh":
			if shouldReject(w, r, http.MethodPost) {
				return
			}
			shell(wr, r)

		case r.URL.Path == "/_new":
			switch r.Method {
			case http.MethodGet:
				new_(wr, r)
			case http.MethodPost:
				write(wr, r)
			}

		case strings.HasPrefix(r.URL.Path, "/edit"):
			switch r.Method {
			case http.MethodGet:
				edit(wr, r)
			case http.MethodPost:
				update(wr, r)
			}
		default:
			switch r.Method {
			case http.MethodGet:
				get(wr, r)
			case http.MethodPost:
				write(wr, r)
			}
		}
		since := time.Since(start)
		log.Println(since, wr.StatusCode(), r.Method, r.URL)
	}
}

func NewHandler(s *Store, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpls.ExecuteTemplate(w, "new.html", nil); err != nil {
			log.Printf("executing list template: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func EditHandler(s *Store, tmpls *template.Template, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(path):]
		file, err := s.Get(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := tmpls.ExecuteTemplate(w, "new.html", file); err != nil {
			log.Printf("executing list template: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func FileGetHandler(s *Store, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimLeft(r.URL.Path, "/")
		if len(name) < 1 {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		f, err := s.Get(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer f.Close()
		if f.Type == "text/plain" {
			if err := tmpls.ExecuteTemplate(w, "file.html", f); err != nil {
				log.Printf("executing file template: %s", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		} else {
			http.ServeFile(w, r, s.path(name))
		}
	}
}

func MenuHandler(s *Store, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := s.Get(s.config.MenuFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer f.Close()
		if err := tmpls.ExecuteTemplate(w, "menu.html", f); err != nil {
			log.Printf("executing file template: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func FileWriteHandler(s *Store, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil { //10 MB
			log.Printf("couldn't parse form: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		var rdr io.Reader
		name := r.URL.Path[len(path):]
		text := r.PostForm.Get("text")
		if text != "" {
			if name == "" {
				name = r.PostForm.Get("name")
			}
			// This causes scripts to act randomly. awk fails not understanding a \r
			text = strings.ReplaceAll(text, "\r", "")
			rdr = strings.NewReader(text)
		} else {
			file, meta, err := r.FormFile("file")
			if err != nil {
				log.Printf("Error Retrieving the File from form: %s", err)
				writeError(w, http.StatusBadRequest)
				return
			}
			defer file.Close()
			if name == "" {
				name = meta.Filename
			}
			rdr = file
		}
		if name == "" {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		if err := s.Write(name, rdr); err != nil {
			log.Printf("creating file: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/"+name, http.StatusSeeOther)
	}
}

func ShellHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Cmd string `json:"cmd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		out := run(filepath.Join(s.root, "objects"), input.Cmd)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"output": out}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func shouldReject(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return false
	}
	writeError(w, http.StatusMethodNotAllowed)
	return true
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// ResponseWriter is a wrapper for http.ResponseWriter to get the written http status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *ResponseWriter) StatusCode() int {
	if r.statusCode == 0 {
		return http.StatusOK
	}
	return r.statusCode
}
