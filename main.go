package main

import (
	"net/http"
)

func main() {
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("app"))))
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8312", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "app/edi.html")
}
