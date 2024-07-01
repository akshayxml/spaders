package models

type Bullet struct {
	Position  Position
	Direction int
	Speed     int
	IsActive  bool
}

func (b *Bullet) Fire() {
	if b.IsActive == false {
		b.IsActive = true
	}
}
