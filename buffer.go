package main

import (
	"io"
	"os"
)

type Buffer interface {
	io.ReadWriteSeeker
	io.Closer

	Name() string
}

type defaultBuffer struct {
	name string
	file *os.File
}

func NewBuffer(name string) (Buffer, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &defaultBuffer{name: name, file: f}, nil
}

func (b *defaultBuffer) Read(p []byte) (int, error) {
	return b.file.Read(p)
}

func (b *defaultBuffer) Write(p []byte) (int, error) {
	return b.file.Write(p)
}

func (b *defaultBuffer) Seek(offset int64, whence int) (int64, error) {
	return b.file.Seek(offset, whence)
}

func (b *defaultBuffer) Close() error {
	return b.file.Close()
}

func (b *defaultBuffer) Name() string {
	return b.name
}
