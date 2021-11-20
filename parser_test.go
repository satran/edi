package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Name string
		In   string
		Out  string
	}{
		{
			Name: "empty",
			In:   ``,
			Out:  ``,
		},
		{
			Name: "heading",
			In:   `# hello`,
			Out:  `<h1 id="hello">hello</h1>`,
		},
		{
			Name: "code eval",
			In:   "`!printf hello`",
			Out:  `<p>hello</p>`,
		},
		{
			Name: "block code eval",
			In:   "```!\n printf hello | tr 'a-z' 'A-Z'```",
			Out:  `<p>HELLO</p>`,
		},
		{
			Name: "filename-env",
			In:   "`!printf $FILE`",
			Out:  `<p>filename-env</p>`,
		},
		{
			Name: "code no eval",
			In:   "`printf hello`",
			Out:  `<p><code>printf hello</code></p>`,
		},
		{
			Name: "block no code eval",
			In:   "```\n printf hello | tr 'a-z' 'A-Z'```",
			Out:  "<p><code>\n printf hello | tr 'a-z' 'A-Z'</code></p>",
		},
	}

	for _, test := range tests {
		p := NewParser(".", test.Name)
		if got := strings.TrimSuffix(p.Parse(test.In), "\n"); got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}

func TestParseInternalLinks(t *testing.T) {
	tests := []struct {
		Name string
		In   string
		Out  string
	}{
		{
			Name: "empty",
			In:   ``,
			Out:  ``,
		},
		{
			Name: "normal",
			In:   `hello world, this is a ((link)) to a file`,
			Out:  `hello world, this is a <a href="link">link</a> to a file`,
		},
		{
			Name: "nested no parse",
			In:   `Would this (nested (file (work)))`,
			Out:  `Would this (nested (file (work)))`,
		},
		{
			Name: "nested",
			In:   `this should create a (nested (link \((file)))`,
			Out:  `this should create a (nested (link \<a href="file)">file)</a>`,
		},
		{
			Name: "special characters",
			In:   `((hello-world with space %?))`,
			Out:  `<a href="hello-world with space %?">hello-world with space %?</a>`,
		},
	}
	for _, test := range tests {
		if got := strings.TrimSuffix(parseInternalLinks(test.In), "\n"); got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}
