package game

// #cgo pkg-config: glfw3 glew
// #include<Renderer.h>
import "C"
import (
	"log"
	"net/url"
	"time"
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
	time.Sleep(time.Second * 4)
	C.terminateRenderer()
}
