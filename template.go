package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	root string
	fns  template.FuncMap
	*template.Template
}

func NewTemplate(dir string) *Template {
	t := Template{root: dir}
	t.fns = template.FuncMap{
		"title": strings.Title,
		"link":  t.Link,
		"sh":    t.Shell,
	}
	t.Template = template.New("engine").Funcs(t.fns).Delims("((", "))")
	return &t
}

func (t *Template) Link(url string, args ...string) string {
	if len(args) >= 1 {
		return fmt.Sprintf(`<a href=%q>%s</a>`, url, args[0])
	}
	return fmt.Sprintf(`<a href=%q>%s</a>`, url, url)
}

func (t *Template) Shell(args string) string {
	objDir := filepath.Join(t.root, "objects")
	cmd := exec.Command("bash", "-c", args)
	// todo: this is a simple hack to ensure the scripts in the
	// object directory is in the PATH.
	path := os.Getenv("PATH")
	path += ":" + objDir
	os.Setenv("PATH", path)
	cmd.Env = append(os.Environ())
	cmd.Dir = objDir
	// do nothing as it just shows error code
	//return err.Error()
	out, _ := cmd.CombinedOutput()
	return string(out)
}
