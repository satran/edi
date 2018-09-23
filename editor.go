package main

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// Editor handles FUSE file system and rendering the editor
type Editor struct {
	terminal Terminal
	current  *BufferView
	console  *BufferView
	views    map[string]*BufferView
}

func newEditor(t Terminal, names ...string) (*Editor, error) {
	e := &Editor{
		terminal: &terminal{},
		views:    make(map[string]*BufferView),
	}
	for i, name := range names {
		b, err := NewBuffer(name)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %s", name, err)
		}
		v := newBufferView(t, b)
		e.views[name] = v
		if i == 0 {
			e.current = v
		}
	}
	// todo: create buffer when names are empty
	return e, nil
}

func (e *Editor) Close() {
	for _, v := range e.views {
		v.Close()
	}
}

// ListenAndServe renders the UI and starts the FUSE filesystem
func (e *Editor) ListenAndServe() error {
	_, h := e.terminal.Size()
	if err := e.terminal.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		return err
	}
	err := e.current.Render(0, h)
	if err != nil {
		return err
	}
	e.terminal.Flush()
	time.Sleep(5 * time.Second)
	return nil
}

/*func (e *Editor) handleKeypress() error {
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			if err := renderBuffer(e.terminal, e.current); err != nil {
				return err
			}
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				e.current.line--
				if err := renderBuffer(e.terminal, e.current); err != nil {
					return err
				}
			case termbox.KeyArrowDown:
				e.current.line++
				if err := renderBuffer(e.terminal, e.current); err != nil {
					return err
				}
			case termbox.KeyCtrlQ:
				return nil
			default:
			}

		case termbox.EventInterrupt:
			return nil
		}
	}
}

type Event struct{}
*/
