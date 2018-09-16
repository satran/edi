package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/satran/edi/file"
)

type bufferView struct {
	file.Buffer
	line    int
	topline int // line number for the first visible line
	column  int
	offset  int64
	cursorx int
	cursory int
}

func newBufferView(b file.Buffer) *bufferView {
	return &bufferView{
		Buffer: b,
		line:   1,
	}
}
func (b *bufferView) seekLine(n int) (int64, error) {
	offset, err := b.Seek(b.offset, io.SeekStart)
	if err != nil {
		return -1, fmt.Errorf("seek to start: %s", err)
	}
	br := bufio.NewReader(b)
	if n == b.topline {
		return offset, nil
	}
	for {
		r, size, err := br.ReadRune()
		if err != nil {
			return offset, fmt.Errorf("read rune: %s", err)
		}
		b.offset += int64(size)
		if r == '\n' {
			b.topline++
		}
		if b.topline == n {
			_, err = b.Seek(b.offset, io.SeekStart)
			return b.offset, nil
		}
		if b.topline > n {
			return b.offset, fmt.Errorf("can't read line %d", n)
		}
	}
}
