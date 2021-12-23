package parser

import (
	"errors"
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
	var args []arg
	for i := range c {
		switch i.typ {
		case itemText:
			ret += i.val
		case itemLeftMeta:
			args = []arg{}
		case itemArg:
			args = append(args, arg{value: i.val})
		case itemArgMultiLine:
			args = append(args, arg{value: i.val, multiline: true})
		case itemArgQuoted:
			args = append(args, arg{value: i.val, quoted: true})
		case itemRightMeta:
			parsed, err := p.eval(args)
			if err != nil {
				ret += "\n" + err.Error()
				return ret
			}
			ret += parsed
		case itemHeader:
			ret += fmt.Sprintf("<strong>%s</strong>", i.val)
		}
	}
	return ret
}

type arg struct {
	value     string
	multiline bool
	quoted    bool
}

func (p *Parser) eval(args []arg) (string, error) {
	if len(args) == 0 {
		return "", errors.New("eval block must contain arguments")
	}
	if args[0].multiline {
		return "", errors.New("function cannot be multiline")
	}
	fn := args[0].value
	switch {
	case fn == "!" || fn == "sh":
		content, err := shell(p.root, p.filename, args[1:])
		if err != nil {
			return content, err
		}
		nested := New(p.root, p.filename)
		return nested.Parse(content), nil
	case fn == "i" || fn == "image":
		if len(args) == 1 {
			return "", errors.New(`use ((img url "optional alt text"))`)
		}
		return image(args[1:])
	case len(args) == 1 && args[0].quoted:
		return link(args[0].value)
	}
	return "", fmt.Errorf("unknown function: %s", fn)
}

func shell(root string, filename string, args []arg) (string, error) {
	if len(args) == 1 && args[0].multiline {
		return exec.RunStdin(root, filename, []byte(args[0].value)), nil
	}
	params := ""
	for _, a := range args {
		if a.multiline {
			return "", errors.New("one one multiline argument allowed")
		}
		if a.quoted {
			params += fmt.Sprintf(" %q", a.value)
		} else {
			params += " " + a.value
		}

	}
	return exec.Run(root, filename, params), nil
}

// link takes a single quoted argument.
// A link can be represented as url|optional string
func link(url string) (string, error) {
	var link, text string
	chunks := strings.SplitN(url, "|", 2)
	if len(chunks) == 1 {
		link = chunks[0]
		text = chunks[0]
	} else if len(chunks) == 2 {
		link = chunks[0]
		text = chunks[1]
	} else {
		return "", errors.New(`use (("url[|optional text]"))`)
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, link, text), nil
}

func image(args []arg) (string, error) {
	var url, alt string
	if len(args) == 1 {
		url = args[0].value
		alt = url
	} else if len(args) == 2 {
		url = args[0].value
		alt = args[1].value
	} else {
		return "", errors.New(`use ((img url "optional alt text"))`)
	}
	return fmt.Sprintf(`<img src="%s" alt="%s"/>`, url, alt), nil
}
