package main

import termbox "github.com/nsf/termbox-go"

// Terminal is defined to interface termbox's functions to make it testable
type Terminal interface {
	Clear(fg, bg termbox.Attribute) error
	Flush() error
	HideCursor()
	Interrupt()
	SetCell(x, y int, ch rune, fg, bg termbox.Attribute)
	SetCellDefault(x, y int, ch rune)
	SetCellInverse(x, y int, ch rune)
	SetCursor(x, y int)
	Size() (width int, height int)
	Sync() error
	CellBuffer() []termbox.Cell
	ParseEvent(data []byte) termbox.Event
	PollEvent() termbox.Event
	PollRawEvent(data []byte) termbox.Event
	SetInputMode(mode termbox.InputMode) termbox.InputMode
	SetOutputMode(mode termbox.OutputMode) termbox.OutputMode
}

type terminal struct{}

func (t *terminal) Clear(fg, bg termbox.Attribute) error {
	return termbox.Clear(fg, bg)
}

func (t *terminal) Flush() error {
	return termbox.Flush()
}

func (t *terminal) HideCursor() {
	termbox.HideCursor()
}

func (t *terminal) Interrupt() {
	termbox.Interrupt()
}

func (t *terminal) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
}

func (t *terminal) SetCellDefault(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
}

func (t *terminal) SetCellInverse(x, y int, ch rune) {
	termbox.SetCell(x, y, ch,
		termbox.ColorDefault|termbox.AttrReverse,
		termbox.ColorDefault|termbox.AttrReverse)
}

func (t *terminal) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (t *terminal) Size() (width int, height int) {
	return termbox.Size()
}

func (t *terminal) Sync() error {
	return termbox.Sync()
}

func (t *terminal) CellBuffer() []termbox.Cell {
	return termbox.CellBuffer()
}

func (t *terminal) ParseEvent(data []byte) termbox.Event {
	return termbox.ParseEvent(data)
}

func (t *terminal) PollEvent() termbox.Event {
	return termbox.PollEvent()
}

func (t *terminal) PollRawEvent(data []byte) termbox.Event {
	return termbox.PollRawEvent(data)
}

func (t *terminal) SetInputMode(mode termbox.InputMode) termbox.InputMode {
	return termbox.SetInputMode(mode)
}

func (t *terminal) SetOutputMode(mode termbox.OutputMode) termbox.OutputMode {
	return termbox.SetOutputMode(mode)
}
