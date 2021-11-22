package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
)

// ReadWriteSeekCloser is the interface that groups the basic Read, Write, Seek and Close methods.
type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

type File struct {
	ReadWriteSeekCloser
	parser *Parser
	Name   string
	Type   string
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
	return template.HTML(f.parser.Parse(f.Content()))
}

func (f *File) Content() string {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		// todo: return a better error. For convenience this works
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

func Dummy(root string, name string) *File {
	return &File{
		ReadWriteSeekCloser: &Buffer{&bytes.Buffer{}},
		Name:                name,
		Type:                "text/plain",
		parser:              NewParser(root, name),
	}
}

type Buffer struct {
	*bytes.Buffer
}

func (b *Buffer) Close() error {
	// nothing here just a holder to help implement the Closer
	return nil
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
