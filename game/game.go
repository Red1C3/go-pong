package game

// #cgo pkg-config: glfw3 glew
// #include<Renderer.h>
import "C"
import (
	"log"
	"net/url"
)

var isOnline bool
var isRunning bool
var closingChannel chan (bool)

func Run(u *url.URL) {
	closingChannel = make(chan bool, 1)
	if u != nil {
		isOnline = true
	} else {
		isOnline = false
	}
	if isOnline {
		log.Fatal("Not implemented yet")
	}
	if err := C.initRenderer(); err != 0 {
		log.Fatalf("Failed to init renderer, error code: %v", err)
	}
	isRunning = true
	go eventsHandler()
	if err := <-closingChannel; !err {
		log.Fatal("An error occured")
	}
	terminate()
}
func terminate() {
	C.terminateRenderer()
	isRunning = false
}
func eventsHandler() {
	for {
		event := C.loop()
		switch event.code {
		case 0:
		case 1:
			closingChannel <- true
			return
		default:
			log.Fatalf("Unknown event code: %v", event.code)
		}
	}
}
