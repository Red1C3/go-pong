package game

import (
	"log"
	"math"
	"math/rand"
)

type ball struct {
	pos      [2]float64
	velocity [2]float64
	radius   float64
}

func newBall() ball {
	b := ball{pos: [2]float64{0, 0}}
	angle := rand.Float64()*90.0 - 45.0
	if rand.Float32() < 0.5 {
		angle += 180
	}
	angle = angle * math.Pi / 180
	intialSpeed := rand.Float64()*25 + 75 //intial speed is between 75-100 units
	b.velocity[0] = math.Cos(angle) * intialSpeed
	b.velocity[1] = math.Sin(angle) * intialSpeed
	b.radius = 0.5 //should be the same here and C renderer to avoid wrong behavior
	return b
}

//game physics, called every frame
func (b *ball) update(dt float64, p []player) {
	//checks if any of the players scored
	if b.pos[0] < -32 && b.velocity[0] < 0 {
		reset(-1)
		p[1].score++
		log.Printf("%v : %v", p[0].score, p[1].score)
		if p[1].score > 9 {
			log.Print("Player Right Won !")
			isRunning = false
		}
	}
	if b.pos[0] > 32 && b.velocity[0] > 0 {
		reset(1)
		p[0].score++
		log.Printf("%v : %v", p[0].score, p[1].score)
		if p[0].score > 9 {
			log.Print("Player Left Won !")
			isRunning = false
		}
	}
	//calculates new velocity
	b.pos[0] += dt * b.velocity[0]
	b.pos[1] += dt * b.velocity[1]
	b.resolveCollisions(p, dt)

}
func (b *ball) resolveCollisions(p []player, dt float64) {
	//checks floor and ceiling collisions
	if b.pos[1]+b.radius > 18 && b.velocity[1] > 0 {
		b.velocity[1] *= -1
	}
	if b.pos[1]-b.radius < -18 && b.velocity[1] < 0 {
		b.velocity[1] *= -1
	}
	//checks players collision and increases speed if bounced back
	b.resolvePlayer(p[0], dt)
	b.resolvePlayer(p[1], dt)
}
func (b *ball) resolvePlayer(p player, dt float64) {
	var ballSidePoint [2]float64
	if p.pos[0] > 0 {
		ballSidePoint[0] = b.pos[0] + b.radius
	} else {
		ballSidePoint[0] = b.pos[0] - b.radius
	}
	ballSidePoint[1] = b.pos[1]
	/*line A and B are determined by the player's edge point and the ball's
	velocity, in order to determine player reflecting area and
	prevent false scoring when ball speed is too high*/
	var lineA, lineB line
	playerEdge := p.pos
	if playerEdge[0] > 0 {
		playerEdge[0] -= p.width

	} else {
		playerEdge[0] += p.width
	}
	playerEdge[1] += p.length
	if b.velocity[0] != 0 {
		lineA.a = b.velocity[1] / b.velocity[0]
		lineA.b = -(lineA.a)*playerEdge[0] + playerEdge[1]
	} else {
		lineA.a = 0
		lineA.b = playerEdge[1]
	}
	playerEdge[1] -= 2 * p.length
	if b.velocity[0] != 0 {
		lineB.a = b.velocity[1] / b.velocity[0]
		lineB.b = -(lineB.a)*playerEdge[0] + playerEdge[1]
	} else {
		lineB.a = 0
		lineB.b = playerEdge[1]
	}
	ballBottom := [2]float64{b.pos[0], b.pos[1] - b.radius}
	ballTop := [2]float64{b.pos[0], b.pos[1] + b.radius}
	//checks if ball is inside the player reflecting area
	if ballBottom[1] <= lineA.a*ballBottom[0]+lineA.b &&
		ballTop[1] >= lineB.a*ballTop[0]+lineB.b {
		if playerEdge[0] > 0 && ballSidePoint[0] >= playerEdge[0] && b.velocity[0] > 0 {
			b.velocity[0] *= -1 * reflectionGain
			b.velocity[1] *= reflectionGain
		}
		if playerEdge[0] < 0 && ballSidePoint[0] <= playerEdge[0] && b.velocity[0] < 0 {
			b.velocity[0] *= -1 * reflectionGain
			b.velocity[1] *= reflectionGain
		}
	}
}
