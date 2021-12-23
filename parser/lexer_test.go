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
			name:    "header",
			content: `# hello world`,
			items: []item{
				{itemHeader, "# hello world"},
				{itemEOF, ""},
			},
		},
		{
			name: "header multiline",
			content: `a new line
# hello world
this is a test`,
			items: []item{
				{itemText, "a new line\n"},
				{itemHeader, "# hello world"},
				{itemText, "\nthis is a test"},
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
			content: "hello world `arg1 arg2`",
			items: []item{
				{itemText, "hello world "},
				{itemCode, "arg1 arg2"},
				{itemEOF, ""},
			},
		},
		{
			name:    "embedded quoted arg",
			content: "hello world `arg1 arg2 \"arg3 value\"`",
			items: []item{
				{itemText, "hello world "},
				{itemCode, "arg1 arg2 \"arg3 value\""},
				{itemEOF, ""},
			},
		},
		{
			name:    "embedded multi line",
			content: "hello world\nthis is\n```a multi\nline argument```",
			items: []item{
				{itemText, "hello world\nthis is\n"},
				{itemCodeMultiLine, "a multi\nline argument"},
				{itemEOF, ""},
			},
		},
		{
			name:    "embedded nested",
			content: "```l=`date````",
			items: []item{
				{itemCodeMultiLine, "l=`date`"},
				{itemEOF, ""},
			},
		},
		{
			name:    "simple link",
			content: "[[this is a link]]",
			items: []item{
				{itemLink, "this is a link"},
				{itemEOF, ""},
			},
		},
		{
			name:    "link with description",
			content: "[[this is a link|with description]]",
			items: []item{
				{itemLink, "this is a link|with description"},
				{itemEOF, ""},
			},
		},
		{
			name:    "link with text around",
			content: "this is a line [[this is a link]]\and another line",
			items: []item{
				{itemText, "this is a line "},
				{itemLink, "this is a link"},
				{itemText, "\and another line"},
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
