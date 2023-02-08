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
	"sync"
)

// structure used for server recieving
type exchangeData struct {
	mutex                sync.RWMutex
	P1, P2, BallX, BallY float64
}

const (
	CLOSE_MSG = "close"
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
	_, err = client.connection.Write([]byte("")) //Connecting message
	if err != nil {
		log.Fatal("Failed to send connecting message to server")
	}

	var id [1]byte
	_, err = client.connection.Read(id[:])
	if err != nil {
		log.Fatal("Failed to recieve id from server, error:", err.Error())
	}
	client.ID = int(id[0])
	fmt.Printf("Connected as Player %v \n", client.ID)
	fmt.Println("Waiting for other players to join")
	closeConnection()
	return
	/*for {
		msgType, p, err := client.connection.ReadMessage()
		if err != nil {
			fmt.Printf("Failed to read msg from server %v", err)
			closeConnection()
			return
		}
		if msgType == websocket.TextMessage {
			if string(p) == "Ready" {
				break
			}
		}
	}*/
	fmt.Println("Players connected, starting game...")
	go msgsHandler()
	startGame()
	closeConnection()
	game.Terminate()
}
func closeConnection() {
	fmt.Println("Closing connection...")
	_, err := client.connection.Write([]byte(CLOSE_MSG))
	if err != nil {
		log.Print("Failed to send closing message to server, error:", err.Error())
	}
	err = client.connection.Close()
	if err != nil {
		log.Fatal("Failed to close connection to server")
	}
}
func startGame() {
	err := game.CInitRenderer()
	if err != nil {
		log.Fatal(err)
	}
	for close := false; !close; {
		select {
		case <-closeChannel:
			close = true
		default:
		}
		updateDrawInfo()
		if eventsHandler(drawInfo) == 1 {
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
		/*switch event.Key {
		case 'u':
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{1})
		case 'd':
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{0})
		}*/

		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.Code)
		return -1
	}
}
func msgsHandler() {
	/*for {
		msgType, p, err := client.connection.ReadMessage()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				if ce.Code == websocket.CloseNormalClosure {
					fmt.Println("Connection closed from server")
					closeChannel <- true
					return
				}
			}
			log.Fatalf("error while reading msg from client %v : %v", client.ID, err)
		}
		if msgType == websocket.BinaryMessage {
			buffer.Reset()
			_, err := buffer.Write(p)
			if err != nil {
				log.Fatalf("Error while writing to buffer %v", err)
			}
			var structure struct {
				P1, P2, BallX, BallY float64
			}
			err = decoder.Decode(&structure)
			if err != nil {
				log.Fatalf("Decoder err %v", err)
			}
			data.mutex.Lock()
			data.P1 = structure.P1
			data.P2 = structure.P2
			data.BallX = structure.BallX
			data.BallY = structure.BallY
			data.mutex.Unlock()
		}
		if msgType == websocket.TextMessage {
			fmt.Println(string(p))
			if len(string(p)) == 5 {
				scores[0], err = strconv.Atoi(string(p[0]))
				if err != nil {
					log.Fatal("Invalid score")
				}
				scores[1], err = strconv.Atoi(string(p[4]))
				if err != nil {
					log.Fatal("Invalid score")
				}
			}
		}
	}*/
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
