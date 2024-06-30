package models

type Player struct {
	Position Position
	Lives    int
	Speed    float64
}

func (p *Player) MoveLeft() {
	if p.Position.X >= 50 {
		p.Position.X -= p.Speed
	}
}

func (p *Player) MoveRight(rightBoundary float64) {
	if p.Position.X <= rightBoundary {
		p.Position.X += p.Speed
	}
}
