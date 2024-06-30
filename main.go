package main

import (
	"fmt"
	"github.com/akshayxml/spaders/lib/sprites"
	"github.com/akshayxml/spaders/models"
	"github.com/akshayxml/spaders/models/EntityState"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	_ "image/jpeg"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
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
	player        *models.Player
	enemies       []models.Enemy
	playerBullet  *bullet
	score         int
	enemyBullets  []bullet
	bunkerSprites []models.Rectangle
}

type bullet struct {
	position  models.Position
	direction int
	speed     int
	isActive  bool
}

func (b *bullet) fire() {
	if b.isActive == false {
		b.isActive = true
	}
}

func moveEnemySideways(g *Game) {
	leftMostEnemy := windowWidth
	rightMostEnemy := 0.0
	for i := range g.enemies {
		if g.enemies[i].State == EntityState.Alive {
			leftMostEnemy = min(leftMostEnemy, g.enemies[i].Position.X)
			rightMostEnemy = max(rightMostEnemy, g.enemies[i].Position.X)
		}
	}
	moveLeft := rightMostEnemy >= rightBoundary
	moveRight := leftMostEnemy < leftBoundary

	for i := range g.enemies {
		if g.enemies[i].State == EntityState.Alive {
			if moveRight {
				g.enemies[i].HorizontalDirection = 1
			} else if moveLeft {
				g.enemies[i].HorizontalDirection = -1
			}
			g.enemies[i].Position.X += float64(g.enemies[i].HorizontalDirection)
		}
	}
}

func moveBullets(g *Game) {
	if g.playerBullet.isActive {
		g.playerBullet.position.Y += float64(g.playerBullet.speed * g.playerBullet.direction)
	}

	for i := range g.enemyBullets {
		g.enemyBullets[i].position.Y += float64(g.enemyBullets[i].speed * g.enemyBullets[i].direction)
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
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.playerBullet.isActive {
		g.playerBullet.position = models.Position{g.player.Position.X + 20, g.player.Position.Y}
		g.playerBullet.fire()
	}
	moveBullets(g)

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) == 0 {
		var enemyNumber = rand.Intn(len(g.enemies))
		if g.enemies[enemyNumber].State == EntityState.Alive {
			var enemyWidth = g.enemies[enemyNumber].GetEnemyWidth()
			var enemyHeight = g.enemies[enemyNumber].GetEnemyHeight()
			var bullet = bullet{models.Position{X: g.enemies[enemyNumber].Position.X + enemyWidth/2,
				Y: g.enemies[enemyNumber].Position.Y + enemyHeight/2},
				1, 1, true}
			g.enemyBullets = append(g.enemyBullets, bullet)
			bullet.fire()
		}
	}
	return nil
}

func renderPlayerImages(g *Game, screen *ebiten.Image) {
	playerImgPositions := []struct{ x, y float64 }{
		{g.player.Position.X, g.player.Position.Y},
		{480, 10},
		{530, 10},
		{580, 10},
	}

	for i := range playerImgPositions {
		if i <= g.player.Lives {
			for _, rect := range sprites.GetPlayerRectangles() {
				ebitenutil.DrawRect(screen, rect.Position.X+playerImgPositions[i].x, rect.Position.Y+playerImgPositions[i].y,
					rect.Width, rect.Height, rect.Color)
			}
		}
	}
}

func renderBunkerImages(g *Game, screen *ebiten.Image, neonGreen color.RGBA) {
	for _, bunkerSprite := range g.bunkerSprites {
		ebitenutil.DrawRect(screen, bunkerSprite.Position.X, bunkerSprite.Position.Y,
			bunkerSprite.Width, bunkerSprite.Height, bunkerSprite.Color)
	}
}

func renderBullets(g *Game, screen *ebiten.Image) {
	if g.playerBullet.isActive {
		for _, rect := range sprites.GetPlayerBulletRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+g.playerBullet.position.X, rect.Position.Y+g.playerBullet.position.Y,
				rect.Width, rect.Height, rect.Color)
		}
	}
	for _, enemyBullet := range g.enemyBullets {
		for _, rect := range sprites.GetEnemyBulletRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+enemyBullet.position.X, rect.Position.Y+enemyBullet.position.Y,
				rect.Width, rect.Height, rect.Color)
		}
	}
}

func getEnemies(rows int, yPosStart float64, img *ebiten.Image, scale float64) (float64, []models.Enemy) {
	var cols = 10
	var imgWidth float64 = float64(img.Bounds().Dx()) * scale
	var imgHeight float64 = float64(img.Bounds().Dy()) * scale
	var xPosStart float64 = float64(windowWidth/2) - (imgWidth / 2) - 40 - imgWidth*5
	var xGap float64 = imgWidth + 10
	var yGap float64 = imgHeight + 10
	var enemies = []models.Enemy{}

	for y, rowCnt := float64(yPosStart), 0; y < float64(windowHeight) && rowCnt < rows; y, rowCnt = y+yGap, rowCnt+1 {
		for x, colCnt := float64(xPosStart), 0; x < float64(windowWidth) && colCnt < cols; x, colCnt = x+xGap, colCnt+1 {
			var enemy = models.Enemy{models.Position{x, y}, img, scale, EntityState.Alive, 1}
			enemies = append(enemies, enemy)
		}
	}

	return yGap * float64(rows), enemies
}

func renderEnemies(g *Game, screen *ebiten.Image) {
	for _, enemy := range g.enemies {
		if enemy.State == EntityState.Alive {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(enemy.Scale, enemy.Scale)
			opts.GeoM.Translate(enemy.Position.X, enemy.Position.Y)
			screen.DrawImage(enemy.Img, opts)
		}
	}
}

func hasCollided(leftEdge, rightEdge, topEdge, bottomEdge float64, bulletPosition models.Position) bool {
	return bulletPosition.X >= leftEdge && bulletPosition.X <= rightEdge &&
		bulletPosition.Y <= bottomEdge && bulletPosition.Y >= topEdge
}

func detectCollision(g *Game) {
	for i, enemy := range g.enemies {
		if enemy.State == EntityState.Alive && g.playerBullet.isActive {
			var enemyLeftEdge = enemy.Position.X
			var enemyRightEdge = enemy.Position.X + enemy.GetEnemyWidth()
			var enemyTopEdge = enemy.Position.Y
			var enemyBottomEdge = enemy.Position.Y + enemy.GetEnemyHeight()
			if hasCollided(enemyLeftEdge, enemyRightEdge, enemyTopEdge, enemyBottomEdge, g.playerBullet.position) {
				g.playerBullet.isActive = false
				g.enemies[i].State = EntityState.Dead
				g.score++
			}
		}
	}

	for i := range g.bunkerSprites {
		if g.bunkerSprites[i].Width > 0 {
			if g.playerBullet.isActive {
				var bunkerSpriteLeft = g.bunkerSprites[i].Position.X
				var bunkerSpriteRight = g.bunkerSprites[i].Position.X + g.bunkerSprites[i].Width
				var bunkerSpriteTop = g.bunkerSprites[i].Position.Y
				var bunkerSpriteBottom = g.bunkerSprites[i].Position.Y + g.bunkerSprites[i].Height
				//println(int(bunkerSpriteLeft))
				//println(int(bunkerSpriteRight))
				if hasCollided(bunkerSpriteLeft, bunkerSpriteRight, bunkerSpriteTop, bunkerSpriteBottom, g.playerBullet.position) {
					g.bunkerSprites[i].Width = 0
					g.playerBullet.isActive = false
				}
			}
		}
	}

	if g.playerBullet.position.Y <= 5 {
		g.playerBullet.isActive = false
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(bgImg, imgOp)

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
	text.Draw(screen, strconv.Itoa(g.score), &text.GoTextFace{
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

	detectCollision(g)
	renderPlayerImages(g, screen)
	renderEnemies(g, screen)
	renderBullets(g, screen)
	renderBunkerImages(g, screen, neonGreen)

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
	fmt.Println("SPADERS")
	ebiten.SetWindowSize(int(windowWidth), int(windowHeight))
	ebiten.SetWindowTitle("Spaders")

	var allEnemies = []models.Enemy{}
	var enemyYPos float64 = 60
	var yGap, enemies = getEnemies(1, enemyYPos, enemyThreeImg, 0.5)
	allEnemies = append(allEnemies, enemies[:]...)

	enemyYPos += yGap
	yGap, enemies = getEnemies(2, enemyYPos, enemyTwoImg, 0.6)
	allEnemies = append(allEnemies, enemies[:]...)

	enemyYPos += yGap
	yGap, enemies = getEnemies(2, enemyYPos, enemyOneImg, 0.6)
	allEnemies = append(allEnemies, enemies[:]...)

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
	var bunkerSprites = []models.Rectangle{}
	for i := range bunkerPositions {
		for _, bunkerSprite := range sprites.GetBunkerRectangles() {
			bunkerSprite.Position.X += bunkerPositions[i].x
			bunkerSprite.Position.Y += bunkerPositions[i].y
			bunkerSprites = append(bunkerSprites, bunkerSprite)
		}
	}

	g := &Game{
		player: &models.Player{
			Position: models.Position{
				X: float64(windowWidth / 2),
				Y: float64(windowHeight - 40),
			},
			Lives: 3,
			Speed: 1.5,
		},
		playerBullet: &bullet{
			position: models.Position{
				X: -1,
				Y: -1,
			},
			direction: -1,
			speed:     3,
			isActive:  false,
		},
		enemies:       allEnemies,
		score:         0,
		bunkerSprites: bunkerSprites,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
