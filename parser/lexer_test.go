package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name    string
		content string
		items   []item
	}{
		{
			name:    "simple",
			content: `hello world`,
			items: []item{
				{itemText, "hello world"},
				{itemEOF, ""},
			},
		},
		{
			name: "simple new line",
			content: `
hello world`,
			items: []item{
				{itemText, "\nhello world"},
				{itemEOF, ""},
			},
		},
		{
			name:    "embedded",
			content: `hello world ((arg1 arg2))`,
			items: []item{
				{itemText, "hello world "},
				{itemLeftMeta, leftMeta},
				{itemArg, "arg1"},
				{itemArg, "arg2"},
				{itemRightMeta, rightMeta},
				{itemEOF, ""},
			},
		},
		{
			name: "embedded multi line",
			content: `hello world
((this is """a multi
line argument"""))`,
			items: []item{
				{itemText, "hello world\n"},
				{itemLeftMeta, leftMeta},
				{itemArg, "this"},
				{itemArg, "is"},
				{itemMultiLineArg, "a multi\nline argument"},
				{itemRightMeta, rightMeta},
				{itemEOF, ""},
			},
		},
	}
	for _, test := range tests {
		_, c := lex(test.name, test.content)
		pos := 0
		for i := range c {
			if pos >= len(test.items) {
				t.Errorf("%s: got more items: %#v", test.name, i)
				continue
			}
			if test.items[pos] != i {
				t.Errorf("%s: expected %#v, got %#v at %d", test.name, test.items[pos], i, pos)
			}
			pos++
		}
	}
}
