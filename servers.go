package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type SocketServer struct {
	sessions       map[int]*Session
	upgrader       websocket.Upgrader
	sessionCounter int
}

// NewSocketServer create a new WS session
func NewSocketServer() *SocketServer {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s := SocketServer{
		sessions:       make(map[int]*Session),
		sessionCounter: 0,
		upgrader:       upgrader,
	}
	return &s
}

func (s *SocketServer) Init(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := getSessionID(r)
	session := s.GetOrCreateSession(id)

	if session == nil {
		log.Println("Session could not be found.")
		return
	}

	conn := NewConn(ws)
	session.AddConn(conn)
	session.Listen()
	conn.Listen(session)
}

func getSessionID(r *http.Request) int {
	var id int
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		return id
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Println("Session ids must be integers")
	}
	return id
}

// GetOrCreateSession fetches the Session given the id or creates a new Sesison
// when id is 0.
func (s *SocketServer) GetOrCreateSession(id int) *Session {
	if id == 0 {
		return s.NewSession()
	} else {
		return s.GetSession(id)
	}
}

func (s *SocketServer) NewSession() *Session {
	session := NewSession()
	s.sessionCounter++
	s.sessions[s.sessionCounter] = session
	return session
}

func (s *SocketServer) GetSession(id int) *Session {
	return nil
}
