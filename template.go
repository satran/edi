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
	fns template.FuncMap
	*template.Template
}

func NewTemplate(dir string) *Template {
	t := Template{root: dir,}
	t.fns = template.FuncMap{
		"title": strings.Title,
		"link": t.Link,
		"sh": t.Shell,
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
	cmd := exec.Command("bash", "-c", args)
	path := os.Getenv("PATH")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PATH=%s", path),
	)
	cmd.Dir = filepath.Join(t.root, "objects")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
