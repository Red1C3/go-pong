package game

// #cgo pkg-config: glfw3 glew cglm
// #include<Renderer.h>
import "C"
import (
	"log"
	"net/url"
	"time"
)

var isOnline bool
var isRunning bool
var drawInfo C.DrawInfo
var deltaTime float64
var frameBeg time.Time
var speed = float64(10.0)

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
		frameBeg = time.Now()
		if eventsHandler(drawInfo) == 1 {
			return
		}
		deltaTime = float64(time.Since(frameBeg).Seconds())
	}
}
func eventsHandler(dI C.DrawInfo) int {
	event := C.loop(dI)
	switch event.code {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		switch event.key {
		case 'w':
			drawInfo.p1 += C.float(speed * deltaTime)
		case 's':
			drawInfo.p1 -= C.float(speed * deltaTime)
		case 'u':
			drawInfo.p2 += C.float(speed * deltaTime)
		case 'd':
			drawInfo.p2 -= C.float(speed * deltaTime)
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.code)
		return -1
	}
}
