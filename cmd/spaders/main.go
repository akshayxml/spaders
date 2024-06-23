package main

import (
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
	baseImg         *ebiten.Image
)

const (
	bgImgLocation     string  = "./assets/bg.jpg"
	fontLocation      string  = "./assets/CosmicAlien.ttf"
	enemyImg1Location string  = "./assets/enemyOne.png"
	enemyImg2Location string  = "./assets/enemyTwo.png"
	enemyImg3Location string  = "./assets/enemyThree.png"
	playerImgLocation string  = "./assets/player.png"
	baseImgLocation   string  = "./assets/base.png"
	normalFontSize    float64 = 18
	bigFontSize       float64 = 48
	windowWidth       int     = 640
	windowHeight      int     = 480
)

type Game struct{}

type Enemy struct {
	rows  int
	yPos  float64
	img   *ebiten.Image
	scale float64
}

func (g *Game) Update() error {
	return nil
}

func renderPlayerImages(screen *ebiten.Image) {
	playerImgPositions := []struct{ x, y float64 }{
		{float64(windowWidth / 2), float64(windowHeight - 40)},
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

func renderBaseImages(screen *ebiten.Image) {
	playerImgPositions := []struct{ x, y float64 }{
		{float64(148), float64(windowHeight - 100)},
		{float64(256), float64(windowHeight - 100)},
		{float64(364), float64(windowHeight - 100)},
		{float64(472), float64(windowHeight - 100)},
	}

	for _, pos := range playerImgPositions {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(pos.x, pos.y)
		screen.DrawImage(baseImg, opts)
	}
}

func renderEnemies(screen *ebiten.Image, enemy Enemy) float64 {
	var rows = enemy.rows
	var cols = 10
	var img = enemy.img
	var scale float64 = enemy.scale
	var imgWidth float64 = float64(img.Bounds().Dx()) * scale
	var imgHeight float64 = float64(img.Bounds().Dy()) * scale
	var yPosStart = enemy.yPos
	var xPosStart float64 = float64(windowWidth/2) - (imgWidth / 2) - 40 - imgWidth*5
	var xGap float64 = imgWidth + 10
	var yGap float64 = imgHeight + 10

	for y, rowCnt := float64(yPosStart), 0; y < float64(windowHeight) && rowCnt < rows; y, rowCnt = y+yGap, rowCnt+1 {
		for x, colCnt := float64(xPosStart), 0; x < float64(windowWidth) && colCnt < cols; x, colCnt = x+xGap, colCnt+1 {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(scale, scale)
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, opts)
		}
	}

	return yGap * float64(rows)
}

func (g *Game) Draw(screen *ebiten.Image) {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}
	var enemyYPos float64 = 60

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(bgImg, imgOp)

	renderPlayerImages(screen)

	renderBaseImages(screen)

	var enemyThree = Enemy{1, enemyYPos, enemyThreeImg, 0.5}
	var yGap = renderEnemies(screen, enemyThree)

	enemyYPos += yGap
	var enemyTwo = Enemy{2, enemyYPos, enemyTwoImg, 0.6}
	yGap = renderEnemies(screen, enemyTwo)

	enemyYPos += yGap
	var enemyOne = Enemy{2, enemyYPos, enemyOneImg, 0.7}
	yGap = renderEnemies(screen, enemyOne)

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
	return windowWidth, windowHeight
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

	baseImg, _, err = ebitenutil.NewImageFromFile(baseImgLocation)
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
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Spaders")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
