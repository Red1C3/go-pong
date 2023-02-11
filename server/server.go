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
	"bytes"
	"encoding/gob"
	"fmt"
	"go-pong/client"
	"go-pong/game"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

var (
	port            string
	udpServerHandle net.PacketConn
	gameLobby       = newLobby()
	players         [2]game.Player
	playersMutex    [2]sync.RWMutex
	gameBall        game.Ball
	pauseTime       time.Time
	savedVelocity   [2]float64
	deltaTime       float64
	encoder         *gob.Encoder
	buffer          bytes.Buffer
	closeChannel    = make(chan bool, 1)
)

func Start() {
	var err error
	port = os.Args[2]
	log.Print("Starting server...")
	udpServerHandle, err = net.ListenPacket("udp", ":"+port)
	if err != nil {
		log.Fatal("Failed to start listening error: ", err.Error())
	}
	waitForPlayers()
	gameLobby.start()
	err = udpServerHandle.Close()
	if err != nil {
		log.Print("Failed to close UDP server, error: ", err.Error())
	}
}

func waitForPlayers() {
	var buffer [64]byte
	log.Print("Waiting for players...")
	for{
        if len(gameLobby.clients)==2{
            return
        }
        n,addr,err:=udpServerHandle.ReadFrom(buffer[:])
        if err!=nil{
            if addr!=nil{
                log.Print("Failed to read from address:",addr.String(),", error:",err.Error())
            }else{
                log.Print("Failed to read, error:",err.Error())
            }
        }
        if c,ok:=gameLobby.clients[addr.String()];ok{
            if n==1 && buffer[0]==client.CLOSE_MSG{
                delete(gameLobby.clients,addr.String())
                fmt.Printf("A Player has disconnected with ID %v \n", c.ID)
                broadcast([]byte{client.OTHER_DISCONNECT_MSG,c.ID})
            }else{
                log.Print("Client sent an unexpected message:",string(buffer[:n]))
            }
        }else if n==1 && buffer[0]==client.CONNECT_MSG{
            newClient:=&clientStr{
                ID:byte(len(gameLobby.clients)),
                address: addr,
            }
            gameLobby.clients[addr.String()] = newClient
            fmt.Printf("A Player has connected with ID %v \n", newClient.ID)
            sendToAddress(c.address, []byte{client.ID_MSG,newClient.ID})
            broadcast([]byte{client.OTHER_CONNECT_MSG,newClient.ID})
        }else{
            log.Print("Non-client sent an unexpected message, address:",addr.String()," messages:",string(buffer[:n]))
        }
    }
}
func broadcast(content []byte) {
	for _, c := range gameLobby.clients {
		_, err := udpServerHandle.WriteTo(content, c.address)
		if err != nil {
			log.Print("Failed to send message to address ", c.address.String(), " error:", err.Error())
		}
	}
}

func sendToAddress(addr net.Addr, msg []byte) {
	_, err := udpServerHandle.WriteTo(msg, addr)
	if err != nil {
		log.Print("Failed to send message to address ", addr.String(), " error:", err.Error())
	}
}
func startGame() {
	time.Sleep(time.Second * 3)
	encoder = gob.NewEncoder(&buffer)
	rand.Seed(time.Now().UnixNano())
	players[0] = game.NewPlayer(-30, 6, 0.5)
	players[1] = game.NewPlayer(30, 6, 0.5)
	gameBall = game.NewBall()
	netTicker := time.NewTicker(time.Second / 120)
	gameTicker := time.NewTicker(time.Second / 80)
	deltaTime = ((time.Second) / 80).Seconds()
	for close := false; !close; {
		select {
		case <-closeChannel:
			close = true
		case <-netTicker.C:
			playersMutex[0].RLock()
			playersMutex[1].RLock()
			broadcastData()
			playersMutex[0].RUnlock()
			playersMutex[1].RUnlock()
		case <-gameTicker.C:
			if gameBall.Velocity[0] == 0 && gameBall.Velocity[1] == 0 &&
				time.Since(pauseTime).Seconds() > game.ResetTime {
				gameBall.Velocity = savedVelocity
			}
			playersMutex[0].Lock()
			playersMutex[1].Lock()
			gameBall.Update(deltaTime, players[:], reset)
			playersMutex[0].Unlock()
			playersMutex[1].Unlock()
		}
	}
	broadcast([]byte(client.CLOSE_MSG))
}
func broadcastData() {
	var structure struct {
		P1, P2, BallX, BallY float64
	}
	structure.P1 = players[0].Pos[1]
	structure.P2 = players[1].Pos[1]
	structure.BallX = gameBall.Pos[0]
	structure.BallY = gameBall.Pos[1]
	buffer.Reset()
	err := encoder.Encode(structure)
	if err != nil {
		log.Fatal(err)
	}
	broadcast(buffer.Bytes())
}
func reset(i float64) {
	gameBall.Pos = [2]float64{i * 25, 0}
	//create a new Velocity vector with the same speed of the current one
	//but with a different angle
	velocityLength := math.Sqrt(math.Pow(gameBall.Velocity[0], 2) + math.Pow(gameBall.Velocity[1], 2))
	angle := rand.Float64()*120 - 60
	if gameBall.Velocity[0] > 0 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	savedVelocity[0] = math.Cos(angle) * velocityLength * game.ScoreGain
	savedVelocity[1] = math.Sin(angle) * velocityLength * game.ScoreGain
	//pause Ball until reset time is passed
	gameBall.Velocity = [2]float64{0, 0}
	players[0].Pos[1] = 0
	players[1].Pos[1] = 0
	pauseTime = time.Now()
	broadcast([]byte(fmt.Sprintf("%v : %v", players[0].Score, players[1].Score)))
	if players[0].Score > 9 {
		broadcast([]byte("Player left won !"))
		closeChannel <- true
	}
	if players[1].Score > 9 {
		broadcast([]byte("Player right won !"))
		closeChannel <- true
	}
}
