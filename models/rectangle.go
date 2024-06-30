package models

import (
	"image/color"
)

type Rectangle struct {
	Position      Position
	Width, Height float64
	Color         color.Color
}
