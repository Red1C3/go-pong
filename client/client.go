package client

import (
	"bytes"
	"encoding/gob"
	"go-pong/game"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type exchangeData struct {
	mutex                sync.RWMutex
	P1, P2, BallX, BallY float64
}

var data exchangeData
var client struct {
	ID         int
	connection *websocket.Conn
}
var buffer bytes.Buffer
var decoder *gob.Decoder
var drawInfo game.CDrawInfo

func Start() {
	var err error
	decoder = gob.NewDecoder(&buffer)
	drawInfo = game.CDrawInfo{
		Ball: [2]float64{0, 0},
		P1:   0,
		P2:   0,
	}
	serverURL := url.URL{Scheme: "ws", Host: os.Args[1], Path: "/"}
	client.connection, _, err = websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	_, p, err := client.connection.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	client.ID = int(p[0])
	log.Printf("Connected as Player %v", client.ID)
	log.Print("Waiting for other players to join")
	for {
		msgType, p, err := client.connection.ReadMessage()
		if err != nil {
			log.Printf("Failed to read msg from server %v", err)
			closeConnection()
		}
		if msgType == websocket.TextMessage {
			if string(p) == "Ready" {
				break
			}
			log.Print(string(p))
		}
	}
	log.Print("Players connected, starting game...")
	go msgsHandler()
	startGame()
	closeConnection()
}
func closeConnection() {
	log.Print("Closing connection...")
	err := client.connection.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second))
	if err != nil {
		log.Printf("Failed to close connection %v", err)
	}
	client.connection.Close()
}
func startGame() {
	err := game.CInitRenderer()
	if err != nil {
		log.Fatal(err)
	}
	for {
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
		switch event.Key {
		case 'u':
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{1})
		case 'd':
			client.connection.WriteMessage(websocket.BinaryMessage, []byte{0})
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.Code)
		return -1
	}
}
func msgsHandler() {
	for {
		msgType, p, err := client.connection.ReadMessage()
		if err != nil {
			log.Fatalf("Error occured while reciving a msg %v", err)
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
			log.Print(string(p))
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
}
