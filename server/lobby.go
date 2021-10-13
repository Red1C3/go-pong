package server

import (
	"log"
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
			log.Printf("new player connected with ID %v", client.ID)
		case client := <-l.disconnected:
			log.Printf("new player disconnected with ID %v", client.ID)
		}
	}
}
