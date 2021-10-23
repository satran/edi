package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func (s *Store) Update(id int64, r io.ReadSeeker) error {
	hash, err := writeObject(s.root, r)
	if err != nil {
		return err
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
	_, err = tx.Exec(`insert into log (file_id, object_id, updated_at) values (?, ?, ?)`, id, hash, now)
	if err != nil {
		return fmt.Errorf("inserting to log: %w", err)
	}

	_, err = tx.Exec(`update files set object_id=?, updated_at=?, content_type=? where id=?`,
		hash, now, contentType, id)
	if err != nil {
		return fmt.Errorf("inserting into table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		// no need to delete the file, if the person tries to recreate the file, nothing happens
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *Store) Create(r io.ReadSeeker) (int64, error) {
	hash, err := writeObject(s.root, r)
	if err != nil {
		return 0, err
	}
	contentType, err := fileContentType(r)
	if err != nil {
		return 0, err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, fmt.Errorf("error seeking to begin: %w", err)
	}

	tx, err := s.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	now := time.Now().Unix()
	res, err := tx.Exec(`insert into files (object_id, created_at, updated_at, content_type) values (?, ?, ?, ?)`,
		hash, now, now, contentType)
	if err != nil {
		return 0, fmt.Errorf("inserting into table: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fetching last id: %w", err)
	}
	res, err = tx.Exec(`insert into log (file_id, object_id, updated_at) values (?, ?, ?)`,
		id, hash, now)
	if err != nil {
		return 0, fmt.Errorf("inserting to log: %w", err)
	}
	if err := tx.Commit(); err != nil {
		// no need to delete the file, if the person tries to recreate the file, nothing happens
		return 0, fmt.Errorf("commit: %w", err)
	}
	return id, nil
}

func writeObject(rootDir string, r io.ReadSeeker) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("creating hash: %w", err)
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error seeking to begin: %w", err)
	}

	objectPath := getObjectPath(rootDir, hash)
	if err := os.MkdirAll(filepath.Dir(objectPath), os.ModePerm); err != nil {
		return "", fmt.Errorf("create object dir: %w", err)
	}

	w, err := os.OpenFile(objectPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("error creating object file: %w", err)
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return "", fmt.Errorf("error writing object: %w", err)
	}
	return hash, nil
}

func (s *Store) Get(id int64) (File, error) {
	stmt := `SELECT object_id, created_at, updated_at, content_type from files where id=?`
	var objectID, contentType string
	var createdAt, updatedAt int64
	err := s.QueryRow(stmt, id).Scan(&objectID, &createdAt, &updatedAt, &contentType)
	if err != nil {
		return File{}, fmt.Errorf("could query row: %w", err)
	}
	f := File{
		ID:        id,
		ObjectID:  objectID,
		CreatedAt: time.Unix(createdAt, 0),
		UpdatedAt: time.Unix(updatedAt, 0),
		Type:      contentType,
	}
	b, err := os.Open(getObjectPath(s.root, objectID))
	if err != nil {
		return File{}, err
	}
	defer b.Close()
	// len of text/plain==10
	if f.Type[:10] == "text/plain" {
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
	return http.DetectContentType(raw), nil
}

func getObjectPath(rootDir, hash string) string {
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(rootDir, "objects", dir, file)
}

type File struct {
	ID        int64     `json:"id"`
	ObjectID  string    `json:"object_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Type      string    `json:"type"`
	// Content holds either the content of the file if text or the link to the file if it is an image
	Content string `json:"content"`
}
