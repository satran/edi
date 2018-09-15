// edi is a very minimal text editor. It provides basic functionality
// for text editing. It allows you to extend functionality through the
// environment that it runs. Executing shell commands is first
// class. It also provides a fuse file system which provides all
// functionality that you can do with the editor.
//
// All settings are done by setting environment variables. All
// environment variables are prefixed with edi_. They can also be set
// by flags, which are without the prefix. These are the following
// variables you can set:
//
// edi_fuse_mount - where to mount the fuse file system. If the
// directory does not exist it will create it. Defaults to /tmp/edi
// edi_tab_size - size of tabs, defaults to 8
package main

import (
	"flag"
	"log"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	debug := flag.Bool("debug", false, "print debug statements")
	flag.Parse()
	if *debug {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	e, err := newEditor(&terminal{}, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	if err := e.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
