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
package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-pong/game"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

// structure used for server recieving
type exchangeData struct {
	mutex                sync.RWMutex
	P1, P2, BallX, BallY float64
}

const (
	CLOSE_MSG            = byte(0x0F)
	READY_MSG            = byte(0x1F)
	CONNECT_MSG          = byte(0x2F)
	OTHER_DISCONNECT_MSG = byte(0x3F)
	ID_MSG               = byte(0x4F)
	OTHER_CONNECT_MSG    = byte(0x5F)
	DATA_MSG             = byte(0x6F)
	SCORE_MSG            = byte(0x7F)
)

var (
	data   exchangeData
	client struct {
		ID         int
		connection *net.UDPConn
	}
	scores [2]int
	buffer bytes.Buffer

	//decodes received data
	decoder      *gob.Decoder
	drawInfo     game.CDrawInfo
	closeChannel = make(chan bool)
	readyChannel = make(chan bool)
)

func Start() {
	decoder = gob.NewDecoder(&buffer)
	drawInfo = game.CDrawInfo{
		Ball: [2]float64{0, 0},
		P1:   0,
		P2:   0,
	}
	udpServer, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		log.Fatal("Failed to resolve address ", os.Args[1], " error: ", err.Error())
	}
	client.connection, err = net.DialUDP("udp", nil, udpServer)
	if err != nil {
		log.Fatal("Failed to dial server, error:", err.Error())
	}
	log.Println("Connected to server")
	_, err = client.connection.Write([]byte{CONNECT_MSG}) //Connecting message
	if err != nil {
		log.Fatal("Failed to send connecting message to server")
	}

	go listenToServer()
	<-readyChannel
	startGame()
	game.Terminate()
}
func closeConnection(sendCloseMsg bool) {
	fmt.Println("Closing connection...")
    if sendCloseMsg{
        _, err := client.connection.Write([]byte(CLOSE_MSG))
        if err != nil {
            log.Print("Failed to send closing message to server, error:", err.Error())
        }
    }
	err := client.connection.Close()
	if err != nil {
		log.Fatal("Failed to close connection to server")
	}
}
func startGame() {
	err := game.CInitRenderer()
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-closeChannel:
            closeConnection(false)
			return
		default:
		}
		updateDrawInfo()
		if eventsHandler(drawInfo) == 1 {
            closeChannel<-true
            closeConnection(true)
			return
		}
	}
}
func eventsHandler(dI game.CDrawInfo) int {
	event := game.CLoop(dI)
	switch event.Code {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		var err error
		switch event.Key {
		case 'u':
			_, err = client.connection.Write([]byte{1})
		case 'd':
			_, err = client.connection.Write([]byte{0})
		}
		if err != nil {
			log.Fatal("Failed to send direction byte to server, error:", err.Error())
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.Code)
		return -1
	}
}
func listenToServer() {
	var structure struct {
		P1, P2, BallX, BallY float64
	}

	var b [64]byte

	for {
        select{
        case <-closeChannel:
            return
        default:
        }
		n, err := client.connection.Read(b[:])
		if err != nil {
			log.Print("Failed to read from server, error:", err.Error())
		}
		switch b[0] {
		case CLOSE_MSG:
            closeChannel<-true
            return
		case DATA_MSG:
			buffer.Reset()
			_, err = buffer.Write(b[1:])
			if err != nil {
				log.Fatal("Failed to write message content to buffer, error:", err.Error())
			}
			err = decoder.Decode(&structure)
			if err != nil {
				log.Print("Failed to decode data msg, error: ", err.Error())
			}
			data.mutex.Lock()
			data.P2 = structure.P2
			data.P1 = structure.P1
			data.BallX = structure.BallX
			data.BallY = structure.BallY
			data.mutex.Unlock()
		case ID_MSG:
			client.ID = int(b[1])
		case OTHER_CONNECT_MSG:
			log.Print("A player has joined with ID ", b[1])
		case OTHER_DISCONNECT_MSG:
			log.Print("A player has disconnected with ID ", b[1])
		case READY_MSG:
			readyChannel <- true
		case SCORE_MSG:
			log.Print(string(b[1:n]))
			scores[0], err = strconv.Atoi(string(b[1]))
			if err != nil {
				log.Fatal("Invalid score")
			}
			scores[1], err = strconv.Atoi(string(b[5]))
			if err != nil {
				log.Fatal("Invalid score")
			}
		default:
			log.Print("Server sent unknown message type:", b[0])
		}
	}
}
func updateDrawInfo() {
	data.mutex.RLock()
	defer data.mutex.RUnlock()
	drawInfo.P1 = data.P1
	drawInfo.P2 = data.P2
	drawInfo.Ball[0] = data.BallX
	drawInfo.Ball[1] = data.BallY
	drawInfo.Scores[0] = scores[0]
	drawInfo.Scores[1] = scores[1]
}
