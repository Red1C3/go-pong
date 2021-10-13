package client

import (
	"go-pong/game"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var client struct {
	ID         int
	connection *websocket.Conn
}

func Start() {
	var err error
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
	var drawInfo game.CDrawInfo
	for {
		if eventsHandler(drawInfo) == 1 {
			return
		}
	}
}
func eventsHandler(dI game.CDrawInfo) int {
	event := game.CLoop(dI)

	switch event.Code {
	case 0:
		updateDrawInfo()
		return 0
	case 1:
		return 1
	case 2:
		switch event.Key {
		case 'u':
			//players[1].move(playerSpeed, deltaTime)
		case 'd':
			//players[1].move(-playerSpeed, deltaTime)
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.Code)
		return -1
	}
}
func updateDrawInfo() {

}
