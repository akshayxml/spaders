package sprites

import (
	"github.com/akshayxml/spaders/models"
	"image/color"
)

func GetPlayerRectangles() []models.Rectangle {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}
	return []models.Rectangle{
		{Position: models.Position{X: 0, Y: 8}, Width: 4, Height: 8, Color: neonGreen},
		{Position: models.Position{X: 4, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 8, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 12, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 16, Y: 0}, Width: 4, Height: 16, Color: neonGreen},
		{Position: models.Position{X: 20, Y: 0}, Width: 4, Height: 16, Color: neonGreen},
		{Position: models.Position{X: 24, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 28, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 32, Y: 4}, Width: 4, Height: 12, Color: neonGreen},
		{Position: models.Position{X: 36, Y: 8}, Width: 4, Height: 8, Color: neonGreen},
	}
}
