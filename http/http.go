package http

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/satran/edi/exec"
	"github.com/satran/edi/store"
)

func getH(s *store.Store, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if shouldReject(w, r, http.MethodGet) {
				return
			}
			http.Redirect(w, r, "/"+s.Index(), http.StatusTemporaryRedirect)
			return
		}
		name := strings.TrimLeft(r.URL.Path, "/")
		if name == "" {
			writeError(w, http.StatusNotFound)
			return
		}

		f, err := s.Get(name)
		if err != nil {
			// File doesn't exist, let's render a template with the name
			f = store.Dummy(s.Root, name)
		}
		if f.IsText() {
			if err := tmpls.ExecuteTemplate(w, "file.html", f); err != nil {
				log.Printf("executing file template: %s", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		} else {
			http.ServeFile(w, r, s.Path(name))
		}
	}

}

const blobDir = "_blob"

func getBlobH(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimLeft(r.URL.Path, "/")
		http.ServeFile(w, r, s.Path(name))
	}
}

func addBlobH(s *store.Store, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil { //10 MB
			log.Printf("couldn't parse form: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		var rdr io.Reader
		name := r.URL.Path[len(path):]
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
		if name == "" {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		name = filepath.Join(blobDir, name)
		if err := s.Write(name, rdr); err != nil {
			log.Printf("creating file: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/"+name, http.StatusSeeOther)
	}
}

func editH(s *store.Store, tmpls *template.Template, path string) http.HandlerFunc {
	update := fileWriteHandler(s, path)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			name := r.URL.Path[len(path):]
			file, err := s.Get(name)
			if err != nil {
				// File doesn't exist, let's render a template with the name
				file = store.Dummy(s.Root, name)
			}
			if err := tmpls.ExecuteTemplate(w, "edit.html", file); err != nil {
				log.Printf("executing list template: %s", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			update(w, r)
		}
	}
}

func newH(s *store.Store, tmpls *template.Template, path string) http.HandlerFunc {
	newHandler := fileWriteHandler(s, path)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if err := tmpls.ExecuteTemplate(w, "edit.html", nil); err != nil {
				log.Printf("executing list template: %s", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			newHandler(w, r)
		}
	}
}

func fileWriteHandler(s *store.Store, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil { //10 MB
			log.Printf("couldn't parse form: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		var rdr io.Reader
		name := r.URL.Path[len(path):]
		text := r.PostForm.Get("text")
		if name == "" {
			name = r.PostForm.Get("name")
		}
		if name == "" {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// This causes scripts to act randomly. awk fails not understanding a \r
		text = strings.ReplaceAll(text, "\r", "")
		rdr = strings.NewReader(text)
		if err := s.Write(name, rdr); err != nil {
			log.Printf("creating file: %s", err)
			writeError(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/"+name, http.StatusSeeOther)
	}
}

func shellH(s *store.Store, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := tmpls.ExecuteTemplate(w, "shell.html", nil); err != nil {
				log.Printf("executing shell template: %s", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
			return
		}
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed)
		}
		var input struct {
			Cmd string `json:"cmd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		out := exec.Run(s.Root, "", input.Cmd)
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

func logRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &ResponseWriter{ResponseWriter: w}
		fn(wr, r)
		since := time.Since(start)
		log.Println(since, wr.StatusCode(), r.Method, r.URL)
	}
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
