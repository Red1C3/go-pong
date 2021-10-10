package game

// #cgo pkg-config: glfw3 glew cglm
// #cgo LDFLAGS:  -lm
// #include<Renderer.h>
import "C"
import (
	"log"
	"math"
	"math/rand"
	"net/url"
	"time"
)

type line struct {
	a, b float64
}

var (
	playerSpeed    = 30.0
	speedGain      = 1.01
	reflectionGain = 1.005
)
var (
	isOnline      bool
	isRunning     bool
	drawInfo      C.DrawInfo
	deltaTime     float64
	frameBeg      time.Time
	gameBall      ball
	players       [2]player
	savedVelocity [2]float64
	pauseTime     time.Time
)

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
	rand.Seed(time.Now().UnixNano())
	players[0] = newPlayer(-30, 6, 0.5)
	players[1] = newPlayer(30, 6, 0.5)
	gameBall = newBall()
	for isRunning {
		frameBeg = time.Now()
		if eventsHandler(drawInfo) == 1 {
			return
		}
	}
}
func updateDrawInfo() {
	if gameBall.pos[0] > players[0].pos[0] && gameBall.pos[0] < players[1].pos[0] {
		drawInfo.ball[0] = C.float(gameBall.pos[0])
		drawInfo.ball[1] = C.float(gameBall.pos[1])
	} else {
		drawInfo.ball[0] = -100
		drawInfo.ball[1] = -100
	}
	drawInfo.p1 = C.float(players[0].pos[1])
	drawInfo.p2 = C.float(players[1].pos[1])
}
func eventsHandler(dI C.DrawInfo) int {
	event := C.loop(dI)
	if gameBall.velocity[0] == 0 && gameBall.velocity[1] == 0 &&
		time.Since(pauseTime).Seconds() > 0.5 {
		gameBall.velocity = savedVelocity
	}
	switch event.code {
	case 0:
		gameBall.update(deltaTime, players[:])
		updateDrawInfo()
		deltaTime = float64(time.Since(frameBeg).Seconds())
		return 0
	case 1:
		return 1
	case 2:

		switch event.key {
		case 'w':
			players[0].move(playerSpeed, deltaTime)
		case 's':
			players[0].move(-playerSpeed, deltaTime)
		case 'u':
			players[1].move(playerSpeed, deltaTime)
		case 'd':
			players[1].move(-playerSpeed, deltaTime)
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.code)
		return -1
	}
}
func reset(i float64) {
	gameBall.pos = [2]float64{i * 25, 0}
	velocityLength := math.Sqrt(math.Pow(gameBall.velocity[0], 2) + math.Pow(gameBall.velocity[1], 2))
	angle := rand.Float64()*120 - 60
	if gameBall.velocity[0] > 0 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	savedVelocity[0] = math.Cos(angle) * velocityLength * speedGain
	savedVelocity[1] = math.Sin(angle) * velocityLength * speedGain
	gameBall.velocity = [2]float64{0, 0}
	players[0].pos[1] = 0
	players[1].pos[1] = 0
	pauseTime = time.Now()
}
