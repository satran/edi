package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", initWS)
	http.HandleFunc("/", IndexHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

var indexTmpl *template.Template

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: ensure this is not loaded everytime the handler is called. It is great for developing.
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("executing index template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

var upgrader = websocket.Upgrader{}

func initWS(w http.ResponseWriter, r *http.Request) {
	log.Println("upgrader")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %s", err)
		return
	}
	defer c.Close()

	in := make(chan *Message)

	go write(c, in)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			return
		}
		msg := Message{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal: ", err)
			return
		}
		msg.mt = mt
		msg.raw = message
		in <- &msg
	}
}

func write(c *websocket.Conn, in chan *Message) {
	for {
		msg := <-in
		log.Println(msg.mt, string(msg.raw))
		go run(c, msg)
	}
}

func run(c *websocket.Conn, msg *Message) {
	args := strings.Split(msg.Cmd, " ")
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("stdout: ", err)
		writeMsg(c, msg.ID, msg.mt, err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		writeMsg(c, msg.ID, msg.mt, err.Error())
		return
	}
	raw := make([]byte, 1024)
	for {
		n, err := stdout.Read(raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			writeMsg(c, msg.ID, msg.mt, err.Error())
			return
		}
		writeMsg(c, msg.ID, msg.mt, string(raw[:n]))
	}

	if err := cmd.Wait(); err != nil {
		writeMsg(c, msg.ID, msg.mt, err.Error())
		return
	}
}

func writeMsg(c *websocket.Conn, id int, msgType int, msg string) {
	out := Output{
		ID:     id,
		Output: msg,
	}
	raw, err := json.Marshal(&out)
	if err != nil {
		log.Println("marshal: ", err)
		return
	}
	if err := c.WriteMessage(msgType, raw); err != nil {
		log.Println("write: ", err)
		return
	}
}

type Message struct {
	mt  int
	raw []byte
	Cmd string `json:"cmd"`
	ID  int    `json:"id"`
}

type Output struct {
	ID     int    `json:"id"`
	Output string `json:"output"`
}
