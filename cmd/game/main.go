package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 240
	screenHeight = 240

	tileSize = 16
	tileXNum = 32

	// SpriteBack - sprite for background
	SpriteBack = 129

	// SpriteSnake - sprite for snake piece
	SpriteSnake = 126
)

var (
	tilesImage       *ebiten.Image
	snakePosition    [][]int
	currentDirection = ebiten.KeyRight
	isFinished       = true
	pieceToPick      = [2]int{-1, -1}
)

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("assets/sprites.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
}

// Game ...
type Game struct {
	counter int
}

// Update state
func (g *Game) Update(screen *ebiten.Image) error {
	if isFinished {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			initNewGame()
		} else {
			return nil
		}
	}
	g.counter++
	doStep := g.counter%10 == 1
	if doStep {
		const xNum = screenWidth / tileSize
		const yNum = screenHeight / tileSize
		oldFirstSnakePiece := snakePosition[len(snakePosition)-1]
		newFirstSnakePiece := []int{oldFirstSnakePiece[0], oldFirstSnakePiece[1]}
		isFinished = true
		if currentDirection == ebiten.KeyUp {
			if newFirstSnakePiece[1] > 0 {
				newFirstSnakePiece[1]--
				isFinished = false
			}
		} else if currentDirection == ebiten.KeyDown {
			if newFirstSnakePiece[1] < (yNum - 1) {
				newFirstSnakePiece[1]++
				isFinished = false
			}
		} else if currentDirection == ebiten.KeyLeft {
			if newFirstSnakePiece[0] > 0 {
				newFirstSnakePiece[0]--
				isFinished = false
			}
		} else if currentDirection == ebiten.KeyRight {
			if newFirstSnakePiece[0] < (xNum - 1) {
				newFirstSnakePiece[0]++
				isFinished = false
			}
		}
		if !isFinished {
			// Check self-crossing
			for i, piece := range snakePosition {
				if i == 0 {
					continue // Skip first item (it is last snake piece and will be deleted anyway)
				}
				if piece[0] == newFirstSnakePiece[0] && piece[1] == newFirstSnakePiece[1] {
					isFinished = true
					break
				}
			}
		}
		if isFinished {
			return nil
		}
		startIndex := 1
		if newFirstSnakePiece[0] == pieceToPick[0] && newFirstSnakePiece[1] == pieceToPick[1] {
			startIndex = 0
			pieceToPick = getPieceToPick()
		}
		snakePosition = append(snakePosition[startIndex:len(snakePosition)], newFirstSnakePiece)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && currentDirection != ebiten.KeyDown {
		currentDirection = ebiten.KeyUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) && currentDirection != ebiten.KeyUp {
		currentDirection = ebiten.KeyDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && currentDirection != ebiten.KeyRight {
		currentDirection = ebiten.KeyLeft
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) && currentDirection != ebiten.KeyLeft {
		currentDirection = ebiten.KeyRight
	}
	return nil
}

func initNewGame() {
	snakePosition = [][]int{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 0},
	}
	currentDirection = ebiten.KeyRight
	pieceToPick = getPieceToPick()
	isFinished = false
}

func getPieceToPick() [2]int {
	const xMax = screenWidth / tileSize
	const yMax = screenHeight / tileSize
	var newPieceX, newPieceY int
	for {
		newPieceIndex := rand.Intn(xMax * yMax)
		newPieceX = newPieceIndex % xMax
		newPieceY = newPieceIndex / xMax
		isEmptyCell := true
		for _, piece := range snakePosition {
			if piece[0] == newPieceX && piece[1] == newPieceY {
				isEmptyCell = false
				break
			}
		}
		if isEmptyCell {
			break
		}
	}
	return [2]int{newPieceX, newPieceY}
}

// Draw scene
func (g *Game) Draw(screen *ebiten.Image) {
	const xNum = screenWidth / tileSize
	const yNum = screenHeight / tileSize
	// Draw background
	sx := (SpriteBack % tileXNum) * tileSize
	sy := (SpriteBack / tileXNum) * tileSize
	for x := 0; x < xNum; x++ {
		for y := 0; y < yNum; y++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
	// Draw snake
	sx = (SpriteSnake % tileXNum) * tileSize
	sy = (SpriteSnake / tileXNum) * tileSize
	for _, p := range snakePosition {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p[0]*tileSize), float64(p[1]*tileSize))
		screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	}
	// Draw piece to pick
	if pieceToPick[0] != -1 {
		sx = (SpriteSnake % tileXNum) * tileSize
		sy = (SpriteSnake / tileXNum) * tileSize
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pieceToPick[0]*tileSize), float64(pieceToPick[1]*tileSize))
		screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	}
	// Draw finish text
	if isFinished {
		ebitenutil.DebugPrint(screen, "You've lost :( Press R to restart")
	}
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Go Snake")
	initNewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
