package file

import (
	"errors"
	"io"
	"os"
)

type Buffer interface {
	io.ReadWriteSeeker
	io.Closer

	Line(n int64) ([]byte, error)
	Name() string
}

type defaultBuffer struct {
	name string
	file *os.File
}

func New(name string) (Buffer, error) {
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

func (b *defaultBuffer) Line(n int64) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (b *defaultBuffer) Name() string {
	return b.name
}
