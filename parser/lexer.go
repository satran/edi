package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemHeader
	itemText
	itemLink
	itemCode
	itemCodeMultiLine
	itemEOF
)

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	return fmt.Sprintf("%q", i.val)
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// run lexes the input by executing state functions
// until the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {
	for {
		switch {
		case l.hasPrefix(codeMultiLineBlock):
			l.emitWhenNotStart(itemText)
			return lexInsideMultiLineBlock
		case l.hasPrefix(leftLink):
			l.emitWhenNotStart(itemText)
			return lextLink
		}

		// generate header
		if n := l.peek(); n == '#' {
			if l.pos == 0 {
				return lexHeader
			}
			if l.input[l.pos-1] == '\n' {
				l.emit(itemText)
				return lexHeader
			}
		} else if n == '`' {
			if l.pos != 0 {
				l.emit(itemText)
			}
			l.next() // read and ignore the `
			l.ignore()
			return lextInsideCode
		}
		r := l.next()
		if r == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil
}

func (l *lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(l.input[l.pos:], prefix)
}

func (l *lexer) emitWhenNotStart(t itemType) {
	if l.pos > l.start {
		l.emit(t)
	}
}

func lexInsideMultiLineBlock(l *lexer) stateFn {
	l.ignoreN(len(codeMultiLineBlock))
	var nested bool
	for {
		if l.hasPrefix(codeMultiLineBlock) && !nested {
			l.emitWhenNotStart(itemCodeMultiLine)
			l.ignoreN(len(codeMultiLineBlock))
			return lexText
		}
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed action")
		case r == codeBlock:
			nested = !nested
		}
	}
}

func lextInsideCode(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case r == codeBlock:
			l.backup()
			l.emit(itemCode)
			l.next()
			l.ignore()
			return lexText
		}
	}
}

func lextLink(l *lexer) stateFn {
	l.ignoreN(len(leftLink))
	nested := 0
	for {
		if l.hasPrefix(rightLink) && nested == 0 {
			l.emitWhenNotStart(itemLink)
			l.ignoreN(len(rightLink))
			return lexText
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case r == '[':
			nested++
		case r == ']':
			if nested > 0 {
				nested--
			}
		}
	}
}

func lexHeader(l *lexer) stateFn {
	for {
		r := l.next()
		if r == '\n' || r == eof {
			l.backup()
			l.emit(itemHeader)
			return lexText
		}
	}
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) ignoreN(n int) {
	l.pos += n
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

const (
	eof                = 0
	codeBlock          = '`'
	codeMultiLineBlock = "```"
	leftLink           = "[["
	rightLink          = "]]"
)
