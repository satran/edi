package main

import (
	"bytes"
	"regexp"

	bf "github.com/russross/blackfriday/v2"
)

type Parser struct {
	root     string
	filename string
}

func NewParser(dir string, name string) *Parser {
	return &Parser{root: dir, filename: name}
}

func (p *Parser) Parse(content string) string {
	content = parseInternalLinks(content)
	m := bf.New(bf.WithExtensions(bf.CommonExtensions |
		bf.HeadingIDs |
		bf.Footnotes |
		bf.NoEmptyLineBeforeBlock |
		bf.AutoHeadingIDs))
	ast := m.Parse([]byte(content))
	var buf bytes.Buffer
	renderer := bf.NewHTMLRenderer(bf.HTMLRendererParameters{})
	listDepth := 0
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
		case bf.List:
			if entering {
				listDepth++
			} else {
				listDepth--
			}
		case bf.Text:
			if listDepth == 0 {
				break
			}
			if isTask(node.Literal) {
				matched = true
				buf.Write(toTask(node.Literal, []byte(`<span class="task">$1</span>&nbsp;`)))
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

var (
	taskR  = regexp.MustCompile(`^\[([ a-zA-Z\?/]*)\] `)
	isTask = taskR.Match
	toTask = taskR.ReplaceAll
)

var linkR = regexp.MustCompile(`\(\((.*)\)\)`)

func parseInternalLinks(content string) string {
	return linkR.ReplaceAllString(content, `<a href="$1">$1</a>`)
}
