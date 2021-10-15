package server

import (
	"github.com/gorilla/websocket"
)

type client struct {
	ID         int
	connection *websocket.Conn
}

func (c *client) start() {
	defer func() {
		gameLobby.disconnected <- c
	}()
	for {
		msgType, p, err := c.connection.ReadMessage()
		if err != nil {
			return
		}
		if msgType == websocket.BinaryMessage && len(p) == 1 {
			var dir int
			if p[0] == 0 {
				dir = -1
			} else {
				dir = 1
			}
			playersMutex[c.ID-1].Lock()
			players[c.ID-1].Move(float64(dir)*30, deltaTime)
			playersMutex[c.ID-1].Unlock()
		}
	}
}
