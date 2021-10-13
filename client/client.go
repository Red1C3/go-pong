package client

import (
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
	//start game
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
