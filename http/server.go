package http

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/satran/edi/store"
)

type options struct {
	root     string
	start    string
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

//go:embed _s templates
var contents embed.FS

// Server returns a http.Server configured to run a webserver
func Server(opts ...Opts) (*http.Server, error) {
	o := &options{
		root:  ".",
		start: "Start",
		addr:  "localhost:8080",
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

	http.Handle("/_s/", http.FileServer(http.FS(contents)))
	http.Handle("/_edit/", basicAuth(editH(s, tmpls, "/_edit/"), o.username, o.password))
	http.Handle("/_new", basicAuth(newH(s, tmpls, "/_new"), o.username, o.password))
	http.Handle("/_add/", basicAuth(addBlobH(s, "/_add/"), o.username, o.password))
	http.Handle("/_blob/", getBlobH(s))
	http.Handle("/_sh/", basicAuth(shellH(s, tmpls), o.username, o.password))
	http.Handle("/", getH(s, tmpls))

	return srv, nil
}
