package game

import (
	"image/color"
	"golang.org/x/image/font/basicfont"
)

const (
	ScreenWidth    = 640
	ScreenHeight   = 480
	GridSize       = 32
	GridWidth      = ScreenWidth / GridSize
	GridHeight     = ScreenHeight / GridSize
	PlayerSpeed    = 5
	LaneSpacing    = GridSize * 1.5 // Расстояние между полосами
	TextAreaHeight = 5              // Высота области для текста
)

// Константы для уровней сложности
const (
	Easy = iota
	Medium
	Hard
)

var (
	Font = basicfont.Face7x13
)

func GetDifficultyName(level int) string {
	switch level {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}

func GetDifficultyColor(level int) color.RGBA {
	switch level {
	case Easy:
		return color.RGBA{0, 255, 0, 255} // Green
	case Medium:
		return color.RGBA{255, 165, 0, 255} // Orange
	case Hard:
		return color.RGBA{255, 0, 0, 255} // Red
	default:
		return color.RGBA{255, 255, 255, 255} // White
	}
}