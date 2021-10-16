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

import (
	"fmt"
	"math"
	"math/rand"
)

type Ball struct {
	Pos      [2]float64
	Velocity [2]float64
	radius   float64
}

func NewBall() Ball {
	b := Ball{Pos: [2]float64{0, 0}}
	angle := rand.Float64()*90.0 - 45.0
	if rand.Float32() < 0.5 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	intialSpeed := rand.Float64()*25 + 75 //intial speed is between 75-100 units
	b.Velocity[0] = math.Cos(angle) * intialSpeed
	b.Velocity[1] = math.Sin(angle) * intialSpeed
	b.radius = 0.5 //should be the same here and C renderer to avoid wrong behavior
	return b
}

//game physics, called every frame
func (b *Ball) Update(dt float64, p []Player, resetFun func(i float64)) {
	//checks if any of the players scored
	if b.Pos[0] < -32 && b.Velocity[0] < 0 {
		p[1].Score++
		resetFun(-1)
		fmt.Printf("%v : %v \n", p[0].Score, p[1].Score)
		if p[1].Score > 9 {
			fmt.Println("Player Right Won !")
			isRunning = false
		}
	}
	if b.Pos[0] > 32 && b.Velocity[0] > 0 {
		p[0].Score++
		resetFun(1)
		fmt.Printf("%v : %v \n", p[0].Score, p[1].Score)
		if p[0].Score > 9 {
			fmt.Println("Player Left Won !")
			isRunning = false
		}
	}
	//calculates new Velocity
	b.Pos[0] += dt * b.Velocity[0]
	b.Pos[1] += dt * b.Velocity[1]
	b.resolveCollisions(p, dt)

}
func (b *Ball) resolveCollisions(p []Player, dt float64) {
	//checks floor and ceiling collisions
	if b.Pos[1]+b.radius > 18 && b.Velocity[1] > 0 {
		b.Velocity[1] *= -1
	}
	if b.Pos[1]-b.radius < -18 && b.Velocity[1] < 0 {
		b.Velocity[1] *= -1
	}
	//checks players collision and increases speed if bounced back
	b.resolvePlayer(p[0], dt)
	b.resolvePlayer(p[1], dt)
}
func (b *Ball) resolvePlayer(p Player, dt float64) {
	var ballSidePoint [2]float64
	if p.Pos[0] > 0 {
		ballSidePoint[0] = b.Pos[0] + b.radius
	} else {
		ballSidePoint[0] = b.Pos[0] - b.radius
	}
	ballSidePoint[1] = b.Pos[1]
	/*line A and B are determined by the Player's edge point and the Ball's
	Velocity, in order to determine Player reflecting area and
	prevent false scoring when Ball speed is too high*/
	var lineA, lineB line
	playerEdge := p.Pos
	if playerEdge[0] > 0 {
		playerEdge[0] -= p.width

	} else {
		playerEdge[0] += p.width
	}
	playerEdge[1] += p.length
	if b.Velocity[0] != 0 {
		lineA.a = b.Velocity[1] / b.Velocity[0]
		lineA.b = -(lineA.a)*playerEdge[0] + playerEdge[1]
	} else {
		lineA.a = 0
		lineA.b = playerEdge[1]
	}
	playerEdge[1] -= 2 * p.length
	if b.Velocity[0] != 0 {
		lineB.a = b.Velocity[1] / b.Velocity[0]
		lineB.b = -(lineB.a)*playerEdge[0] + playerEdge[1]
	} else {
		lineB.a = 0
		lineB.b = playerEdge[1]
	}
	ballBottom := [2]float64{b.Pos[0], b.Pos[1] - b.radius}
	ballTop := [2]float64{b.Pos[0], b.Pos[1] + b.radius}
	//checks if Ball is inside the Player reflecting area
	if ballBottom[1] <= lineA.a*ballBottom[0]+lineA.b &&
		ballTop[1] >= lineB.a*ballTop[0]+lineB.b {
		if playerEdge[0] > 0 && ballSidePoint[0] >= playerEdge[0] && b.Velocity[0] > 0 {
			b.Velocity[0] *= -1 * reflectionGain
			b.Velocity[1] *= reflectionGain
		}
		if playerEdge[0] < 0 && ballSidePoint[0] <= playerEdge[0] && b.Velocity[0] < 0 {
			b.Velocity[0] *= -1 * reflectionGain
			b.Velocity[1] *= reflectionGain
		}
	}
}
