package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/phrozen/game-of-life/pkg/conway"
)

type Game struct {
	life *conway.GameOfLife
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (game *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		os.Exit(0)
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		game.life.Fill(25)
	} else {
		game.life.Step()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (game *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for i, cell := range game.life.State() {
		if cell {
			screen.Set(i%game.life.Width(), i/game.life.Width(), color.White)
		}
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (game *Game) Layout(width, height int) (int, int) {
	return game.life.Width(), game.life.Height()
}

func main() {
	width, height := ebiten.ScreenSizeInFullscreen()
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println(width, "x", height)

	game := &Game{life: conway.NewGameOfLife(width/10, height/10)}

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Game of Life")
	ebiten.SetMaxTPS(10)
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
