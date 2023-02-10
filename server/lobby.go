/*
MIT License

# Copyright (c) 2021 Mohammad Issawi

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
	"fmt"
    "go-pong/client"
)


type lobby struct {
	connected    chan *clientStr
	disconnected chan *clientStr
	clients      map[string]*clientStr
}

func newLobby() lobby {
	return lobby{
		connected:    make(chan *clientStr),
		disconnected: make(chan *clientStr),
		clients:      make(map[string]*clientStr),
	}
}
func (l *lobby) start() {
	for {
		select {
		case c := <-l.connected:
			l.clients[c.address.String()] = c
			fmt.Printf("A Player has connected with ID %v \n", c.ID)
			sendToAddress(c.address, []byte{c.ID})
			broadcast([]byte("A new Player connected"))
			if len(l.clients) == 2 {
				broadcast([]byte(client.READY_MSG))
				startGame()
			}
		case c := <-l.disconnected:
			delete(l.clients, c.address.String())
			fmt.Printf("A Player has disconnected with ID %v \n", c.ID)
			broadcast([]byte("A Player has disconnected"))
		}
	}
}
