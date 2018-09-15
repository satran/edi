package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/satran/edi/buffer"
)

func newEditor(t Terminal, names ...string) (*editor, error) {
	buffers := make([]buffer.Buffer, 0, len(names))
	for _, name := range names {
		b, err := buffer.New(name)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %s", name, err)
		}
		buffers = append(buffers, b)
	}
	// todo: create buffer when names are empty
	return &editor{terminal: t, buffers: buffers, current: buffers[0]}, nil
}

type editor struct {
	terminal Terminal
	current  buffer.Buffer
	buffers  []buffer.Buffer
}

// ListenAndServe renders the UI and starts the FUSE filesystem
func (e *editor) ListenAndServe() error {
	f, err := os.Open(e.current.Name())
	if err != nil {
		return err
	}
	defer f.Close()
	if err := setStatus(e.terminal, e.current.Name()); err != nil {
		return err
	}
	if err := render(e.terminal, f); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	return nil
}

// render draws the main frame for the editor
func render(t Terminal, r io.Reader) error {
	br := bufio.NewReader(r)
	w, h := termbox.Size()
	// the first line is always for the status
	for x, y := 0, 1; ; {
		rn, _, err := br.ReadRune()
		if err != nil && err != io.EOF {
			return err
		}
		// always wrap the long lines
		if x > w {
			y++
			x = 0
		}
		if rn == '\n' {
			y++
			x = 0
			continue
		}
		// we have run out of lines
		if y > h {
			break
		}

		if rn == '\t' {
			for i := 0; i < 8; i, x = i+1, x+1 {
				t.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}
			continue
		}

		t.SetCell(x, y, rn, termbox.ColorDefault, termbox.ColorDefault)
		x++

	}
	return t.Flush()
}

func setStatus(t Terminal, status string) error {
	w, _ := termbox.Size()
	var i int
	var rn rune
	// todo: check if ranging of multibyte characters renders this incorrectly
	for i, rn = range status {
		if i > w {
			break
		}
		t.SetCell(i, 0, rn,
			termbox.ColorDefault|termbox.AttrReverse,
			termbox.ColorDefault|termbox.AttrReverse)
	}
	// Ensure the whole line is with an inversed color
	if i < w {
		for i++; i < w; i++ {
			t.SetCell(i, 0, ' ',
				termbox.ColorDefault|termbox.AttrReverse,
				termbox.ColorDefault|termbox.AttrReverse)
		}
	}
	return t.Flush()
}
