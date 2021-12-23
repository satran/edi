package parser

import (
	"fmt"
	"strings"

	"github.com/satran/edi/exec"
)

type Parser struct {
	root     string
	filename string
}

func New(dir string, name string) *Parser {
	return &Parser{root: dir, filename: name}
}

func (p *Parser) Parse(content string) string {
	ret := ""
	_, c := lex(p.filename, content)
	for i := range c {
		switch i.typ {
		case itemText:
			ret += i.val
		case itemHeader:
			ret += fmt.Sprintf("<strong>%s</strong>", i.val)
		case itemCode:
			nested := New(p.root, p.filename)
			ret += nested.Parse(exec.Run(p.root, p.filename, i.val))
		case itemCodeMultiLine:
			nested := New(p.root, p.filename)
			ret += nested.Parse(exec.RunStdin(p.root, p.filename, []byte(i.val)))
		case itemLink:
			var url, opt string
			chunks := strings.SplitN(i.val, "|", 2)
			url = chunks[0]
			opt = chunks[0]
			if len(chunks) == 2 {
				opt = chunks[1]
			}
			if isImageLink(url) {
				url = url[1:]
				if len(chunks) == 1 {
					ret += fmt.Sprintf(`<img src="%s" />`, url)
				} else {
					ret += fmt.Sprintf(`<img src="%s" alt="%s" />`, url, opt)
				}

			} else {
				ret += fmt.Sprintf(`<a href="%s">%s</a>`, url, opt)
			}
		}
	}
	return ret
}

func isImageLink(link string) bool {
	return strings.HasPrefix(link, "!")
}
