package parser

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
			Out:  `<strong># hello</strong>`,
		},
		{
			Name: "code eval",
			In:   "((! printf hello))",
			Out:  `hello`,
		},
		{
			Name: "code eval with text",
			In:   "((! printf hello)) world",
			Out:  `hello world`,
		},
		{
			Name: "block code eval",
			In: `((! """
printf hello | tr 'a-z' 'A-Z'"""))`,
			Out: `HELLO`,
		},
		{
			Name: "filename-env",
			In:   "((! printf $FILE))",
			Out:  `filename-env`,
		},
	}

	for _, test := range tests {
		p := New(".", test.Name)
		if got := strings.TrimSuffix(p.Parse(test.In), "\n"); got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}
