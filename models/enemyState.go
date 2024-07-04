package models

type EnemyState struct {
	EnemyCount          int
	HorizontalDirection int
	HorizontalSpeed     float64
	BulletCount         int
	EnemyFireRate       int
	EnemyBullets        []Bullet
}
