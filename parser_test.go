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
