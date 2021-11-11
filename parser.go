package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Parser struct {
	root string
	fns  template.FuncMap
	*template.Template
}

func NewParser(dir string) *Parser {
	p := Parser{root: dir}
	p.fns = template.FuncMap{
		"title": strings.Title,
		"link":  p.Link,
		"sh":    p.Shell,
	}
	p.Template = template.New("engine").Funcs(p.fns).Delims("((", "))")
	return &p
}

func (p *Parser) Link(url string, args ...string) string {
	if len(args) >= 1 {
		return fmt.Sprintf(`<a href=%q>%s</a>`, url, args[0])
	}
	return fmt.Sprintf(`<a href=%q>%s</a>`, url, url)
}

func (p *Parser) Shell(args string) string {
	cmd := exec.Command("bash", "-c", args)
	// todo: this is a simple hack to ensure the scripts in the
	// object directory is in the PATH.
	path := os.Getenv("PATH")
	path += ":" + p.root
	os.Setenv("PATH", path)
	cmd.Env = append(os.Environ())
	cmd.Dir = p.root
	// do nothing as it just shows error code
	//return err.Error()
	out, _ := cmd.CombinedOutput()
	return string(out)
}
