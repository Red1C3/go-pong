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

var port string
var upgrader = websocket.Upgrader{
	WriteBufferSize: 128,
	ReadBufferSize:  128,
}
var gameLobby = newLobby()
var players [2]game.Player
var playersMutex [2]sync.RWMutex
var gameBall game.Ball
var pauseTime time.Time
var savedVelocity [2]float64
var deltaTime float64
var encoder *gob.Encoder
var buffer bytes.Buffer
var closeChannel = make(chan bool, 1)

const resetTime = 0.8
const scoreGain = 1.01

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
				time.Since(pauseTime).Seconds() > resetTime {
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
	savedVelocity[0] = math.Cos(angle) * velocityLength * scoreGain
	savedVelocity[1] = math.Sin(angle) * velocityLength * scoreGain
	//pause Ball until reset time is passed
	gameBall.Velocity = [2]float64{0, 0}
	players[0].Pos[1] = 0
	players[1].Pos[1] = 0
	pauseTime = time.Now()
	broadcast(websocket.TextMessage, []byte(fmt.Sprintf("%v : %v", players[0].Score, players[1].Score)))
	if players[0].Score > 9 {
		broadcast(websocket.TextMessage, []byte("Player right won !"))
		closeChannel <- true
	}
	if players[1].Score > 9 {
		broadcast(websocket.TextMessage, []byte("Player left won !"))
		closeChannel <- true
	}
}
