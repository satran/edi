package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		Name string
		In   string
		Out  string
	}{
		{
			Name: "echo",
			In:   "echo 'hello'",
			Out:  "hello",
		},
	}
	for _, test := range tests {
		got := strings.TrimSuffix(run(".", test.Name, test.In), "\n")
		if got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}

func TestRunStdin(t *testing.T) {
	tests := []struct {
		Name string
		In   string
		Out  string
	}{
		{
			Name: "multi line posix",
			In: `name="world"
echo $name | awk '{print "hello " $0}'
`,
			Out: "hello world",
		},
	}
	for _, test := range tests {
		got := strings.TrimSuffix(runstdin(".", test.Name, []byte(test.In)), "\n")
		if got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}
