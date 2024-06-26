package store

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/satran/edi/parser"
)

var fileExt = ".txt"

type Store struct {
	Root      string
	startFile string
}

func New(root string, start string) *Store {
	s := &Store{Root: root, startFile: start}
	return s
}

func (s *Store) List() ([]string, error) {
	files, err := ioutil.ReadDir(s.objpath())
	if err != nil {
		return nil, fmt.Errorf("couldn't list directory: %w", err)
	}
	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

func (s *Store) Get(name string) (*File, error) {
	f, err := os.OpenFile(s.Path(name), os.O_RDONLY, 0600)
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
		parser:              parser.New(s.objpath(), name),
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
	ext := filepath.Ext(name)
	if ext == "" {
		// Assume all extension-less files are the txt files
		name = name + fileExt
	}
	var mode fs.FileMode = 0600
	if ext == "sh" {
		mode = 0700
	}
	path := s.Path(name)
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
	return s.startFile
}

func (s *Store) Path(name string) string {
	return filepath.Join(s.Root, name)
}

func (s *Store) PathTxt(name string) string {
	return filepath.Join(s.Root, name+fileExt) // .txt is used so that editors that don't support text files without an extension also works.
}

func (s *Store) objpath() string {
	return s.Root
}
