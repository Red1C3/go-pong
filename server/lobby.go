package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type lobby struct {
	connected    chan *client
	disconnected chan *client
	clients      map[*client]bool
}

func newLobby() lobby {
	return lobby{
		connected:    make(chan *client),
		disconnected: make(chan *client),
		clients:      make(map[*client]bool),
	}
}
func (l *lobby) start() {
	for {
		select {
		case client := <-l.connected:
			l.clients[client] = true
			log.Printf("A player has connected with ID %v", client.ID)
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{byte(client.ID)})
			broadcast(websocket.TextMessage, []byte("A new player connected"))
			if len(l.clients) == 2 {
				broadcast(websocket.TextMessage, []byte("Ready"))
				//start game
			}
		case client := <-l.disconnected:
			delete(l.clients, client)
			log.Printf("A player has disconnected with ID %v", client.ID)
			broadcast(websocket.TextMessage, []byte("A player has disconnected"))
		}
	}
}
