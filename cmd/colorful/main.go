package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/phrozen/game-of-life/pkg/conway"
	"gopkg.in/fogleman/gg.v1"
)

const (
	RESOLUTION = 8
	HALF       = RESOLUTION / 2
)

type Game struct {
	life    *conway.GameOfLife
	img     *ebiten.Image
	context *gg.Context
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (game *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && runtime.GOOS != "js" {
		os.Exit(0)
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ebiten.SetMaxTPS(60)
		x, y := ebiten.CursorPosition()
		game.life.Seed(x/RESOLUTION, y/RESOLUTION)
	} else {
		ebiten.SetMaxTPS(10)
		game.life.Step()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (game *Game) Draw(screen *ebiten.Image) {
	game.context.SetColor(color.Black)
	game.context.Clear()
	x, y := 0.0, 0.0
	for i, cell := range game.life.State() {
		if cell {
			x = float64((i%game.life.Width())*RESOLUTION + HALF)
			y = float64((i/game.life.Width())*RESOLUTION + HALF)
			game.context.DrawCircle(x, y, HALF)
			game.context.SetColor(colorAt(i, game.life.Width()*game.life.Height()))
			game.context.Stroke()
		}
	}
	game.img = ebiten.NewImageFromImage(game.context.Image())
	screen.DrawImage(game.img, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%3.2f FPS", ebiten.CurrentTPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (game *Game) Layout(width, height int) (int, int) {
	return game.life.Width() * RESOLUTION, game.life.Height() * RESOLUTION
}

func main() {
	fmt.Println("OS:", runtime.GOOS)
	width, height := ebiten.ScreenSizeInFullscreen()

	game := &Game{
		life:    conway.NewGameOfLife(width/RESOLUTION, height/RESOLUTION),
		context: gg.NewContext(width, height),
	}

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Game of Life")

	if runtime.GOOS != "js" {
		ebiten.SetFullscreen(true)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func interpolate(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*((value-low1)/(high1-low1))
}

func colorAt(i, max int) color.Color {
	hue := interpolate(float64(i), 0, float64(max), 0, 360)
	return colorful.Hsv(hue, 1.0, 1.0)
}
