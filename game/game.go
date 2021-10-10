/*MIT License

Copyright (c) 2021 Mohammad Issawi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package game

// #cgo pkg-config: glfw3 glew cglm
// #cgo LDFLAGS:  -lm
// #include<Renderer.h>
import "C"
import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"time"
)

//geomatric line, formatted as y = a * x + b
type line struct {
	a, b float64
}

//game constants, change for different difficulty
const (
	playerSpeed    = 30.0
	speedGain      = 1.01
	reflectionGain = 1.005
	resetTime      = 0.8
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
	fmt.Println("Created with fuzzy kittens, with the help of RedDeadAlice")
}
func terminate() {
	C.terminateRenderer()
}
func gameLogic() {
	//feeds random
	rand.Seed(time.Now().UnixNano())
	players[0] = newPlayer(-30, 6, 0.5)
	players[1] = newPlayer(30, 6, 0.5)
	gameBall = newBall()
	for isRunning {
		frameBeg = time.Now()
		if eventsHandler(drawInfo) == 1 {
			isRunning = false
			return
		}
	}
}

//Updates the structure sent to C code to draw properlys
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

//Handles C events
func eventsHandler(dI C.DrawInfo) int {
	event := C.loop(dI)
	//if resat, wait for (resetTime) seconds before starting...
	if gameBall.velocity[0] == 0 && gameBall.velocity[1] == 0 &&
		time.Since(pauseTime).Seconds() > resetTime {
		gameBall.velocity = savedVelocity
	}
	switch event.code {
	case 0:
		gameBall.update(deltaTime, players[:])
		updateDrawInfo()
		deltaTime = time.Since(frameBeg).Seconds()
		return 0
	case 1:
		return 1
	case 2:
		//Handle input
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

//called when a player scores
func reset(i float64) {
	gameBall.pos = [2]float64{i * 25, 0}
	//create a new velocity vector with the same speed of the current one
	//but with a different angle
	velocityLength := math.Sqrt(math.Pow(gameBall.velocity[0], 2) + math.Pow(gameBall.velocity[1], 2))
	angle := rand.Float64()*120 - 60
	if gameBall.velocity[0] > 0 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	savedVelocity[0] = math.Cos(angle) * velocityLength * speedGain
	savedVelocity[1] = math.Sin(angle) * velocityLength * speedGain
	//pause ball until reset time is passed
	gameBall.velocity = [2]float64{0, 0}
	players[0].pos[1] = 0
	players[1].pos[1] = 0
	pauseTime = time.Now()
}
