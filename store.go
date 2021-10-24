package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewStore(root string) (*Store, error) {
	s := Store{root: root}
	path := filepath.Join(root, "index.db")
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	f.Close()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("could open db: %w", err)
	}
	db.SetMaxOpenConns(1)
	s.DB = db

	return &s, nil
}

type Store struct {
	root string // path for both sqlite db and objects
	*sql.DB
}

func (s *Store) Update(id string, r io.ReadSeeker) error {
	return s.createOrUpdate(r, id)
}

func (s *Store) Create(r io.ReadSeeker) (string, error) {
	id := randID()
	return id, s.createOrUpdate(r, id)
}

func (s *Store) createOrUpdate(r io.ReadSeeker, id string) error {
	objectPath := getObjectPath(s.root, id)
	if err := os.MkdirAll(filepath.Dir(objectPath), os.ModePerm); err != nil {
		return fmt.Errorf("create object dir: %w", err)
	}

	f, err := os.OpenFile(objectPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error creating object file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("error writing object: %w", err)
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to begin: %w", err)
	}
	contentType, err := fileContentType(r)
	if err != nil {
		return err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to begin: %w", err)
	}

	tx, err := s.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	now := time.Now().Unix()
	_, err = tx.Exec(`
insert into files (id, created_at, updated_at, content_type) values ($1, $2, $2, $3)
on conflict (id) do
update set updated_at=$2, content_type=$3 where id=$1
`, id, now, contentType)
	if err != nil {
		return fmt.Errorf("inserting into table: %w", err)
	}
	if err := tx.Commit(); err != nil {
		// todo: is it necessary to delete the file here?
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *Store) Search(params Query) ([]File, error) {
	stmt := `select id, created_at, updated_at, content_type from files order by created_at`
	var files []File
	rows, err := s.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("searching files: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		f := File{}
		var createdAt, updatedAt int64
		err := rows.Scan(&f.ID, &createdAt, &updatedAt, &f.Type)
		if err != nil {
			return nil, fmt.Errorf("scanning file: %w", err)
		}
		f.CreatedAt = time.Unix(createdAt, 0)
		f.UpdatedAt = time.Unix(updatedAt, 0)
		files = append(files, f)
	}
	return files, nil
}

type Query struct {
	Page     int
	PageSize int
	FromDate *time.Time
	ToDate   *time.Time
	Tags     []string
}

func (s *Store) Get(id string) (File, error) {
	stmt := `SELECT created_at, updated_at, content_type from files where id=?`
	var contentType string
	var createdAt, updatedAt int64
	err := s.QueryRow(stmt, id).Scan(&createdAt, &updatedAt, &contentType)
	if err != nil {
		return File{}, fmt.Errorf("could query row: %w", err)
	}
	f := File{
		ID:        id,
		CreatedAt: time.Unix(createdAt, 0),
		UpdatedAt: time.Unix(updatedAt, 0),
		Type:      contentType,
	}
	// len of text/plain==10
	if f.Type == "text/plain" {
		b, err := os.Open(getObjectPath(s.root, id))
		if err != nil {
			return File{}, err
		}
		defer b.Close()
		raw, err := ioutil.ReadAll(b)
		if err != nil {
			return File{}, err
		}
		f.Content = string(raw)
	}
	return f, nil
}

func fileContentType(in io.Reader) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	raw, err := ioutil.ReadAll(&(io.LimitedReader{R: in, N: 512}))
	if err != nil {
		return "", err
	}
	fileType, _, err := mime.ParseMediaType(http.DetectContentType(raw))
	if err != nil {
		return "", err
	}
	return fileType, nil
}

func getObjectPath(rootDir, name string) string {
	return filepath.Join(rootDir, "objects", name)
}

type File struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Type      string    `json:"type"`
	// Content is provided only when the data was a text file
	Content string `json:"content"`
}

func randID() string {
	return RandString(6)
}
