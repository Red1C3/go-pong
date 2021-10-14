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
	defer func() {
		gameLobby.disconnected <- c
	}()
	for {
		msgType, p, err := c.connection.ReadMessage()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				if ce.Code == websocket.CloseNormalClosure {
					return
				}
			}
			log.Fatalf("error while reading msg from client %v : %v", c.ID, err)
		}
		if msgType == websocket.BinaryMessage && len(p) == 1 {
			var dir int
			if p[0] == 0 {
				dir = -1
			} else {
				dir = 1
			}
			timeMutex.RLock()
			playersMutex.Lock()
			players[c.ID-1].Move(float64(dir)*30*10000, deltaTime)
			log.Print(deltaTime)
			playersMutex.Unlock()
			timeMutex.RUnlock()
		}
	}
}
