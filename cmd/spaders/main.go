package main

import (
	"fmt"
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
	playerImg       *ebiten.Image
	enemyOneImg     *ebiten.Image
	enemyTwoImg     *ebiten.Image
	enemyThreeImg   *ebiten.Image
	bunkerImg       *ebiten.Image
)

const (
	bgImgLocation     string  = "./assets/bg.jpg"
	fontLocation      string  = "./assets/CosmicAlien.ttf"
	enemyImg1Location string  = "./assets/enemyOne.png"
	enemyImg2Location string  = "./assets/enemyTwo.png"
	enemyImg3Location string  = "./assets/enemyThree.png"
	playerImgLocation string  = "./assets/player.png"
	bunkerImgLocation string  = "./assets/bunker.png"
	normalFontSize    float64 = 18
	bigFontSize       float64 = 48
	windowWidth       float64 = 640
	windowHeight      float64 = 480
)

type Game struct {
	player  *player
	enemies []Enemy
}

type player struct {
	x float64
	y float64
}

type Enemy struct {
	x         float64
	y         float64
	img       *ebiten.Image
	scale     float64
	state     EntityState
	moveRight int
}

type EntityState int

const (
	alive EntityState = iota
	dead  EntityState = iota
	dying EntityState = iota
)

func (p *player) moveLeft() {
	if p.x >= 50 {
		p.x--
	}
}

func (p *player) moveRight() {
	if p.x <= windowWidth-80 {
		p.x++
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.moveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.moveRight()
	}
	var leftMostEnemy float64 = windowWidth
	var rightMostEnemy float64 = 0
	for i := range g.enemies {
		leftMostEnemy = min(leftMostEnemy, g.enemies[i].x)
		rightMostEnemy = max(rightMostEnemy, g.enemies[i].x)
	}
	var switchSidewaysMovement = false
	fmt.Print(rightMostEnemy)

	if leftMostEnemy <= 50 || rightMostEnemy >= windowWidth-80 {
		switchSidewaysMovement = true
	}
	for i := range g.enemies {
		if switchSidewaysMovement {
			g.enemies[i].moveRight = 1 - g.enemies[i].moveRight
			fmt.Println(g.enemies[i].moveRight)
		}
		g.enemies[i].x += float64(g.enemies[i].moveRight)
	}
	return nil
}

func (p *player) Position() (float64, float64) {
	return p.x, p.y
}

func renderPlayerImages(g *Game, screen *ebiten.Image) {
	playerImgPositions := []struct{ x, y float64 }{
		{g.player.x, g.player.y},
		{480, 10},
		{520, 10},
		{560, 10},
	}

	for _, pos := range playerImgPositions {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(0.5, 0.5)
		opts.GeoM.Translate(pos.x, pos.y)
		screen.DrawImage(playerImg, opts)
	}
}

func renderBunkerImages(screen *ebiten.Image) {
	var imgWidth = bunkerImg.Bounds().Dy()
	playerImgPositions := []struct{ x, y float64 }{
		{float64(windowWidth/5 - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(2*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(3*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
		{float64(4*(windowWidth/5) - float64(imgWidth/2)), float64(windowHeight - 100)},
	}

	for _, pos := range playerImgPositions {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(pos.x, pos.y)
		screen.DrawImage(bunkerImg, opts)
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

func (g *Game) Draw(screen *ebiten.Image) {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(bgImg, imgOp)

	renderPlayerImages(g, screen)
	renderEnemies(g, screen)
	renderBunkerImages(screen)

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

	vector.StrokeLine(screen, 50, float32(windowHeight-10),
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

	playerImg, _, err = ebitenutil.NewImageFromFile(playerImgLocation)
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

	bunkerImg, _, err = ebitenutil.NewImageFromFile(bunkerImgLocation)
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
	yGap, enemies = getEnemies(2, enemyYPos, enemyOneImg, 0.7)
	allEnemies = append(allEnemies, enemies[:]...)

	g := &Game{
		player: &player{
			x: float64(windowWidth / 2),
			y: float64(windowHeight - 40),
		},
		enemies: allEnemies,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
