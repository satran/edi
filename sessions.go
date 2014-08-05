package main

import (
	"encoding/json"
	"log"
)

type Session struct {
	connections    []*Conn        // List of all the active connections for this session
	buffers        map[string]int // Hash of all buffers open: file_name:buffer_id
	current        string         // Abs file name of the current buffer
	processes      map[string]int // Process ids for any command executed in the session
	cwd            string         // the current working directory of the session
	bufferCounter  int            // counter to auto generate ids for buffers
	commandCounter int            // counter to auto generate command ids from server

	input  chan []byte      //Channel that listens to inputs from various clients
	output chan interface{} // a channel where commands push result into
}

func NewSession() *Session {
	s := Session{}
	s.connections = []*Conn{}
	s.buffers = make(map[string]int)
	s.processes = make(map[string]int)

	s.input = make(chan []byte)
	s.output = make(chan interface{})

	return &s
}

func (s *Session) Listen() {
	go s.listenInput()
	go s.listenOutput()
}

func (s *Session) listenInput() {
	for {
		message := <-s.input
		cmd := Command{}
		if err := json.Unmarshal(message, &cmd); err != nil {
			log.Println("Can't decode", err)
			continue
		}
		go cmd.Exec(s)
	}
}

func (s *Session) listenOutput() {
	for {
		response := <-s.output
		for _, conn := range s.connections {
			conn.Write(response)
		}
	}
}

func (s *Session) AddConn(c *Conn) {
	s.connections = append(s.connections, c)
}
