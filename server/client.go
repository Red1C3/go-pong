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
	"go-pong/client"
	"log"
	"net"
)

type clientStr struct {
	ID      byte
	address net.Addr
}

func listenToClient() {
	var buffer [64]byte
	for {
		_, addr, err := udpServerHandle.ReadFrom(buffer[:])
		if addr == nil {
			continue
		}
		if err != nil {
			log.Print("Failed to read input from ", addr.String(), "error:", err.Error())
			continue
		}
		switch buffer[0] {
		case client.CLOSE_MSG:
			closeChannel <- true
		case client.DATA_MSG:
			var dir int
			if buffer[1] == 0 {
				dir = -1
			} else {
				dir = 1
			}
			playersMutex[gameLobby.clients[addr.String()].ID].Lock()
			players[gameLobby.clients[addr.String()].ID].Move(float64(dir)*30, deltaTime)
			playersMutex[gameLobby.clients[addr.String()].ID].Unlock()
		default:
			log.Print("Unknown message type recieved from ", addr.String(), " message type:", buffer[0])
		}
	}
}