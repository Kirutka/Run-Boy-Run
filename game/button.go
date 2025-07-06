package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Button struct {
	X, Y, Width, Height float64
	Text                string
	Hovered             bool
	Action              func()
	Font                font.Face
}

func (b *Button) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

func (b *Button) Draw(screen *ebiten.Image) {
	// Цвет кнопки - градиент от синего к более светлому
	btnColor := color.RGBA{65, 105, 225, 255} // Royal Blue
	if b.Hovered {
		btnColor = color.RGBA{100, 149, 237, 255} // Cornflower Blue
	}

	// Основной прямоугольник кнопки
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), btnColor, false)
	
	// Тонкая рамка вокруг кнопки
	borderColor := color.RGBA{255, 255, 255, 100}
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), 1, borderColor, false)
	
	// Эффект тени при наведении
	if b.Hovered {
		shadowColor := color.RGBA{255, 255, 255, 50}
		vector.DrawFilledRect(screen, float32(b.X)+2, float32(b.Y)+2, float32(b.Width), float32(b.Height), shadowColor, false)
	}

	// Текст кнопки
	textBounds := text.BoundString(b.Font, b.Text)
	textHeight := b.Font.Metrics().Height.Ceil()
	textColor := color.White
	
	// Полужирный эффект для текста при наведении
	if b.Hovered {
		textColor = color.Black 
	}
	
	text.Draw(screen, b.Text, b.Font, 
		int(b.X)+int(b.Width)/2-textBounds.Max.X/2, 
		int(b.Y)+int(b.Height)/2+textHeight/2-2, 
		textColor)
}