package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func Handler(s *Store) http.HandlerFunc {
	appHandler := AppHandler(s)
	getHandler := FileGetHandler(s)
	metaHandler := FileMetaHandler(s, "/meta/")
	createHandler := FileCreateHandler(s)
	updateHandler := FileUpdateHandler(s)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &ResponseWriter{ResponseWriter: w}
		switch {
		case r.URL.Path == "/":
			switch r.Method {
			case http.MethodGet:
				appHandler(wr, r)
			case http.MethodPost:
				createHandler(wr, r)
			}
		case strings.HasPrefix(r.URL.Path, "/meta"):
			metaHandler(wr, r)
		default:
			switch r.Method {
			case http.MethodGet:
				getHandler(wr, r)
			case http.MethodPost:
				updateHandler(wr, r)
			}

		}
		since := time.Since(start)
		log.Println(since, wr.StatusCode(), r.Method, r.URL)
	}
}

func AppHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "list.html")))

		files, err := s.Search(Query{})
		if err != nil {
			log.Printf("listing files: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
		ret := make([]File, 0, len(files))
		for _, f := range files {
			if f.Type == "text/plain" {
				content, err := s.GetText(f.ID)
				if err != nil {
					log.Printf("cant get text: %s", err)
					writeError(w, http.StatusInternalServerError)
					return
				}
				f.Content = content
			} else {
				f.Content = f.ID
			}
			ret = append(ret, f)
		}

		if err := listTemplate.ExecuteTemplate(w, "list.html",
			map[string]interface{}{
				"Files": ret,
			},
		); err != nil {
			log.Printf("executing list template: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func FileMetaHandler(s *Store, metaPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Path[len(metaPath):]
		f, err := s.Get(id)
		if err != nil {
			log.Printf("couldn't get file(%s): %s", id, err)
			writeError(w, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(&f); err != nil {
			log.Printf("writing json: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func FileSearchHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := Query{}
		files, err := s.Search(params)
		if err != nil {
			log.Printf("writing json: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&files); err != nil {
			log.Printf("writing json: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
}

func FilesHandler(s *Store, path string) http.HandlerFunc {
	getHandler := FileGetHandler(s)
	createHandler := FileCreateHandler(s)
	updateHandler := FileUpdateHandler(s)
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &ResponseWriter{ResponseWriter: w}
		switch r.Method {
		case http.MethodGet:
			getHandler(wr, r)
		case http.MethodPost:
			createHandler(wr, r)
		case http.MethodPut:
			updateHandler(wr, r)
		default:
			writeError(w, http.StatusNotImplemented)
		}
		since := time.Since(start)
		log.Println(since, wr.StatusCode(), r.Method, r.URL)
	}
}

func FileGetHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimLeft(r.URL.Path, "/")
		log.Println("get:", id)
		if len(id) < 1 {
			log.Println("empty ID requested")
			writeError(w, http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, getObjectPath(s.root, id))
	}
}

func FileCreateHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("couldn't parse form: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		text := r.PostForm.Get("text")
		if text == "" {
			return
		}
		sr := strings.NewReader(text)
		_, err := s.Create(sr)
		if err != nil {
			log.Printf("error creating: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 301)
	}
}

func FileUpdateHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[1:]
		f := struct {
			Content string `json:"content"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			log.Printf("can't decode json: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		sr := strings.NewReader(f.Content)
		if err := s.Update(id, sr); err != nil {
			log.Printf("error updating: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"id": "%s"}`, id)))
	}
}

func writeError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
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
