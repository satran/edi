package http

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/satran/edi/defaults"
	"github.com/satran/edi/store"
)

type options struct {
	root     string
	start    string
	static   string
	addr     string
	username string
	password string
}

type Opts func(o *options) *options

func WithRootDir(dir string) Opts {
	return func(o *options) *options {
		if dir != "" {
			o.root = dir
		}
		return o
	}
}

func WithStartFile(filename string) Opts {
	return func(o *options) *options {
		if filename != "" {
			o.start = filename
		}
		return o
	}
}

func WithStaticDir(dir string) Opts {
	return func(o *options) *options {
		if dir != "" {
			o.static = dir
		}
		return o
	}
}

func WithServerAddr(addr string) Opts {
	return func(o *options) *options {
		if addr != "" {
			o.addr = addr
		}
		return o
	}
}

func WithBasicAuth(username string, password string) Opts {
	return func(o *options) *options {
		o.username = username
		o.password = password
		return o
	}
}

//go:embed s templates
var contents embed.FS

// Server returns a http.Server configured to run a webserver
func Server(opts ...Opts) (*http.Server, error) {
	defroot, err := defaults.Root()
	if err != nil {
		return nil, err
	}
	o := &options{
		root:   defroot,
		start:  "Start",
		static: "",
		addr:   "localhost:8080",
	}
	for _, fn := range opts {
		o = fn(o)
	}
	root, err := filepath.Abs(o.root)
	if err != nil {
		log.Fatal(err)
	}
	s := store.New(root, o.start)

	tmpls := template.Must(template.ParseFS(contents, "templates/*"))

	srv := &http.Server{Addr: o.addr}
	if o.static != "" {
		http.Handle("/s/", http.StripPrefix("/s/",
			http.FileServer(http.Dir(o.static))))
	} else {
		http.Handle("/s/", http.FileServer(http.FS(contents)))
	}

	if o.username != "" || o.password != "" {
		http.Handle("/", basicAuth(handler(s, tmpls), o.username, o.password))
	} else {
		http.Handle("/", handler(s, tmpls))
	}

	return srv, nil
}
