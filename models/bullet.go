package models

import "math"

type Bullet struct {
	Position  Position
	Direction int
	Speed     int
	IsActive  bool
	Height    float64
}

func (b *Bullet) Fire() {
	if b.IsActive == false {
		b.IsActive = true
	}
}

func (b *Bullet) HasCollided(leftEdge, rightEdge, topEdge, bottomEdge float64) bool {
	var bulletPositionY = b.Position.Y
	if b.Direction > 0 {
		bulletPositionY += b.Height
	}
	return b.Position.X >= leftEdge && b.Position.X <= rightEdge &&
		bulletPositionY <= bottomEdge && bulletPositionY >= topEdge
}

func (b *Bullet) HasCollidedBullets(otherBullet Bullet) bool {
	var playerBulletX = b.Position.X
	var enemyBulletX = otherBullet.Position.X
	var playerBulletTopEdge = int(b.Position.Y)
	var enemyBulletTopEdge = int(otherBullet.Position.Y)
	var enemyBulletBottomEdge = int(otherBullet.Position.Y + otherBullet.Height)
	return math.Abs(playerBulletX-enemyBulletX) <= 4 && playerBulletTopEdge <= enemyBulletBottomEdge && playerBulletTopEdge >= enemyBulletTopEdge
}
