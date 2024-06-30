package models

import (
	"github.com/akshayxml/spaders/models/EntityState"
	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Position            Position
	Img                 *ebiten.Image
	Scale               float64
	State               EntityState.EntityState
	HorizontalDirection int
}

func (e *Enemy) GetEnemyWidth() float64 {
	return float64(e.Img.Bounds().Dx()) * e.Scale
}

func (e *Enemy) GetEnemyHeight() float64 {
	return float64(e.Img.Bounds().Dy()) * e.Scale
}
