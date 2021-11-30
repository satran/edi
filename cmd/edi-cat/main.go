package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/satran/edi/defaults"
	"github.com/satran/edi/parser"
)

func main() {
	if len(os.Args) <= 1 {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		catFile("stdin", string(raw))
		return
	}
	for _, arg := range os.Args[1:] {
		raw, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		catFile("stdin", string(raw))
	}
}

func catFile(filename string, content string) {
	dir := os.Getenv("EDI_DIR")
	if dir == "" {
		d, err := defaults.Root()
		if err != nil {
			fmt.Println("can't read home directory:", err.Error())
			return
		}
		dir = d
	}
	p := parser.New(dir, filename)
	fmt.Print(string(p.Parse(content)))
}
