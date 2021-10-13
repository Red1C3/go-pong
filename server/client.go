package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	ID         int
	connection *websocket.Conn
}

func (c *client) start() {
	for {
		msgType, p, err := c.connection.ReadMessage()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				if ce.Code == websocket.CloseNormalClosure {
					gameLobby.disconnected <- c
					return
				}
			}
			log.Fatalf("error while reading msg from client %v : %v", c.ID, err)
		}
		_ = p
		_ = msgType
	}
}
