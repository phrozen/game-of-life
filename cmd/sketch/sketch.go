package main

import (
	"image"
	"image/color"
	"math/rand"

	"gopkg.in/fogleman/gg.v1"
)

type UserParams struct {
	StrokeRatio              float64
	DestWidth                int
	DestHeight               int
	InitialAlpha             float64
	StrokeReduction          float64
	AlphaIncrease            float64
	StrokeInversionThreshold float64
	StrokeJitter             int
	MinEdgeCount             int
	MaxEdgeCount             int
	RotationRange            float64
	Seed                     int64
}

type Sketch struct {
	UserParams        // embed for easier access
	source            image.Image
	dc                *gg.Context
	width             int
	height            int
	strokeSize        float64
	initialStrokeSize float64
	random            *rand.Rand
}

func NewSketch(source image.Image, userParams UserParams) *Sketch {
	s := &Sketch{UserParams: userParams}
	s.random = rand.New(rand.NewSource(userParams.Seed))
	bounds := source.Bounds()
	s.width, s.height = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.width)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.width, s.height)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.width), float64(s.height))
	canvas.FillPreserve()
	canvas.DrawImage(source, 0, 0)

	s.source = source
	s.dc = canvas
	return s
}

func (s *Sketch) Update() {
	x := s.random.Intn(s.width)
	y := s.random.Intn(s.height)
	r, g, b := rgb255(s.source.At(x, y))

	x += s.Range(s.StrokeJitter)
	y += s.Range(s.StrokeJitter)
	s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
	s.dc.DrawRegularPolygon(4, float64(x), float64(y), s.strokeSize, gg.Radians(s.RangeF(s.RotationRange)))
	s.dc.FillPreserve()

	if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
		if (r+g+b)/3 < 128 {
			s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
		} else {
			s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
		}
	}
	s.dc.Stroke()

	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease
}

func (s *Sketch) Output() image.Image {
	return s.dc.Image()
}

func (s *Sketch) Range(max int) int {
	return -max + s.random.Intn(2*max)
}

func (s *Sketch) RangeF(max float64) float64 {
	a := -max + (s.random.Float64() * 2 * max)
	if a < 0 {
		return 360.0 + a
	}
	return a
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}
