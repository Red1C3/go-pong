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
	"bytes"
	"encoding/gob"
	"fmt"
	"go-pong/game"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	port     string
	upgrader = websocket.Upgrader{
		WriteBufferSize: 128,
		ReadBufferSize:  128,
	}
	gameLobby     = newLobby()
	players       [2]game.Player
	playersMutex  [2]sync.RWMutex
	gameBall      game.Ball
	pauseTime     time.Time
	savedVelocity [2]float64
	deltaTime     float64
	encoder       *gob.Encoder
	buffer        bytes.Buffer
	closeChannel  = make(chan bool, 1)
)

func Start() {
	port = os.Args[2]
	log.Print("Starting server...")
	http.HandleFunc("/", requestsHandler)
	go http.ListenAndServe(":"+port, nil)
	log.Print("Waiting for players...")
	gameLobby.start()
}
func requestsHandler(w http.ResponseWriter, r *http.Request) {
	if len(gameLobby.clients) == 2 {
		return
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to connect to a client %v", err)
	}
	client := &client{
		ID:         len(gameLobby.clients) + 1,
		connection: connection,
	}
	gameLobby.connected <- client
	client.start()
}
func broadcast(msgType int, content []byte) {
	for client := range gameLobby.clients {
		client.connection.WriteMessage(msgType, content)
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
	for client := range gameLobby.clients {
		log.Print("Closing connection...")
		err := client.connection.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(time.Second))
		if err != nil {
			if err != websocket.ErrCloseSent {
				log.Printf("Failed to close connection %v", err)
			}
		}
		client.connection.Close()
	}
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
	broadcast(websocket.BinaryMessage, buffer.Bytes())
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
	broadcast(websocket.TextMessage, []byte(fmt.Sprintf("%v : %v", players[0].Score, players[1].Score)))
	if players[0].Score > 9 {
		broadcast(websocket.TextMessage, []byte("Player left won !"))
		closeChannel <- true
	}
	if players[1].Score > 9 {
		broadcast(websocket.TextMessage, []byte("Player right won !"))
		closeChannel <- true
	}
}
