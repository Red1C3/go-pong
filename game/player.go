package game

type Player struct {
	Pos           [2]float64
	length, width float64
	Score         int
}

func NewPlayer(x float64, l float64, w float64) Player {
	p := Player{Pos: [2]float64{x, 0}, length: l, width: w, Score: 0}
	return p
}
func (p *Player) Move(s float64, dt float64) {
	if p.Pos[1]+p.length > 18 && s > 0 {
		return
	}
	if p.Pos[1]-p.length < -18 && s < 0 {
		return
	}
	p.Pos[1] += s * dt
}
