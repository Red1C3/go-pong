/*MIT License

Copyright (c) 2021 Mohammad Issawi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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
			log.Printf("A Player has connected with ID %v", client.ID)
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{byte(client.ID)})
			broadcast(websocket.TextMessage, []byte("A new Player connected"))
			if len(l.clients) == 2 {
				broadcast(websocket.TextMessage, []byte("Ready"))
				startGame()
			}
		case client := <-l.disconnected:
			delete(l.clients, client)
			log.Printf("A Player has disconnected with ID %v", client.ID)
			broadcast(websocket.TextMessage, []byte("A Player has disconnected"))
		}
	}
}
