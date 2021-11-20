package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var cmdCat = &Command{
	UsageLine: "cat [filename]",
	Short:     "parses the content of the file or stdin",
	Long:      `parses the content of the file or stdin `,
	Run:       runCat,
}

func runCat(cmd *Command, args []string) {
	if len(args) == 0 {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		catFile("stdin", string(raw))
		return
	}
	for _, arg := range args {
		raw, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		catFile("stdin", string(raw))
	}
}

func catFile(filename string, content string) {
	dir := os.Getenv("DABBA_DIR")
	if dir == "" {
		dir = dabbaDefaultDir
	}
	p := NewParser(dir, filename)
	fmt.Print(string(p.Parse(content)))
}
