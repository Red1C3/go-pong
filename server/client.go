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
			delete(gameLobby.clients, addr.String())
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