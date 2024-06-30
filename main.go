package main

import (
	"github.com/akshayxml/spaders/lib/sprites"
	"github.com/akshayxml/spaders/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"
	"strconv"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	bgImg           *ebiten.Image
	enemyOneImg     *ebiten.Image
	enemyTwoImg     *ebiten.Image
	enemyThreeImg   *ebiten.Image
)

const (
	bgImgLocation     string  = "./assets/bg.jpg"
	fontLocation      string  = "./assets/CosmicAlien.ttf"
	enemyImg1Location string  = "./assets/enemyOne.png"
	enemyImg2Location string  = "./assets/enemyTwo.png"
	enemyImg3Location string  = "./assets/enemyThree.png"
	normalFontSize    float64 = 18
	bigFontSize       float64 = 48
	windowWidth       float64 = 640
	windowHeight      float64 = 480
	leftBoundary              = 50
	rightBoundary             = windowWidth - 80
)

type Game struct {
	player       *models.Player
	enemies      []Enemy
	playerBullet *bullet
}

type bullet struct {
	position  models.Position
	direction int
	speed     int
}

type Enemy struct {
	x          float64
	y          float64
	img        *ebiten.Image
	scale      float64
	state      EntityState
	xDirection int
}

type EntityState int

const (
	alive EntityState = iota
	dead  EntityState = iota
	dying EntityState = iota
)

func (b *bullet) fire(startPosition models.Position) {
	if b.position.X == -1 {
		b.position = startPosition
	}
}

func moveEnemySideways(g *Game) {
	leftMostEnemy := windowWidth
	rightMostEnemy := 0.0
	for i := range g.enemies {
		leftMostEnemy = min(leftMostEnemy, g.enemies[i].x)
		rightMostEnemy = max(rightMostEnemy, g.enemies[i].x)
	}
	moveLeft := rightMostEnemy >= rightBoundary
	moveRight := leftMostEnemy < leftBoundary

	for i := range g.enemies {
		if moveRight {
			g.enemies[i].xDirection = 1
		} else if moveLeft {
			g.enemies[i].xDirection = -1
		}
		g.enemies[i].x += float64(g.enemies[i].xDirection)
	}
}

func moveBullet(g *Game) {
	if g.playerBullet.position.Y != -1 {
		g.playerBullet.position.Y += float64(g.playerBullet.speed * g.playerBullet.direction)
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.MoveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.MoveRight(rightBoundary)
	}
	moveEnemySideways(g)
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.playerBullet.fire(models.Position{g.player.Position.X + 20, g.player.Position.Y})
	}
	moveBullet(g)
	return nil
}

func renderPlayerImages(g *Game, screen *ebiten.Image) {
	playerImgPositions := []struct{ x, y float64 }{
		{g.player.Position.X, g.player.Position.Y},
		{480, 10},
		{530, 10},
		{580, 10},
	}

	for i, _ := range playerImgPositions {
		for _, rect := range sprites.GetPlayerRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+playerImgPositions[i].x, rect.Position.Y+playerImgPositions[i].y,
				rect.Width, rect.Height, rect.Color)
		}
	}
}

func renderBunkerImages(screen *ebiten.Image, neonGreen color.RGBA) {
	var imgWidth float64 = 0
	for _, rect := range sprites.GetBunkerRectangles() {
		imgWidth += rect.Width
	}
	bunkerPositions := []struct{ x, y float64 }{
		{float64(windowWidth/5 - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(2*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(3*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(4*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
	}

	for i, _ := range bunkerPositions {
		for _, rect := range sprites.GetBunkerRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+bunkerPositions[i].x, rect.Position.Y+bunkerPositions[i].y,
				rect.Width, rect.Height, rect.Color)
		}
	}
}

func renderBullets(g *Game, screen *ebiten.Image) {
	if g.playerBullet.position.X != -1 {
		for _, rect := range sprites.GetPlayerBulletRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+g.playerBullet.position.X, rect.Position.Y+g.playerBullet.position.Y,
				rect.Width, rect.Height, rect.Color)
		}
	}
}

func getEnemies(rows int, yPosStart float64, img *ebiten.Image, scale float64) (float64, []Enemy) {
	var cols = 10
	var imgWidth float64 = float64(img.Bounds().Dx()) * scale
	var imgHeight float64 = float64(img.Bounds().Dy()) * scale
	var xPosStart float64 = float64(windowWidth/2) - (imgWidth / 2) - 40 - imgWidth*5
	var xGap float64 = imgWidth + 10
	var yGap float64 = imgHeight + 10
	var enemies = []Enemy{}

	for y, rowCnt := float64(yPosStart), 0; y < float64(windowHeight) && rowCnt < rows; y, rowCnt = y+yGap, rowCnt+1 {
		for x, colCnt := float64(xPosStart), 0; x < float64(windowWidth) && colCnt < cols; x, colCnt = x+xGap, colCnt+1 {
			var enemy = Enemy{x, y, img, scale, alive, 1}
			enemies = append(enemies, enemy)
		}
	}

	return yGap * float64(rows), enemies
}

func renderEnemies(g *Game, screen *ebiten.Image) {
	for _, enemy := range g.enemies {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(enemy.scale, enemy.scale)
		opts.GeoM.Translate(enemy.x, enemy.y)
		screen.DrawImage(enemy.img, opts)
	}
}

func detectCollision(g *Game) {
	//var enemyCollisionLocations = []models.Position
	//for _, enemy := range g.enemies{
	//	for _,
	//}
}

func (g *Game) Draw(screen *ebiten.Image) {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(bgImg, imgOp)

	detectCollision(g)
	renderPlayerImages(g, screen)
	renderEnemies(g, screen)
	renderBunkerImages(screen, neonGreen)
	renderBullets(g, screen)

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(50, 13)
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, "SCORE", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, textOp)

	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(150, 13)
	textOp.ColorScale.ScaleWithColor(neonGreen)
	text.Draw(screen, strconv.Itoa(0), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, textOp)

	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(400, 13)
	textOp.ColorScale.ScaleWithColor(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	text.Draw(screen, "LIVES", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, textOp)

	vector.StrokeLine(screen, leftBoundary, float32(windowHeight-10),
		float32(windowWidth-50), float32(windowHeight-10), 2, neonGreen, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(windowWidth), int(windowHeight)
}

func init() {
	var err error
	bgImg, _, err = ebitenutil.NewImageFromFile(bgImgLocation)
	if err != nil {
		log.Fatal(err)
	}

	enemyOneImg, _, err = ebitenutil.NewImageFromFile(enemyImg1Location)
	if err != nil {
		log.Fatal(err)
	}

	enemyTwoImg, _, err = ebitenutil.NewImageFromFile(enemyImg2Location)
	if err != nil {
		log.Fatal(err)
	}

	enemyThreeImg, _, err = ebitenutil.NewImageFromFile(enemyImg3Location)
	if err != nil {
		log.Fatal(err)
	}

	textFile, err := os.Open(fontLocation)
	s, err := text.NewGoTextFaceSource(textFile)
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func main() {
	ebiten.SetWindowSize(int(windowWidth), int(windowHeight))
	ebiten.SetWindowTitle("Spaders")

	var allEnemies = []Enemy{}
	var enemyYPos float64 = 60
	var yGap, enemies = getEnemies(1, enemyYPos, enemyThreeImg, 0.5)
	allEnemies = append(allEnemies, enemies[:]...)

	enemyYPos += yGap
	yGap, enemies = getEnemies(2, enemyYPos, enemyTwoImg, 0.6)
	allEnemies = append(allEnemies, enemies[:]...)

	enemyYPos += yGap
	yGap, enemies = getEnemies(2, enemyYPos, enemyOneImg, 0.6)
	allEnemies = append(allEnemies, enemies[:]...)

	g := &Game{
		player: &models.Player{
			Position: models.Position{
				X: float64(windowWidth / 2),
				Y: float64(windowHeight - 40),
			},
		},
		playerBullet: &bullet{
			position: models.Position{
				X: -1,
				Y: -1,
			},
			direction: -1,
			speed:     3,
		},
		enemies: allEnemies,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
