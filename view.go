package main

import (
	"bufio"
	"fmt"
	"io"

	termbox "github.com/nsf/termbox-go"
)

// BufferView is the view for a buffer. It handles how the buffer must be rendered
type BufferView struct {
	terminal     Terminal
	buffer       Buffer
	cursorLine   int
	cursorColumn int
	topLine      int
}

func newBufferView(t Terminal, buffer Buffer) *BufferView {
	return &BufferView{
		terminal:     t,
		buffer:       buffer,
		cursorLine:   1,
		cursorColumn: 1,
		topLine:      1,
	}
}

func (v *BufferView) Close() error {
	return v.buffer.Close()
}

func (v *BufferView) Render(start, stop int) error {
	v.setStatusLine(start)
	v.setContents(start+1, stop)
	return v.terminal.Flush()
}

func (v *BufferView) setStatusLine(x int) {
	w, _ := v.terminal.Size()
	var i int
	var rn rune
	status := fmt.Sprintf("%s %d:%d", v.buffer.Name, v.cursorLine, v.cursorColumn)
	// todo: check if ranging of multibyte characters renders this incorrectly
	for i, rn = range status {
		if i > w {
			break
		}
		v.terminal.SetCellInverse(i, 0, rn)
	}
	// Ensure the whole line is with an inversed color
	if i < w {
		for i++; i < w; i++ {
			v.terminal.SetCellInverse(i, 0, ' ')
		}
	}
}

func (v *BufferView) setContents(start, stop int) error {
	w, _ := termbox.Size()
	// we just need as many bytes as the screen can render
	size := stop - start*w
	content := make([]byte, 0, size)
	n, err := v.buffer.Read(content)
	if err != nil && err != io.EOF {
		return err
	}

	r := bufio.NewReader(content)
}
