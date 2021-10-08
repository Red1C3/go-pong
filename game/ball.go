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
	deg := rand.Float64()*90.0 - 45.0
	if rand.Float32() < 0.5 {
		deg += 180
	}
	deg = deg * math.Pi / 180
	intialSpeed := rand.Float64()*25 + 75
	b.velocity[0] = math.Cos(deg) * intialSpeed
	b.velocity[1] = math.Sin(deg) * intialSpeed
	b.radius = 0.5
	return b
}
func (b *ball) update(dt float64, p []player) {
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
	b.pos[0] += dt * b.velocity[0]
	b.pos[1] += dt * b.velocity[1]
	b.resolveCollisions(p, dt)

}
func (b *ball) resolveCollisions(p []player, dt float64) {
	if b.pos[1]+b.radius > 18 && b.velocity[1] > 0 {
		b.velocity[1] *= -1
	}
	if b.pos[1]-b.radius < -18 && b.velocity[1] < 0 {
		b.velocity[1] *= -1
	}
	b.resolvePlayer(p[0], dt)
	b.resolvePlayer(p[1], dt)
}
func (b *ball) resolvePlayer(p player, dt float64) {
	var testPoint [2]float64
	if p.pos[0] > 0 {
		testPoint[0] = b.pos[0] + b.radius
	} else {
		testPoint[0] = b.pos[0] - b.radius
	}
	testPoint[1] = b.pos[1]
	var lineA, lineB line
	point := p.pos
	if point[0] > 0 {
		point[0] -= p.width

	} else {
		point[0] += p.width
	}
	point[1] += p.length
	if b.velocity[0] != 0 {
		lineA.a = b.velocity[1] / b.velocity[0]
		lineA.b = -(lineA.a)*point[0] + point[1]
	} else {
		lineA.a = 0
		lineA.b = point[1]
	}
	point[1] -= 2 * p.length
	if b.velocity[0] != 0 {
		lineB.a = b.velocity[1] / b.velocity[0]
		lineB.b = -(lineB.a)*point[0] + point[1]
	} else {
		lineB.a = 0
		lineB.b = point[1]
	}
	bottomPoint := [2]float64{b.pos[0], b.pos[1] - b.radius}
	topPoint := [2]float64{b.pos[0], b.pos[1] + b.radius}
	if bottomPoint[1] <= lineA.a*bottomPoint[0]+lineA.b &&
		topPoint[1] >= lineB.a*topPoint[0]+lineB.b {
		if point[0] > 0 && testPoint[0] >= point[0] && b.velocity[0] > 0 {
			b.velocity[0] *= -1
		}
		if point[0] < 0 && testPoint[0] <= point[0] && b.velocity[0] < 0 {
			b.velocity[0] *= -1
		}
	}
}
