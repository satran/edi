package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
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
		ReadWriteSeeker: f,
		Name:            name,
		Type:            type_,
		path:            s.path(name),
		close:           f.Close,
		parser:          s.parser,
	}
	if err := file.SeekStart(); err != nil {
		file.Close()
		return nil, err
	}

	return &file, err
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

type File struct {
	io.ReadWriteSeeker
	parser *Parser
	Name   string
	Type   string

	path  string
	close func() error
}

func (f *File) Close() error {
	return f.close()
}

func (f *File) IsText() bool {
	return strings.HasPrefix(f.Type, "text/plain")
}

var imageMime = map[string]bool{
	"image/avif":    true,
	"image/gif":     true,
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/svg+xml": true,
	"image/webp":    true,
}

func (f *File) IsImage() bool {
	return imageMime[f.Type]
}

func (f *File) Parse() template.HTML {
	t, err := f.parser.Clone()
	if err != nil {
		return template.HTML(fmt.Sprintf("couldn't load parser %q: %w", f.Name, err))
	}
	t, err = t.Parse(f.Content())
	if err != nil {
		return template.HTML(fmt.Sprintf("couldn't parse template %q: %w", f.Name, err))
	}
	wr := &bytes.Buffer{}
	if err := t.Execute(wr, nil); err != nil {
		return template.HTML(fmt.Sprintf("couldn't execute template %q: %w", f.Name, err))
	}
	return template.HTML(wr.String())
}

func (f *File) Content() string {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		// for now return the parser error as content
		return err.Error()
	}
	return string(content)
}

func (f *File) SeekStart() error {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to begin: %w", err)
	}
	return nil
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
