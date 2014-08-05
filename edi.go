// All IDEs suck. This one sucks less.
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	server := NewSocketServer()

	r.HandleFunc("/{path:.+}", fileHandler)
	r.HandleFunc("/", rootHandler)

	http.Handle("/", r)
	http.HandleFunc("/ws", server.Init)
	http.HandleFunc("/ws/{id:[0-9]*}", server.Init)

	log.Println("Listening on port 8312")
	log.Fatal(http.ListenAndServe(":8312", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeFile(w, "app/edi.html")
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["path"]
	writeFile(w, "app/"+file)
}

func writeFile(w http.ResponseWriter, file string) {
	data, err := Asset(file)
	if err != nil {
		log.Println(err)
	}
	w.Write(data)
}
