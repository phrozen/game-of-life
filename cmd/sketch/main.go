package main

import (
	"image"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const size = 2048
const totalCycleCount = 5000

var DefaultUserParams = UserParams{
	StrokeRatio:              0.75,
	InitialAlpha:             0.1,
	StrokeReduction:          0.002,
	AlphaIncrease:            0.06,
	StrokeInversionThreshold: 0.00,
	StrokeJitter:             size / 10,
	RotationRange:            4.5,
	Seed:                     time.Now().Unix(),
}

type Game struct {
	sketch *Sketch
	img    *image.RGBA
}

func (game *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && runtime.GOOS != "js" {
		os.Exit(0)
	}
	game.sketch.Update()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	pixels := make([]byte, len(game.img.Pix))
	copy(pixels, game.img.Pix)
	screen.ReplacePixels(pixels)
}

func (game *Game) Layout(width, height int) (int, int) {
	return game.sketch.width, game.sketch.height
}

func LoadRandomUnsplashImage() (image.Image, error) {
	res, err := http.Get("https://source.unsplash.com/random/2048")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}

func main() {
	img, err := LoadRandomUnsplashImage()
	if err != nil {
		log.Panicln(err)
	}
	sketch := NewSketch(img, DefaultUserParams)
	game := &Game{sketch: sketch, img: sketch.Output().(*image.RGBA)}

	width, height := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Sketch")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
