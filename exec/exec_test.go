package exec

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
		got := strings.TrimSuffix(Run(".", test.Name, test.In), "\n")
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
		got := strings.TrimSuffix(RunStdin(".", test.Name, []byte(test.In)), "\n")
		if got != test.Out {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Out, got)
		}
	}
}

func TestAddPath(t *testing.T) {
	tests := []struct {
		Name string
		Path string
		Dir  string
		Exp  string
	}{
		{Name: "blank", Path: "", Dir: "", Exp: ""},
		{Name: "single", Path: "/home/satran", Dir: "/home/satran/bin", Exp: "/home/satran:/home/satran/bin"},
		{Name: "simple", Path: "/home/satran:/home/satran/bin", Dir: "/home/satran/go/bin", Exp: "/home/satran:/home/satran/bin:/home/satran/go/bin"},
		{Name: "simple-exists", Path: "/home/satran:/home/satran/bin", Dir: "/home/satran/bin", Exp: "/home/satran:/home/satran/bin"},
		{Name: "exists", Path: "/home/satran:/usr/local/bin:/home/satran/bin", Dir: "/usr/local/bin", Exp: "/home/satran:/usr/local/bin:/home/satran/bin"},
	}
	for _, test := range tests {
		got := addPath(test.Path, test.Dir)
		if got != test.Exp {
			t.Errorf("failed %s\nexpected: \n%q\ngot: \n%q", test.Name, test.Exp, got)
		}
	}
}
