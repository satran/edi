package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	_ "net/http/pprof"
)

var cmdServer = &Command{
	UsageLine: "server [-static] [-addr] [-dir] [-basic]",
	Short:     "starts a webserver to serve files",
	Long: `starts a webserver to serve files
  -addr string
    	addr and port to serve from (default "127.0.0.1:8080")
  -basic
    	enable Basic Authentication, requires you to set USERNAME and PASSWORD environment variables
  -dir string
    	directory that stores the data (default "~/dabba")
  -start-file string
    	file to render on index (default "Start")
  -static string
    	Render static files from folder
`,
	Run: runServer,
}
var (
	serverStaticDir string
	serverDocDir    string
	serverAddr      string
	serverBasicAuth bool
	serverStartFile string
)

func init() {
	cmdServer.Flag.StringVar(&serverStaticDir, "static", "", "Render static files from folder")
	cmdServer.Flag.StringVar(&serverDocDir, "dir", dabbaDefaultDir, "directory that stores the data")
	cmdServer.Flag.StringVar(&serverAddr, "addr", "127.0.0.1:8080", "addr and port to serve from")
	cmdServer.Flag.BoolVar(&serverBasicAuth, "basic", false, "enable Basic Authentication, requires you to set USERNAME and PASSWORD environment variables")
	cmdServer.Flag.StringVar(&serverStartFile, "start-file", "Start", "file to render on index")

}

//go:embed s templates
var contents embed.FS

func runServer(cmd *Command, args []string) {
	root, err := filepath.Abs(serverDocDir)
	if err != nil {
		log.Fatal(err)
	}
	s := NewStore(root, serverStartFile)

	tmpls := template.Must(template.ParseFS(contents, "templates/*"))

	srv := &http.Server{Addr: serverAddr}
	if serverStaticDir != "" {
		http.Handle("/s/", http.StripPrefix("/s/",
			http.FileServer(http.Dir(serverStaticDir))))
	} else {
		http.Handle("/s/", http.FileServer(http.FS(contents)))
	}

	if serverBasicAuth {
		username := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")
		if username == "" || password == "" {
			log.Fatal("expected USERNAME and PASSWORD enviornment variables to be set")
		}
		http.Handle("/", basicAuth(Handler(s, tmpls), username, password))
	} else {
		http.Handle("/", Handler(s, tmpls))
	}

	go func() {
		log.Printf("Starting server %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("shutting down")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
