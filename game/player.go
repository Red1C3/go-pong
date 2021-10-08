package game

type player struct {
	pos           [2]float64
	length, width float64
	score         int
}

func newPlayer(x float64, l float64, w float64) player {
	p := player{pos: [2]float64{x, 0}, length: l, width: w, score: 0}
	return p
}
func (p *player) move(s float64, dt float64) {
	if p.pos[1]+p.length > 18 && s > 0 {
		return
	}
	if p.pos[1]-p.length < -18 && s < 0 {
		return
	}
	p.pos[1] += s * dt
}
