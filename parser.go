package main

import (
	"bytes"
	"fmt"
	"text/template"

	bf "github.com/russross/blackfriday/v2"
)

type Parser struct {
	root     string
	fns      template.FuncMap
	filename string
}

func NewParser(dir string, name string) *Parser {
	p := Parser{root: dir, filename: name}
	return &p
}

func (p *Parser) Parse(content string) string {
	m := bf.New(bf.WithExtensions(bf.CommonExtensions |
		bf.HeadingIDs |
		bf.Footnotes |
		bf.NoEmptyLineBeforeBlock |
		bf.AutoHeadingIDs))
	ast := m.Parse([]byte(content))
	var buf bytes.Buffer
	renderer := bf.NewHTMLRenderer(bf.HTMLRendererParameters{})
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		var matched bool
		switch node.Type {
		case bf.Code:
			if shouldEval(node.Literal) {
				comm := bytes.TrimPrefix(node.Literal, []byte("!"))
				buf.WriteString(run(p.root, p.filename, string(comm)))
				matched = true
			}
		case bf.CodeBlock:
			// fenced code is only the one with ```<type>```. We want the type to be !
			if !node.CodeBlockData.IsFenced {
				break
			}
			if shouldEval(node.CodeBlockData.Info) {
				buf.WriteString(runstdin(p.root, p.filename, node.Literal))
				matched = true
			}
		}
		if !matched {
			renderer.RenderNode(&buf, node, entering)
		}
		return bf.GoToNext
	})
	return buf.String()
}

func shouldEval(content []byte) bool {
	return bytes.HasPrefix(content, []byte("!"))
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
	return run(p.root, p.filename, args)
}
