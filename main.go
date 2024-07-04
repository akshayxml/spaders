package main

import (
	"fmt"
	"github.com/akshayxml/spaders/models"
	"github.com/akshayxml/spaders/models/EntityState"
	"github.com/akshayxml/spaders/models/Screen"
	"github.com/akshayxml/spaders/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	_ "image/jpeg"
	"io"
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
	audioContext    *audio.Context
	player          *audio.Player
)

const (
	bgImgLocation     string  = "./assets/bg.jpg"
	fontLocation      string  = "./assets/CosmicAlien.ttf"
	enemyImg1Location string  = "./assets/enemyOne.png"
	enemyImg2Location string  = "./assets/enemyTwo.png"
	enemyImg3Location string  = "./assets/enemyThree.png"
	bgAudioLocation   string  = "./assets/audio.mp3"
	normalFontSize    float64 = 18
	bigFontSize       float64 = 36
	windowWidth       float64 = 640
	windowHeight      float64 = 480
	leftBoundary              = 50
	rightBoundary             = windowWidth - 80
	sampleRate                = 44100
)

type Game struct {
	player           *models.Player
	enemies          []models.Enemy
	score            int
	enemyBullets     []models.Bullet
	enemyBulletCount int
	enemyCount       int
	bunkerSprites    []models.Rectangle
	screen           Screen.Screen
	playStartTime    int64
	enemyFireRate    int
	difficulty       int
}

func (g *Game) moveEnemySideways() {
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
			g.enemies[i].Position.X += g.enemies[i].HorizontalSpeed * float64(g.enemies[i].HorizontalDirection)
		}
	}
}

func (g *Game) moveBullets() {
	if g.player.Bullet.IsActive {
		g.player.Bullet.Position.Y += float64(g.player.Bullet.Speed * g.player.Bullet.Direction)
	}

	for i := 0; i < g.enemyBulletCount; i++ {
		g.enemyBullets[i].Position.Y += float64(g.enemyBullets[i].Speed * g.enemyBullets[i].Direction)
	}
}

func (g *Game) addEnemyBullet(bullet models.Bullet) {
	if g.enemyBulletCount < len(g.enemyBullets) {
		g.enemyBullets[g.enemyBulletCount] = bullet
	} else {
		g.enemyBullets = append(g.enemyBullets, bullet)
	}
	g.enemyBulletCount++
}

func (g *Game) removeEnemyBullet(bulletIndex int) {
	g.enemyBullets[bulletIndex], g.enemyBullets[g.enemyBulletCount-1] = g.enemyBullets[g.enemyBulletCount-1], g.enemyBullets[bulletIndex]
	g.enemyBulletCount--
}

func (g *Game) generateEnemyBullets() {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) <= g.enemyFireRate {
		var enemyNumber = rand.Intn(len(g.enemies))
		if g.enemies[enemyNumber].State == EntityState.Alive {
			var enemyWidth = g.enemies[enemyNumber].GetEnemyWidth()
			var enemyHeight = g.enemies[enemyNumber].GetEnemyHeight()
			var bullet = models.Bullet{models.Position{X: g.enemies[enemyNumber].Position.X + enemyWidth/2,
				Y: g.enemies[enemyNumber].Position.Y + enemyHeight/2},
				1, 1, true, getSpritesHeight(sprites.GetEnemyBulletRectangles())}
			g.addEnemyBullet(bullet)
			bullet.Fire()
		}
	}
}

func (g *Game) updateDifficulty(currentTimestamp int64) {
	var elapsedTime = currentTimestamp - g.playStartTime
	var baseHorizontalSpeedLimit = 2.0
	var baseHorizontalSpeedChangeIntervalMs = 10000
	var baseVerticalMoveIntervalMs = 15000
	var baseFireRateLimit = 10
	var baseFireRateLimitChangeIntervalMs = 10000

	for i, _ := range g.enemies {
		var horizontalSpeedChangeIntervalMs = int64(baseHorizontalSpeedChangeIntervalMs - ((baseHorizontalSpeedChangeIntervalMs / 3) * (g.difficulty - 1)))
		var horizontalSpeedLimit = baseHorizontalSpeedLimit + float64(g.difficulty/2)
		g.enemies[i].HorizontalSpeed = min(horizontalSpeedLimit, 1+float64(elapsedTime)/float64(horizontalSpeedChangeIntervalMs*10))

		var verticalMoveIntervalMs = int64(baseVerticalMoveIntervalMs - ((baseVerticalMoveIntervalMs / 3) * (g.difficulty - 1)))
		if elapsedTime >= verticalMoveIntervalMs && elapsedTime%verticalMoveIntervalMs <= 100 {
			g.enemies[i].Position.Y++
		}
	}

	var fireRateLimitChangeIntervalMs = int64(baseFireRateLimitChangeIntervalMs - ((baseFireRateLimitChangeIntervalMs / 3) * (g.difficulty - 1)))
	g.enemyFireRate = min(baseFireRateLimit*g.difficulty, int(elapsedTime/fireRateLimitChangeIntervalMs))
}

func (g *Game) renderScore(screen *ebiten.Image, neonGreen color.RGBA) {

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
}

func (g *Game) renderLives(screen *ebiten.Image) {
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(400, 13)
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, "LIVES", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, textOp)

	playerImgPositions := []struct{ x, y float64 }{
		{480, 10},
		{530, 10},
		{580, 10},
	}

	for i := range playerImgPositions {
		if i < g.player.Lives {
			for _, rect := range sprites.GetPlayerRectangles() {
				ebitenutil.DrawRect(screen, rect.Position.X+playerImgPositions[i].x, rect.Position.Y+playerImgPositions[i].y,
					rect.Width, rect.Height, rect.Color)
			}
		}
	}
}

func (g *Game) renderPlayer(screen *ebiten.Image) {
	for _, rect := range sprites.GetPlayerRectangles() {
		ebitenutil.DrawRect(screen, rect.Position.X+g.player.Position.X, rect.Position.Y+g.player.Position.Y,
			rect.Width, rect.Height, rect.Color)
	}
}

func (g *Game) renderBunker(screen *ebiten.Image) {
	for _, bunkerSprite := range g.bunkerSprites {
		ebitenutil.DrawRect(screen, bunkerSprite.Position.X, bunkerSprite.Position.Y,
			bunkerSprite.Width, bunkerSprite.Height, bunkerSprite.Color)
	}
}

func (g *Game) renderBullets(screen *ebiten.Image) {
	if g.player.Bullet.IsActive {
		for _, rect := range sprites.GetPlayerBulletRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+g.player.Bullet.Position.X, rect.Position.Y+g.player.Bullet.Position.Y,
				rect.Width, rect.Height, rect.Color)
		}
	}
	for i := 0; i < g.enemyBulletCount; i++ {
		for _, rect := range sprites.GetEnemyBulletRectangles() {
			ebitenutil.DrawRect(screen, rect.Position.X+g.enemyBullets[i].Position.X, rect.Position.Y+g.enemyBullets[i].Position.Y,
				rect.Width, rect.Height, rect.Color)
		}
	}
}

func getSpritesWidth(sprites []models.Rectangle) float64 {
	var width = 0.0
	for _, sprite := range sprites {
		width += sprite.Width
	}
	return width
}

func getSpritesHeight(sprites []models.Rectangle) float64 {
	var height = 0.0
	for _, sprite := range sprites {
		height += sprite.Width
	}
	return height
}

func getEnemies(rows int, yPosStart float64, img *ebiten.Image, scale float64) (float64, []models.Enemy) {
	var cols = 10
	var imgWidth = float64(img.Bounds().Dx()) * scale
	var imgHeight = float64(img.Bounds().Dy()) * scale
	var xPosStart = float64(windowWidth/2) - (imgWidth / 2) - 40 - imgWidth*5
	var xGap = imgWidth + 10
	var yGap = imgHeight + 10
	var enemies = []models.Enemy{}

	for y, rowCnt := float64(yPosStart), 0; y < float64(windowHeight) && rowCnt < rows; y, rowCnt = y+yGap, rowCnt+1 {
		for x, colCnt := float64(xPosStart), 0; x < float64(windowWidth) && colCnt < cols; x, colCnt = x+xGap, colCnt+1 {
			var enemy = models.Enemy{models.Position{x, y}, img, scale, EntityState.Alive, 1, 1.0}
			enemies = append(enemies, enemy)
		}
	}

	return yGap * float64(rows), enemies
}

func (g *Game) renderEnemies(screen *ebiten.Image) {
	for _, enemy := range g.enemies {
		if enemy.State == EntityState.Alive {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(enemy.Scale, enemy.Scale)
			opts.GeoM.Translate(enemy.Position.X, enemy.Position.Y)
			screen.DrawImage(enemy.Img, opts)
		}
	}
}

func setupEnemies() []models.Enemy {
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

	return allEnemies
}

func setupBunkers() []models.Rectangle {
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

	return bunkerSprites
}

func (g *Game) detectCollision() {
	if g.player.Bullet.IsActive {
		for i := range g.bunkerSprites {
			if g.bunkerSprites[i].Height > 0 {
				var bunkerSpriteLeft = g.bunkerSprites[i].Position.X
				var bunkerSpriteRight = g.bunkerSprites[i].Position.X + g.bunkerSprites[i].Width
				var bunkerSpriteTop = g.bunkerSprites[i].Position.Y
				var bunkerSpriteBottom = g.bunkerSprites[i].Position.Y + g.bunkerSprites[i].Height
				if g.player.Bullet.HasCollided(bunkerSpriteLeft, bunkerSpriteRight, bunkerSpriteTop, bunkerSpriteBottom) {
					g.bunkerSprites[i].Height -= sprites.GetBunkerRectangles()[0].Height
					g.player.Bullet.IsActive = false
				}
			}
		}

		for i, enemy := range g.enemies {
			if enemy.State == EntityState.Alive {
				var enemyLeftEdge = enemy.Position.X
				var enemyRightEdge = enemy.Position.X + enemy.GetEnemyWidth()
				var enemyTopEdge = enemy.Position.Y
				var enemyBottomEdge = enemy.Position.Y + enemy.GetEnemyHeight()
				if g.player.Bullet.HasCollided(enemyLeftEdge, enemyRightEdge, enemyTopEdge, enemyBottomEdge) {
					g.player.Bullet.IsActive = false
					g.enemies[i].State = EntityState.Dead
					g.score += 5
					g.enemyCount--
					if g.enemyCount == 0 {
						g.screen = Screen.GameOver
					}
				}
			}
		}

		if g.player.Bullet.Position.Y <= 5 {
			g.player.Bullet.IsActive = false
		}
	}

	for i := 0; i < g.enemyBulletCount; i++ {
		if g.enemyBullets[i].IsActive {
			for j := range g.bunkerSprites {
				if g.bunkerSprites[j].Height > 0 {
					var bunkerSpriteLeft = g.bunkerSprites[j].Position.X
					var bunkerSpriteRight = g.bunkerSprites[j].Position.X + g.bunkerSprites[j].Width
					var bunkerSpriteTop = g.bunkerSprites[j].Position.Y
					var bunkerSpriteBottom = g.bunkerSprites[j].Position.Y + g.bunkerSprites[j].Height
					if g.enemyBullets[i].HasCollided(bunkerSpriteLeft, bunkerSpriteRight, bunkerSpriteTop, bunkerSpriteBottom) {
						g.bunkerSprites[j].Height -= sprites.GetBunkerRectangles()[0].Height
						g.bunkerSprites[j].Position.Y += sprites.GetBunkerRectangles()[0].Height
						g.enemyBullets[i].IsActive = false
					}
				}
			}

			for _, playerSprite := range sprites.GetPlayerRectangles() {
				var playerLeftEdge = g.player.Position.X + playerSprite.Position.X
				var playerRightEdge = g.player.Position.X + playerSprite.Position.X + playerSprite.Width
				var playerTopEdge = g.player.Position.Y + playerSprite.Position.Y
				var playerBottomEdge = g.player.Position.Y + playerSprite.Position.Y + playerSprite.Height
				if g.enemyBullets[i].HasCollided(playerLeftEdge, playerRightEdge, playerTopEdge, playerBottomEdge) {
					g.player.Lives--
					g.enemyBullets[i].IsActive = false
					if g.player.Lives == 0 {
						g.screen = Screen.GameOver
					}
				}
			}

			if g.player.Bullet.IsActive {
				if g.player.Bullet.HasCollidedBullets(g.enemyBullets[i]) {
					g.player.Bullet.IsActive = false
					g.enemyBullets[i].IsActive = false
					g.score += 3
				}
			}

			if g.enemyBullets[i].Position.Y >= windowHeight-20 {
				g.enemyBullets[i].IsActive = false
			}
		}
	}

	var i = 0
	for i < g.enemyBulletCount {
		if !g.enemyBullets[i].IsActive {
			g.removeEnemyBullet(i)
		} else {
			i++
		}
	}
}

func (g *Game) DrawMenu(screen *ebiten.Image, neonGreen color.RGBA) {
	msg := "SPADERS"
	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   bigFontSize,
	}
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(windowWidth/2), float64(windowHeight/2))
	textOp.ColorScale.ScaleWithColor(neonGreen)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, msg, face, textOp)

	msg = "PRESS SPACE TO START"
	face = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}
	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(float64(windowWidth/2), float64(windowHeight/2+40))
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, msg, face, textOp)
}

func (g *Game) DrawGameOver(screen *ebiten.Image, neonGreen color.RGBA) {
	msg := "GAME OVER"
	if g.enemyCount == 0 {
		msg = "YOU'VE WON!!!"
	}
	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   bigFontSize,
	}
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(windowWidth/2), float64(windowHeight/2))
	textOp.ColorScale.ScaleWithColor(neonGreen)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, msg, face, textOp)

	msg = "YOUR SCORE IS " + strconv.Itoa(g.score)
	face = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}
	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(float64(windowWidth/2), float64(windowHeight/2+40))
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, msg, face, textOp)

	msg = "PRESS SPACE TO REPLAY"
	face = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}
	textOp = &text.DrawOptions{}
	textOp.GeoM.Translate(float64(windowWidth/2), float64(windowHeight/2+80))
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, msg, face, textOp)
}

func (g *Game) reset() {
	g.screen = Screen.Play
	g.playStartTime = time.Now().UnixMilli()
	g.score = 0
	g.bunkerSprites = setupBunkers()
	g.enemies = setupEnemies()
	g.enemyBulletCount = 0
	g.enemyFireRate = 1
	g.enemyCount = len(g.enemies)
	g.player = &models.Player{
		Position: models.Position{
			X: (windowWidth / 2),
			Y: windowHeight - 40,
		},
		Lives: 3,
		Speed: 2.0,
		Bullet: models.Bullet{
			Direction: -1,
			Speed:     3,
			IsActive:  false,
		},
	}
	g.difficulty = 3
}

func (g *Game) Update() error {
	if g.screen == Screen.Menu || g.screen == Screen.GameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.reset()
		}
	} else if g.screen == Screen.Play {
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.player.MoveLeft()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.player.MoveRight(rightBoundary)
		}
		currentTimestamp := time.Now().UnixMilli()
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.player.Bullet.IsActive && currentTimestamp > g.playStartTime+500 {
			g.player.Bullet.Position = models.Position{g.player.Position.X + 20, g.player.Position.Y}
			g.player.Bullet.Height = getSpritesHeight(sprites.GetPlayerBulletRectangles())
			g.player.Bullet.Fire()
		}

		g.moveEnemySideways()
		g.generateEnemyBullets()
		g.moveBullets()
		g.updateDifficulty(currentTimestamp)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var neonGreen = color.RGBA{0x39, 0xFF, 0x14, 0xFF}

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(bgImg, imgOp)

	g.renderScore(screen, neonGreen)
	g.renderLives(screen)

	if g.screen == Screen.Menu {
		g.DrawMenu(screen, neonGreen)
	} else if g.screen == Screen.GameOver {
		g.DrawGameOver(screen, neonGreen)
		return
	} else {
		g.detectCollision()
		g.renderBullets(screen)
		g.renderEnemies(screen)
		g.renderPlayer(screen)
		g.renderBunker(screen)

		vector.StrokeLine(screen, leftBoundary, float32(windowHeight-10),
			float32(windowWidth-50), float32(windowHeight-10), 2, neonGreen, true)
	}
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

func playAudio(f io.Reader) error {
	audioContext = audio.NewContext(sampleRate)

	d, err := mp3.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		return err
	}

	infiniteLoop := audio.NewInfiniteLoop(d, d.Length())
	player, err := audioContext.NewPlayer(infiniteLoop)
	if err != nil {
		return err
	}

	if err != nil {
		fmt.Println("agb")
		log.Fatal(err)
	}
	player.Play()

	return nil
}

func main() {
	fmt.Println("SPADERS")
	ebiten.SetWindowSize(int(windowWidth), int(windowHeight))
	ebiten.SetWindowTitle("Spaders")

	f, err := os.Open(bgAudioLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = playAudio(f)
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{}
	g.player = &models.Player{
		Lives: 3,
	}
	g.screen = Screen.Menu
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
