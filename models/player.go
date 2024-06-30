package models

type Player struct {
	Position Position
}

func (p *Player) MoveLeft() {
	if p.Position.X >= 50 {
		p.Position.X--
	}
}

func (p *Player) MoveRight(rightBoundary float64) {
	if p.Position.X <= rightBoundary {
		p.Position.X++
	}
}
