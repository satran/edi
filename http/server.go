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
	root  string
	start string
	addr  string
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
	http.Handle("/_edit/", logRequest(editH(s, tmpls, "/_edit/")))
	http.Handle("/_new", logRequest(newH(s, tmpls, "/_new")))
	http.Handle("/_add/", logRequest(addBlobH(s, "/_add/")))
	http.Handle("/_blob/", logRequest(getBlobH(s)))
	http.Handle("/_sh/", logRequest(shellH(s, tmpls)))
	http.Handle("/_ls/", logRequest(listH(s)))
	http.Handle("/", logRequest(getH(s, tmpls)))

	return srv, nil
}
