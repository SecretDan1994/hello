package main

import (
	// "fmt"
	// "github.com/gen2brain/raylib-go/raygui"
	"math"
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	PLAYER_MAX_LIFE = 5
	LINES_OF_BRICKS = 5
	BRICKS_PER_LINE = 20
)

type Player struct {
	position rl.Vector2
	size     rl.Vector2
	life     int
}

type Ball struct {
	position rl.Vector2
	speed    rl.Vector2
	radius   float32
	active   bool
}

type Brick struct {
	position rl.Vector2
	active   bool
}


type Game struct {
	gameOver  bool
	pause     bool
	player    Player
	ball      Ball
	brick     [LINES_OF_BRICKS][BRICKS_PER_LINE]Brick
	brickSize rl.Vector2
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "test game")

	// Init game
	game := NewGame()
	game.gameOver = true

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()

		game.Draw()

	}

	rl.CloseWindow()

}

// NewGame - Create a new game instance
func NewGame() (g Game) {
	g.Init()
	return
}

// Init - initialize game
func (g *Game) Init() {
	g.brickSize = rl.Vector2{float32(rl.GetScreenWidth() / BRICKS_PER_LINE), 40}

	// Initialize player
	g.player.position = rl.Vector2{float32(screenWidth / 2), float32(screenHeight * 7 / 8)}
	g.player.size = rl.Vector2{float32(screenWidth / 10), 20}
	g.player.life = PLAYER_MAX_LIFE

	// Initialize ball
	g.ball.position = rl.Vector2{float32(screenWidth / 2), float32(screenHeight*7/8 - 30)}
	g.ball.speed = rl.Vector2{0, 0}
	g.ball.radius = 7
	g.ball.active = false

	initialDownPosition := int(50)

	for i := 0; i < LINES_OF_BRICKS; i++ {
		for j := 0; j < BRICKS_PER_LINE; j++ {
			g.brick[i][j].position = rl.Vector2{float32(j)*g.brickSize.X + g.brickSize.X/2, float32(i)*g.brickSize.Y + float32(initialDownPosition)}
			g.brick[i][j].active = true
		}
	}
}

// Update - update game state
func (g *Game) Update() {
	if !g.gameOver {

		if rl.IsKeyPressed(rl.KeyP) {
			g.pause = !g.pause
		}

		if !g.pause {

			if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
				g.player.position.X -= 5
			}
			if (g.player.position.X - g.player.size.X/2) <= 0 {
				g.player.position.X = g.player.size.X / 2
			}
			if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
				g.player.position.X += 5
			}
			if (g.player.position.X + g.player.size.X/2) >= screenWidth {
				g.player.position.X = screenWidth - g.player.size.X/2
			}

			if !g.ball.active {
				if rl.IsKeyPressed(rl.KeySpace) {
					g.ball.active = true
					g.ball.speed = rl.Vector2{0, -5}
				}
			}

			if g.ball.active {
				g.ball.position.X += g.ball.speed.X
				g.ball.position.Y += g.ball.speed.Y
			} else {
				g.ball.position = rl.Vector2{g.player.position.X, screenHeight*7/8 - 30}
			}

			// Collision logic: ball vs walls
			if ((g.ball.position.X + g.ball.radius) >= screenWidth) || ((g.ball.position.X - g.ball.radius) <= 0) {
				g.ball.speed.X *= -1
			}
			if (g.ball.position.Y - g.ball.radius) <= 0 {
				g.ball.speed.Y *= -1
			}
			if (g.ball.position.Y + g.ball.radius) >= screenHeight {
				g.ball.speed = rl.Vector2{0, 0}
				g.ball.active = false

				g.player.life--
			}
			if (rl.CheckCollisionCircleRec(g.ball.position, g.ball.radius,
				rl.Rectangle{g.player.position.X - g.player.size.X/2, g.player.position.Y - g.player.size.Y/2, g.player.size.X, g.player.size.Y})) {
				if g.ball.speed.Y > 0 {
					g.ball.speed.Y *= -1
					g.ball.speed.X = (g.ball.position.X - g.player.position.X) / (g.player.size.X / 2) * 5
				}
			}
			// Collision logic: ball vs bricks
			for i := 0; i < LINES_OF_BRICKS; i++ {
				for j := 0; j < BRICKS_PER_LINE; j++ {
					if g.brick[i][j].active {
						if ((g.ball.position.Y - g.ball.radius) <= (g.brick[i][j].position.Y + g.brickSize.Y/2)) &&
							((g.ball.position.Y - g.ball.radius) > (g.brick[i][j].position.Y + g.brickSize.Y/2 + g.ball.speed.Y)) &&
							((float32(math.Abs(float64(g.ball.position.X - g.brick[i][j].position.X)))) < (g.brickSize.X/2 + g.ball.radius*2/3)) &&
							(g.ball.speed.Y < 0) {
							// Hit below
							g.brick[i][j].active = false
							g.ball.speed.Y *= -1
						} else if ((g.ball.position.Y + g.ball.radius) >= (g.brick[i][j].position.Y - g.brickSize.Y/2)) &&
							((g.ball.position.Y + g.ball.radius) < (g.brick[i][j].position.Y - g.brickSize.Y/2 + g.ball.speed.Y)) &&
							((float32(math.Abs(float64(g.ball.position.X - g.brick[i][j].position.X)))) < (g.brickSize.X/2 + g.ball.radius*2/3)) &&
							(g.ball.speed.Y > 0) {
							// Hit above
							g.brick[i][j].active = false
							g.ball.speed.Y *= -1
						} else if ((g.ball.position.X + g.ball.radius) >= (g.brick[i][j].position.X - g.brickSize.X/2)) &&
							((g.ball.position.X + g.ball.radius) < (g.brick[i][j].position.X - g.brickSize.X/2 + g.ball.speed.X)) &&
							((float32(math.Abs(float64(g.ball.position.Y - g.brick[i][j].position.Y)))) < (g.brickSize.Y/2 + g.ball.radius*2/3)) &&
							(g.ball.speed.X > 0) {
							// Hit left
							g.brick[i][j].active = false
							g.ball.speed.X *= -1
						} else if ((g.ball.position.X - g.ball.radius) <= (g.brick[i][j].position.X + g.brickSize.X/2)) &&
							((g.ball.position.X - g.ball.radius) > (g.brick[i][j].position.X + g.brickSize.X/2 + g.ball.speed.X)) &&
							((float32(math.Abs(float64(g.ball.position.Y - g.brick[i][j].position.Y)))) < (g.brickSize.Y/2 + g.ball.radius*2/3)) &&
							(g.ball.speed.X < 0) {
							// Hit right
							g.brick[i][j].active = false
							g.ball.speed.X *= -1
						}
					}
				}
			}
		}

		// Game over logic
		if g.player.life <= 0 {
			g.gameOver = true
		} else {
			g.gameOver = true

			for i := 0; i < LINES_OF_BRICKS; i++ {
				for j := 0; j < BRICKS_PER_LINE; j++ {
					if g.brick[i][j].active {
						g.gameOver = false
					}
				}
			}
		}

	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.Init()
			g.gameOver = false
		}
	}
}

// Draw - draw game
func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	if !g.gameOver {
		// Draw player bar
		rl.DrawRectangle(int32(g.player.position.X-g.player.size.X/2), int32(g.player.position.Y-g.player.size.Y/2), int32(g.player.size.X), int32(g.player.size.Y), rl.Black)

		// Draw player lives
		for i := 0; i < g.player.life; i++ {
			rl.DrawRectangle(int32(20+40*i), screenHeight-30, 35, 10, rl.LightGray)
		}

		// Draw Ball
		rl.DrawCircleV(g.ball.position, g.ball.radius, rl.Maroon)

		for i := 0; i < LINES_OF_BRICKS; i++ {
			for j := 0; j < BRICKS_PER_LINE; j++ {
				if g.brick[i][j].active {
					if (i+j)%2 == 0 {
						rl.DrawRectangle(int32(g.brick[i][j].position.X-g.brickSize.X/2), int32(g.brick[i][j].position.Y-g.brickSize.Y/2), int32(g.brickSize.X), int32(g.brickSize.Y), rl.Gray)
					} else {
						rl.DrawRectangle(int32(g.brick[i][j].position.X-g.brickSize.X/2), int32(g.brick[i][j].position.Y-g.brickSize.Y/2), int32(g.brickSize.X), int32(g.brickSize.Y), rl.DarkGray)
					}
				}
			}
		}

		if g.pause {
			rl.DrawText("GAME PAUSED", screenWidth/2-rl.MeasureText("GAME PAUSED", 40)/2, screenHeight/2+screenHeight/4-40, 40, rl.Gray)
		}

	} else {
		str := "PRESS [ENTER] TO PLAY AGAIN"
		x := int(rl.GetScreenWidth()/2) - int(rl.MeasureText(str, 20)/2)
		y := rl.GetScreenHeight()/2 - 50
		rl.DrawText(str, int32(x), int32(y), 20, rl.Gray)
	}

	rl.EndDrawing()
}




// func main() {
// 	screenWidth := int32(800)
// 	screenHeight := int32(450)

// 	rl.SetConfigFlags(rl.FlagVsyncHint)

// 	rl.InitWindow(screenWidth, screenHeight, "raylib [gui] example - basic controls")

// 	buttonToggle := true
// 	buttonClicked := false
// 	checkboxChecked := false

// 	spinnerValue := 5
// 	sliderValue := float32(10)
// 	sliderBarValue := float32(70)
// 	progressValue := float32(0.5)

// 	comboActive := 0
// 	comboLastActive := 0
// 	toggleActive := 0

// 	toggleText := []string{"Item 1", "Item 2", "Item 3"}
// 	comboText := []string{"default_light", "default_dark", "hello_kitty", "monokai", "obsidian", "solarized_light", "solarized", "zahnrad"}

// 	var inputText string

// 	rl.SetTargetFPS(60)

// 	for !rl.WindowShouldClose() {
// 		if buttonClicked {
// 			progressValue += 0.1
// 			if progressValue >= 1.1 {
// 				progressValue = 0.0
// 			}
// 		}

// 		rl.BeginDrawing()

// 		rl.ClearBackground(rl.Beige)

// 		raygui.Label(rl.NewRectangle(50, 50, 80, 20), "Label")

// 		buttonClicked = raygui.Button(rl.NewRectangle(50, 70, 80, 40), "Button")

// 		raygui.Label(rl.NewRectangle(70, 140, 20, 20), "Checkbox")
// 		checkboxChecked = raygui.CheckBox(rl.NewRectangle(50, 140, 20, 20), checkboxChecked)

// 		raygui.Label(rl.NewRectangle(50, 190, 200, 20), "ProgressBar")
// 		raygui.ProgressBar(rl.NewRectangle(50, 210, 200, 20), progressValue)
// 		raygui.Label(rl.NewRectangle(200+50+5, 210, 20, 20), fmt.Sprintf("%.1f", progressValue))

// 		raygui.Label(rl.NewRectangle(50, 260, 200, 20), "Slider")
// 		sliderValue = raygui.Slider(rl.NewRectangle(50, 280, 200, 20), sliderValue, 0, 100)
// 		raygui.Label(rl.NewRectangle(200+50+5, 280, 20, 20), fmt.Sprintf("%.0f", sliderValue))

// 		buttonToggle = raygui.ToggleButton(rl.NewRectangle(50, 350, 100, 40), "ToggleButton", buttonToggle)

// 		raygui.Label(rl.NewRectangle(500, 50, 200, 20), "ToggleGroup")
// 		toggleActive = raygui.ToggleGroup(rl.NewRectangle(500, 70, 60, 30), toggleText, toggleActive)

// 		raygui.Label(rl.NewRectangle(500, 120, 200, 20), "SliderBar")
// 		sliderBarValue = raygui.SliderBar(rl.NewRectangle(500, 140, 200, 20), sliderBarValue, 0, 100)
// 		raygui.Label(rl.NewRectangle(500+200+5, 140, 20, 20), fmt.Sprintf("%.0f", sliderBarValue))

// 		raygui.Label(rl.NewRectangle(500, 190, 200, 20), "Spinner")
// 		spinnerValue = raygui.Spinner(rl.NewRectangle(500, 210, 200, 20), spinnerValue, 0, 100)

// 		raygui.Label(rl.NewRectangle(500, 260, 200, 20), "ComboBox")
// 		comboActive = raygui.ComboBox(rl.NewRectangle(500, 280, 200, 20), comboText, comboActive)

// 		if comboLastActive != comboActive {
// 			raygui.LoadGuiStyle(fmt.Sprintf("styles/%s.style", comboText[comboActive]))
// 			comboLastActive = comboActive
// 		}

// 		raygui.Label(rl.NewRectangle(500, 330, 200, 20), "TextBox")
// 		inputText = raygui.TextBox(rl.NewRectangle(500, 350, 200, 20), inputText)

// 		rl.EndDrawing()
// 	}

// 	rl.CloseWindow()
// }
