package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var port = os.Args[2]
var upgrader = websocket.Upgrader{
	WriteBufferSize: 128,
	ReadBufferSize:  128,
}
var gameLobby = newLobby()

func Start() {
	log.Print("Starting server...")
	http.HandleFunc("/", requestsHandler)
	go http.ListenAndServe(":"+port, nil)
	log.Print("Waiting for players...")
	gameLobby.start()
}
func requestsHandler(w http.ResponseWriter, r *http.Request) {
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
