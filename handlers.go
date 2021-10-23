package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FilesHandler(s *Store) http.HandlerFunc {
	getHandler := FileGetHandler(s, "")
	createHandler := FileCreateHandler(s, "")
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &ResponseWriter{ResponseWriter: w}
		switch r.Method {
		case http.MethodGet:
			getHandler(wr, r)
		case http.MethodPost:
			createHandler(wr, r)
		default:
			writeError(w, http.StatusNotImplemented)
		}
		since := time.Since(start)
		log.Println(since, wr.StatusCode(), r.Method, r.URL)
	}
}

func FileGetHandler(s *Store, objpath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.URL.Path[len("/files/"):]
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			log.Printf("couldn't parse id: %s, %s", sid, err)
			writeError(w, http.StatusNotFound)
			return
		}
		f, err := s.Get(id)
		if err != nil {
			log.Printf("couldn't get file(%d): %s", id, err)
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

func FileSearchHandler(s *Store, objpath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func FileCreateHandler(s *Store, objpath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := struct {
			Content string `json:"content"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			log.Printf("can't decode json: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		sr := strings.NewReader(f.Content)
		id, err := s.Create(sr)
		if err != nil {
			log.Printf("error creating: %s", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id)))
	}
}

func FileUpdateHandler(s *Store, objpath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
