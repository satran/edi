package buffer

import (
	"errors"
	"io"
)

type Buffer interface {
	io.ReadWriteSeeker
	io.Closer

	Line(n int64) ([]byte, error)
	Name() string
}

type defaultBuffer struct {
	name string
}

func New(name string) (Buffer, error) {
	return &defaultBuffer{name: name}, nil
}

func (b *defaultBuffer) Read(p []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (b *defaultBuffer) Write(p []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (b *defaultBuffer) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("not implemented")
}

func (b *defaultBuffer) Close() error {
	return errors.New("not implemented")
}

func (b *defaultBuffer) Line(n int64) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (b *defaultBuffer) Name() string {
	return b.name
}
