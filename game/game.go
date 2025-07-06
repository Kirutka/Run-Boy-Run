package game

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player         *GameObject
	background     *ebiten.Image
	objects        map[string]*ebiten.Image
	cars           []*GameObject
	currentTime    int
	lastUpdateTime time.Time
	gameState      string // "menu", "playing", "paused", "win", "lose"
	elapsedTime    float64
	buttons        map[string]*Button
	difficulty     int
	levelTime      int // Время для текущего уровня
}

func NewGame() *Game {
	g := &Game{
		objects:     make(map[string]*ebiten.Image),
		gameState:   "menu",
		player:      &GameObject{},
		buttons:     make(map[string]*Button),
		difficulty:  Easy, // Начинаем с легкого уровня
	}
	g.LoadImages()
	g.createButtons()
	g.setDifficulty(Easy) // Устанавливаем начальную сложность
	return g
}

// Установка параметров сложности
func (g *Game) setDifficulty(level int) {
	g.difficulty = level
	
	switch level {
	case Easy:
		g.levelTime = 30 
	case Medium:
		g.levelTime = 25 
	case Hard:
		g.levelTime = 10 
	}
	
	g.currentTime = g.levelTime
	g.initializeGame()
}

// Получение параметров для текущего уровня сложности
func (g *Game) getLevelParams() (int, int, float64, float64, int, int) {
	switch g.difficulty {
	case Easy:
		return 7, 2, 2, 2.5, 7, 11 
	case Medium:
		return 8, 3, 2.5, 3.5, 6, 9  
	case Hard:
		return 9, 4, 3, 4.5, 5, 7  
	default:
		return 7, 2, 2, 2.5, 7, 11
	}
}

func (g *Game) createButtons() {
	// Кнопки выбора уровня сложности
	g.buttons["easy"] = &Button{
		X:      ScreenWidth/2 - 150,
		Y:      ScreenHeight/2 - 20,
		Width:   140,
		Height:  40,
		Text:    "Easy",
		Font:    Font,
		Action: func() {
			g.setDifficulty(Easy)
			g.gameState = "playing"
			g.lastUpdateTime = time.Now()
		},
	}

	g.buttons["medium"] = &Button{
		X:      ScreenWidth/2 - 70,
		Y:      ScreenHeight/2 - 20,
		Width:   140,
		Height:  40,
		Text:    "Medium",
		Font:    Font,
		Action: func() {
			g.setDifficulty(Medium)
			g.gameState = "playing"
			g.lastUpdateTime = time.Now()
		},
	}

	g.buttons["hard"] = &Button{
		X:      ScreenWidth/2 + 10,
		Y:      ScreenHeight/2 - 20,
		Width:   140,
		Height:  40,
		Text:    "Hard",
		Font:    Font,
		Action: func() {
			g.setDifficulty(Hard)
			g.gameState = "playing"
			g.lastUpdateTime = time.Now()
		},
	}

	// Кнопка "Выйти из игры" в главном меню
	g.buttons["exit_menu"] = &Button{
		X:      ScreenWidth/2 - 100,
		Y:      ScreenHeight/2 + 40,
		Width:   200,
		Height:  40,
		Text:    "Exit Game",
		Font:    Font,
		Action: func() {
			os.Exit(0)
		},
	}

	// Кнопка "Выйти в меню" в паузе
	g.buttons["exit_pause"] = &Button{
		X:      ScreenWidth/2 - 100,
		Y:      ScreenHeight/2 + 50,
		Width:   200,
		Height:  40,
		Text:    "Back to Menu",
		Font:    Font,
		Action: func() {
			g.gameState = "menu"
		},
	}

	// Кнопка "Рестарт" при проигрыше/выигрыше
	g.buttons["restart"] = &Button{
		X:      ScreenWidth/2 - 100,
		Y:      ScreenHeight/2 + 50,
		Width:   200,
		Height:  40,
		Text:    "Play Again",
		Font:    Font,
		Action: func() {
			g.setDifficulty(g.difficulty)
			g.gameState = "playing"
			g.lastUpdateTime = time.Now()
		},
	}

	// Кнопка "Меню" при проигрыше/выигрыше
	g.buttons["menu"] = &Button{
		X:      ScreenWidth/2 - 100,
		Y:      ScreenHeight/2 + 100,
		Width:   200,
		Height:  40,
		Text:    "Main Menu",
		Font:    Font,
		Action: func() {
			g.gameState = "menu"
		},
	}
}

func (g *Game) LoadImages() {
	// Загрузка реальных изображений вместо цветных placeholder'ов
	var err error
	
	// Загрузка изображения машины - исправлен путь
	g.objects["car"], err = loadImageFromFile("../image/bus.png", 64, 32)
	if err != nil {
		log.Printf("Failed to load car image: %v, using placeholder", err)
		g.objects["car"] = g.createPlaceholderImage(64, 32, color.RGBA{255, 0, 0, 255})
	}
	
	// Загрузка изображения игрока
	g.objects["player"], err = loadImageFromFile("../image/player.png", 32, 32)
	if err != nil {
		log.Printf("Failed to load player image: %v, using placeholder", err)
		g.objects["player"] = g.createPlaceholderImage(32, 32, color.RGBA{0, 255, 0, 255})
	}
	
	// Загрузка фонового изображения
	g.objects["background"], err = loadImageFromFile("../image/back.png", ScreenWidth, ScreenHeight)
	if err != nil {
		log.Printf("Failed to load background image: %v, using placeholder", err)
		g.objects["background"] = g.createPlaceholderImage(ScreenWidth, ScreenHeight, color.RGBA{200, 200, 200, 255})
	}

	g.background = g.objects["background"]
	g.player.Image = g.objects["player"]
}

func loadImageFromFile(path string, targetWidth, targetHeight int) (*ebiten.Image, error) {
	// Проверяем существование файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Декодируем изображение
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	
	// Конвертируем в ebiten image
	ebitenImg := ebiten.NewImageFromImage(img)
	
	// Масштабируем если нужно
	if targetWidth > 0 && targetHeight > 0 {
		origWidth, origHeight := ebitenImg.Size()
		if origWidth != targetWidth || origHeight != targetHeight {
			scaledImg := ebiten.NewImage(targetWidth, targetHeight)
			op := &ebiten.DrawImageOptions{}
			
			scaleX := float64(targetWidth) / float64(origWidth)
			scaleY := float64(targetHeight) / float64(origHeight)
			op.GeoM.Scale(scaleX, scaleY)
			
			scaledImg.DrawImage(ebitenImg, op)
			return scaledImg, nil
		}
	}
	
	return ebitenImg, nil
}

func (g *Game) createPlaceholderImage(width, height int, col color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(col)
	return img
}

func (g *Game) initializeGame() {
	// Получаем параметры для текущего уровня сложности
	numLanes, numCarsPerLane, carSpeedMin, carSpeedMax, minCarGap, maxCarGap := g.getLevelParams()

	// Настройка начального положения игрока
	g.player = &GameObject{
		X:      float64(GridWidth/2) * GridSize,
		Y:      float64((GridHeight - 1) * GridSize),
		Speed:  PlayerSpeed,
		Image:  g.objects["player"],
		Width:  GridSize,
		Height: GridSize,
	}

	// Очистка существующих автомобилей
	g.cars = []*GameObject{}

	// Инициализация автомобилей на каждой полосе
	for lane := 0; lane < numLanes; lane++ {
		lastCarX := -float64(GridSize)
		for i := 0; i < numCarsPerLane; i++ {
			minGap := lastCarX + float64(minCarGap*GridSize)
			maxGap := lastCarX + float64(maxCarGap*GridSize)
			carX := minGap + rand.Float64()*(maxGap-minGap)

			g.cars = append(g.cars, &GameObject{
				X:       carX,
				Y:       float64(lane)*LaneSpacing + TextAreaHeight,
				Speed:   carSpeedMin + rand.Float64()*(carSpeedMax-carSpeedMin),
				Image:   g.objects["car"],
				Width:   GridSize * 2,
				Height:  GridSize,
				IsRight: rand.Intn(2) == 0,
			})

			lastCarX = carX
		}
	}
}

func (g *Game) Update() error {
	// Обработка паузы
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && (g.gameState == "playing" || g.gameState == "paused") {
		if g.gameState == "playing" {
			g.gameState = "paused"
		} else {
			g.gameState = "playing"
			g.lastUpdateTime = time.Now()
		}
	}

	switch g.gameState {
	case "menu":
		g.updateMenu()
	case "playing":
		g.updateGame()
	case "paused":
		g.updatePaused()
	case "win", "lose":
		g.updateGameOver()
	}

	return nil
}

func (g *Game) updateMenu() {
	// Обновление кнопок в меню
	mx, my := ebiten.CursorPosition()
	g.buttons["easy"].Hovered = g.buttons["easy"].Contains(float64(mx), float64(my))
	g.buttons["medium"].Hovered = g.buttons["medium"].Contains(float64(mx), float64(my))
	g.buttons["hard"].Hovered = g.buttons["hard"].Contains(float64(mx), float64(my))
	g.buttons["exit_menu"].Hovered = g.buttons["exit_menu"].Contains(float64(mx), float64(my))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.buttons["easy"].Hovered {
			g.buttons["easy"].Action()
		}
		if g.buttons["medium"].Hovered {
			g.buttons["medium"].Action()
		}
		if g.buttons["hard"].Hovered {
			g.buttons["hard"].Action()
		}
		if g.buttons["exit_menu"].Hovered {
			g.buttons["exit_menu"].Action()
		}
	}
}

func (g *Game) updateGame() {
	now := time.Now()
	elapsed := now.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = now

	// Управление игроком
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= float64(GridSize) * elapsed * PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += float64(GridSize) * elapsed * PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= float64(GridSize) * elapsed * PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += float64(GridSize) * elapsed * PlayerSpeed
	}

	// Ограничение движения игрока
	g.player.X = clamp(g.player.X, 0, ScreenWidth-float64(GridSize))
	g.player.Y = clamp(g.player.Y, TextAreaHeight, ScreenHeight-float64(GridSize))

	// Проверка победы - достиг верха экрана
	if g.player.Y <= TextAreaHeight {
		g.gameState = "win"
		return
	}

	// Обновление автомобилей
	for _, car := range g.cars {
		car.Update(elapsed, float64(ScreenWidth))
	}

	// Проверка столкновений
	g.checkCollisions()

	// Обновление времени
	g.elapsedTime += elapsed
	if g.elapsedTime >= 1.0 {
		g.currentTime -= 1
		g.elapsedTime = 0
		
		// Проверка на проигрыш по времени
		if g.currentTime <= 0 {
			g.gameState = "lose" // Время вышло - проигрыш
			return
		}
	}
}

func (g *Game) updatePaused() {
	// Обновление кнопок в паузе
	mx, my := ebiten.CursorPosition()
	g.buttons["exit_pause"].Hovered = g.buttons["exit_pause"].Contains(float64(mx), float64(my))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.buttons["exit_pause"].Hovered {
			g.buttons["exit_pause"].Action()
		}
	}
}

func (g *Game) updateGameOver() {
	// Обновление кнопок
	mx, my := ebiten.CursorPosition()
	g.buttons["restart"].Hovered = g.buttons["restart"].Contains(float64(mx), float64(my))
	g.buttons["menu"].Hovered = g.buttons["menu"].Contains(float64(mx), float64(my))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.buttons["restart"].Hovered {
			g.buttons["restart"].Action()
		}
		if g.buttons["menu"].Hovered {
			g.buttons["menu"].Action()
		}
	}
}

func (g *Game) checkCollisions() {
	playerRect := g.player.GetRect()

	for _, car := range g.cars {
		carRect := car.GetRect()
		if playerRect.Overlaps(carRect) {
			g.gameState = "lose"
			return
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Отрисовка фона
	if g.background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.background, op)
	}

	switch g.gameState {
	case "menu":
		g.drawMenu(screen)
	case "playing":
		g.drawGame(screen)
	case "paused":
		g.drawGame(screen)
		g.drawPauseMenu(screen)
	case "win", "lose":
		g.drawGame(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	// Фон меню с легкой прозрачностью
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 180}, false)

	// Заголовок игры
	title := "ROAD ADVENTURE"
	titleBounds := text.BoundString(Font, title)
	text.Draw(screen, title, Font, ScreenWidth/2-titleBounds.Max.X/2, 80, color.RGBA{255, 215, 0, 255}) // Золотой цвет

	// Подзаголовок выбора сложности
	diffText := "SELECT DIFFICULTY"
	diffBounds := text.BoundString(Font, diffText)
	text.Draw(screen, diffText, Font, ScreenWidth/2-diffBounds.Max.X/2, 140, color.White)

	// Расположение кнопок сложности в одну линию с равными промежутками
	buttonY := 180
	g.buttons["easy"].X = ScreenWidth/2 - 220
	g.buttons["easy"].Y = float64(buttonY)
	g.buttons["medium"].X = ScreenWidth/2 - 70
	g.buttons["medium"].Y = float64(buttonY)
	g.buttons["hard"].X = ScreenWidth/2 + 80
	g.buttons["hard"].Y = float64(buttonY)

	// Отрисовка кнопок сложности
	g.buttons["easy"].Draw(screen)
	g.buttons["medium"].Draw(screen)
	g.buttons["hard"].Draw(screen)

	// Разделительная линия
	separatorY := buttonY + 70
	vector.StrokeLine(screen, ScreenWidth/4, float32(separatorY), ScreenWidth*3/4, float32(separatorY), 2, color.RGBA{100, 100, 100, 255}, false)

	// Управление
	controlsTitle := "CONTROLS"
	controlsTitleBounds := text.BoundString(Font, controlsTitle)
	text.Draw(screen, controlsTitle, Font, ScreenWidth/2-controlsTitleBounds.Max.X/2, separatorY+30, color.White)

	controls := []string{
		"W/A/S/D or Arrow Keys - Movement",
		"Space - Pause/Resume",
		"ESC - Back to Menu",
	}
	for i, line := range controls {
		lineBounds := text.BoundString(Font, line)
		text.Draw(screen, line, Font, ScreenWidth/2-lineBounds.Max.X/2, separatorY+60+i*20, color.RGBA{200, 200, 200, 255})
	}

	// Кнопка выхода - центрированная внизу
	g.buttons["exit_menu"].X = ScreenWidth/2 - 100
	g.buttons["exit_menu"].Y = float64(ScreenHeight - 80)
	g.buttons["exit_menu"].Draw(screen)

	// Версия игры или авторские права
	versionText := "v1.0 © 2024"
	versionBounds := text.BoundString(Font, versionText)
	text.Draw(screen, versionText, Font, ScreenWidth-versionBounds.Max.X-10, ScreenHeight-20, color.RGBA{150, 150, 150, 255})
}

func (g *Game) drawGame(screen *ebiten.Image) {
	// Отрисовка автомобилей
	for _, car := range g.cars {
		car.Draw(screen)
	}

	// Отрисовка игрока
	g.player.Draw(screen)

	// Отрисовка времени и уровня сложности
	levelText := GetDifficultyName(g.difficulty)
	
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time: %d", g.currentTime), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Level: %s", levelText), 10, 30)
}

func (g *Game) drawPauseMenu(screen *ebiten.Image) {
	// Полупрозрачный фон
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 150}, false)

	// Текст паузы
	pauseText := "PAUSE"
	pauseBounds := text.BoundString(Font, pauseText)
	text.Draw(screen, pauseText, Font, ScreenWidth/2-pauseBounds.Max.X/2, ScreenHeight/2-50, color.White)

	// Кнопка "Выйти в меню"
	g.buttons["exit_pause"].Draw(screen)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Полупрозрачный фон
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 150}, false)

	// Текст результата
	resultText := ""
	if g.gameState == "win" {
		resultText = "VICTORY!"
	} else {
		resultText = "GAME OVER!"
	}

	reasonText := ""
	if g.gameState == "lose" && g.currentTime <= 0 {
		reasonText = "Time's up!"
	} else if g.gameState == "lose" {
		reasonText = "You got hit!"
	} else {
		reasonText = "You made it!"
	}

	resultBounds := text.BoundString(Font, resultText)
	text.Draw(screen, resultText, Font, ScreenWidth/2-resultBounds.Max.X/2, ScreenHeight/2-80, color.White)
	
	reasonBounds := text.BoundString(Font, reasonText)
	text.Draw(screen, reasonText, Font, ScreenWidth/2-reasonBounds.Max.X/2, ScreenHeight/2-50, color.White)

	// Кнопки
	g.buttons["restart"].Draw(screen)
	g.buttons["menu"].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}