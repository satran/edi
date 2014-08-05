package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Conn stores details about the a WebSocket connection
type Conn struct {
	ws *websocket.Conn
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{
		ws: ws,
	}
}

// handleRead reads data from the client. Every command from the client should be a
// JSON object which can be parsed into Command.
func (c *Conn) Listen(session *Session) {
	defer c.ws.Close()
	go c.pingClient()

	for {
		t, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("Error ", err, t)
			return
		}
		switch t {
		default:
			log.Println("Can't handle current type.")
			continue

		case websocket.TextMessage:
			session.input <- message

		case websocket.CloseMessage:
			log.Println("Client closed connection.")
			return
		}
	}
}

func (c *Conn) pingClient() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		<-ticker.C
		if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
	}
}

func (c *Conn) Write(response interface{}) {
	c.ws.WriteJSON(response)
}
