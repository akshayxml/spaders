package sprites

import (
	"github.com/akshayxml/spaders/models"
	"image/color"
)

func GetPlayerBulletRectangles() []models.Rectangle {
	var white = color.White
	return []models.Rectangle{
		{Position: models.Position{X: 0, Y: 0}, Width: 2, Height: 2, Color: white},
		{Position: models.Position{X: 0, Y: 3}, Width: 2, Height: 9, Color: white},
	}
}

func GetEnemyBulletRectangles() []models.Rectangle {
	var white = color.White
	return []models.Rectangle{
		{Position: models.Position{X: 2, Y: 0}, Width: 2, Height: 12, Color: white},
		{Position: models.Position{X: 0, Y: 0}, Width: 6, Height: 2, Color: white},
	}
}
