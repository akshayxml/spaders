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
