package main

import (
	"bytes"
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
		"l":     p.Link,
		"i":     p.Image,
		"image": p.Image,
		"sh":    p.Shell,
		"parse": p.Parse,
	}
	p.Template = template.New("engine").Funcs(p.fns).Delims("((", "))")
	return &p
}

func (p *Parser) Parse(content string) string {
	t, err := p.Clone()
	if err != nil {
		return fmt.Sprintf("couldn't load parser: %w", err)
	}
	t, err = t.Parse(content)
	if err != nil {
		return fmt.Sprintf("couldn't parse template: %w", err)
	}
	wr := &bytes.Buffer{}
	if err := t.Execute(wr, nil); err != nil {
		return fmt.Sprintf("couldn't execute template: %w", err)
	}
	return wr.String()

}

func (p *Parser) Image(url string, args ...string) string {
	if len(args) >= 1 {
		return fmt.Sprintf(`<img src=%q alt=%q />`, url, args[0])
	}
	return fmt.Sprintf(`<img src=%q />`, url)
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
