package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	root   string
	config *Config
	parser *Parser
}

func NewStore(root string) (*Store, error) {
	// Create the file if it doesn't exist
	f, err := os.OpenFile(filepath.Join(root, "config.json"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("can't open settings file: %w", err)
	}
	defer f.Close()
	c := &Config{}
	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("decode settings file: %w", err)
	}
	s := &Store{root: root, config: c}
	s.parser = NewParser(s.objpath())
	return s, nil
}

type Config struct {
	StartFile string `json:"start-file"`
	MenuFile  string `json:"menu-file"`
}

func (s *Store) Get(name string) (*File, error) {
	f, err := os.OpenFile(s.path(name), os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", name, err)
	}
	type_, err := fileContentType(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	file := File{
		ReadWriteSeekCloser: f,
		Name:                name,
		Type:                type_,
		path:                s.path(name),
		parser:              s.parser,
	}
	// To ensure that further reads don't start at the wrong offset
	if err := file.SeekStart(); err != nil {
		file.Close()
		return nil, err
	}
	return &file, err
}

func fileContentType(r io.ReadSeeker) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	raw, err := ioutil.ReadAll(&(io.LimitedReader{R: r, N: 512}))
	if err != nil {
		return "", err
	}
	fileType, _, err := mime.ParseMediaType(http.DetectContentType(raw))
	if err != nil {
		return "", err
	}
	return fileType, nil
}

func (s *Store) Write(name string, r io.Reader) error {
	var mode fs.FileMode = 0600
	if strings.HasSuffix(name, ".sh") {
		mode = 0700
	}
	path := s.path(name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", name, err)
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("write file %q: %w", name, err)
	}
	return nil
}

func (s *Store) Index() string {
	return s.config.StartFile
}

func (s *Store) path(name string) string {
	return filepath.Join(s.root, "objects", name)
}

func (s *Store) objpath() string {
	return filepath.Join(s.root, "objects")
}
