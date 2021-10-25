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
// #define STB_IMAGE_IMPLEMENTATION
// #include<stb/stb_image.h>
// #include<Renderer.h>
import "C"
import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

//geomatric line, formatted as y = a * x + b
type line struct {
	a, b float64
}

//game constants, change for different difficulty
const (
	playerSpeed    = 30.0
	ScoreGain      = 1.01
	reflectionGain = 1.005
	ResetTime      = 0.8
)

var (
	isRunning     bool
	drawInfo      C.DrawInfo
	deltaTime     float64
	frameBeg      time.Time
	gameBall      Ball
	players       [2]Player
	savedVelocity [2]float64
	pauseTime     time.Time
)

func Run() {
	if err := C.initRenderer(false); err != 0 {
		log.Fatalf("Failed to init renderer, error code: %v", err)
	}
	isRunning = true
	gameLogic()
	Terminate()
}
func Terminate() {
	C.terminateRenderer()
}
func gameLogic() {
	//feeds random
	rand.Seed(time.Now().UnixNano())
	players[0] = NewPlayer(-30, 6, 0.5)
	players[1] = NewPlayer(30, 6, 0.5)
	gameBall = NewBall()
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
	if gameBall.Pos[0] > players[0].Pos[0] && gameBall.Pos[0] < players[1].Pos[0] {
		drawInfo.ball[0] = C.float(gameBall.Pos[0])
		drawInfo.ball[1] = C.float(gameBall.Pos[1])
	} else {
		drawInfo.ball[0] = -100
		drawInfo.ball[1] = -100
	}
	drawInfo.p1 = C.float(players[0].Pos[1])
	drawInfo.p2 = C.float(players[1].Pos[1])
	drawInfo.scores[0] = C.int(players[0].Score)
	drawInfo.scores[1] = C.int(players[1].Score)
}

//Handles C events
func eventsHandler(dI C.DrawInfo) int {
	event := C.loop(dI)
	//if resat, wait for (ResetTime) seconds before starting...
	if gameBall.Velocity[0] == 0 && gameBall.Velocity[1] == 0 &&
		time.Since(pauseTime).Seconds() > ResetTime {
		gameBall.Velocity = savedVelocity
	}
	switch event.code {
	case 0:
		gameBall.Update(deltaTime, players[:], reset)
		updateDrawInfo()
		deltaTime = time.Since(frameBeg).Seconds()
		return 0
	case 1:
		return 1
	case 2:
		//Handle input
		switch event.key {
		case 'w':
			players[0].Move(playerSpeed, deltaTime)
		case 's':
			players[0].Move(-playerSpeed, deltaTime)
		case 'u':
			players[1].Move(playerSpeed, deltaTime)
		case 'd':
			players[1].Move(-playerSpeed, deltaTime)
		}
		return 2
	default:
		log.Fatalf("Unknown event code: %v", event.code)
		return -1
	}
}

//called when a Player scores
func reset(i float64) {
	gameBall.Pos = [2]float64{i * 25, 0}
	//create a new Velocity vector with the same speed of the current one
	//but with a different angle
	velocityLength := math.Sqrt(math.Pow(gameBall.Velocity[0], 2) + math.Pow(gameBall.Velocity[1], 2))
	angle := rand.Float64()*120 - 60
	if gameBall.Velocity[0] > 0 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	savedVelocity[0] = math.Cos(angle) * velocityLength * ScoreGain
	savedVelocity[1] = math.Sin(angle) * velocityLength * ScoreGain
	//pause Ball until reset time is passed
	gameBall.Velocity = [2]float64{0, 0}
	players[0].Pos[1] = 0
	players[1].Pos[1] = 0
	pauseTime = time.Now()
}

type CDrawInfo struct {
	P1, P2 float64
	Ball   [2]float64
	Scores [2]int
}
type CEvent struct {
	Code int
	Key  rune
}

func CInitRenderer() error {
	if err := C.initRenderer(true); err != 0 {
		return fmt.Errorf("failed to init renderer, error code: %v", err)
	}
	return nil
}
func CLoop(dI CDrawInfo) CEvent {
	drawInfo.p1 = C.float(dI.P1)
	drawInfo.p2 = C.float(dI.P2)
	drawInfo.ball[0] = C.float(dI.Ball[0])
	drawInfo.ball[1] = C.float(dI.Ball[1])
	drawInfo.scores[0] = C.int(dI.Scores[0])
	drawInfo.scores[1] = C.int(dI.Scores[1])
	event := C.loop(drawInfo)
	return CEvent{
		Code: int(event.code),
		Key:  rune(event.key),
	}
}
