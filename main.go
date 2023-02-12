package main

import (
	"fmt"
	"go-pong/client"
	"go-pong/game"
	"go-pong/server"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		game.Run() //starts the game offline
	} else if os.Args[1] == "-h" {
		server.Start() //starts host
	} else if len(os.Args) == 2 {
		client.Start() //starts client
	}
	fmt.Println("Created with fuzzy kittens, with the help of Red1C3")
}
