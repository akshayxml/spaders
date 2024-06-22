package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	_ "image/jpeg"
	"log"
)

var bgImg *ebiten.Image
var (
	mplusFaceSource *text.GoTextFaceSource
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const (
		normalFontSize = 12
		bigFontSize    = 48
	)

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.25, 0.25)
	screen.DrawImage(bgImg, imgOp)

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(20, 0)
	//textOp.GeoM.Scale(0.2, 0.2)
	textOp.ColorScale.ScaleWithColor(color.RGBA{0x00, 0xFF, 0x00, 0xFF})
	text.Draw(screen, "SCORE", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, textOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func init() {
	var err error
	bgImg, _, err = ebitenutil.NewImageFromFile("./assets/bg.jpg")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Spaders")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
