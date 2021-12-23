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
	itemLeftMeta
	itemRightMeta
	itemArgMultiLine
	itemArgQuoted
	itemArg
	itemText
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
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta // Next state.
		}
		if r := l.next(); r == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction // Now inside (( ))
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(rightMeta)
	l.emit(itemRightMeta)
	return lexText
}

func lexMultiLine(l *lexer) stateFn {
	l.pos += len(multiLineQuote)
	l.ignore()
	for {
		if strings.HasPrefix(l.input[l.pos:], multiLineQuote) {
			l.emit(itemArgMultiLine)
			l.pos += len(multiLineQuote)
			l.ignore()
			return lexInsideAction
		}
		l.next()
	}
}

func lexInsideQuote(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '"':
			l.backup()
			l.emit(itemArgQuoted)
			l.pos += 1 // skip one for the backup
			l.ignore()
			return lexInsideAction
		}
	}
}

func lexInsideAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate and are ignored.
	// Pipe symbols separate and are emitted.
	nested := 0
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) && nested == 0 {
			if l.pos > l.start {
				l.emit(itemArg)
			}
			return lexRightMeta
		}
		if strings.HasPrefix(l.input[l.pos:], multiLineQuote) {
			if l.pos > l.start {
				l.emit(itemArg)
			}
			return lexMultiLine
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case r == '"':
			l.ignore()
			return lexInsideQuote
		case r == ' ' && nested == 0:
			l.backup()
			l.emit(itemArg)
			l.next()
			l.ignore()
		case r == '(':
			nested++
		case r == ')':
			if nested > 0 {
				nested--
			}
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
	eof            = 0
	leftMeta       = "(("
	rightMeta      = "))"
	quote          = `"`
	multiLineQuote = `"""`
)
