package game

// #cgo pkg-config: glfw3 glew cglm
// #include<Renderer.h>
import "C"
import (
	"log"
	"net/url"
)

var isOnline bool
var isRunning bool

func Run(u *url.URL) {
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
	if !isOnline {
		gameLogic()
	}
	terminate()
}
func terminate() {
	C.terminateRenderer()
	isRunning = false
}
func gameLogic() {
	for {
		var dummy C.DrawInfo
		for loop := true; loop; {
			switch eventsHandler(dummy) {
			case 0:
				loop = false
			case 1:
				return
			default:
				continue
			}
		}
	}
}
func eventsHandler(dI C.DrawInfo) int {
	event := C.loop(dI)
	switch event.code {
	case 0:
		return 0
	case 1:
		return 1
	default:
		log.Fatalf("Unknown event code: %v", event.code)
		return -1
	}
}
