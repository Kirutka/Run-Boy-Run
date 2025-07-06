package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameObject struct {
	X, Y    float64
	Speed   float64
	Image   *ebiten.Image
	Width   int
	Height  int
	IsRight bool
}

func (g *GameObject) GetRect() image.Rectangle {
	return image.Rect(
		int(g.X), 
		int(g.Y), 
		int(g.X)+g.Width, 
		int(g.Y)+g.Height,
	)
}

func (g *GameObject) Draw(screen *ebiten.Image) {
	if g.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.X, g.Y)
		screen.DrawImage(g.Image, op)
	} else {
		// Fallback to colored rectangle if no image
		vector.DrawFilledRect(screen,
			float32(g.X),
			float32(g.Y),
			float32(g.Width),
			float32(g.Height),
			color.RGBA{255, 0, 0, 255},
			false)
	}
}

func (g *GameObject) Update(elapsed float64, screenWidth float64) {
	if g.IsRight {
		g.X += g.Speed * elapsed * float64(GridSize)
		if g.X > screenWidth {
			g.X = -float64(g.Width)
		}
	} else {
		g.X -= g.Speed * elapsed * float64(GridSize)
		if g.X < -float64(g.Width) {
			g.X = screenWidth
		}
	}
}