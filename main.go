package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerCreature struct {
	Xpos    float32
	Ypos    float32
	Speed   float32
	size    float32
	Value   float32 //point value for creature
	color   rl.Color
	texture rl.Texture2D
}

// New Creature containing values from PlayerCreature struct
func NewCreature(xpos, ypos, speed, size, value float32, color rl.Color, texture rl.Texture2D) PlayerCreature {
	return PlayerCreature{
		Xpos:    xpos,
		Ypos:    ypos,
		Speed:   speed,
		size:    size,
		Value:   value,
		color:   color,
		texture: texture,
	}
}

// function to draw creature from the PlayerCreature struct
func (c PlayerCreature) DrawCreature() {
	scale := c.size / float32(c.texture.Width) //scale factor based on size
	rl.DrawTextureEx(c.texture, rl.Vector2{X: c.Xpos, Y: c.Ypos}, 0, scale, c.color)
	rl.DrawText(fmt.Sprintf("%d", int(c.Value)), int32(c.Xpos-10), int32(c.Ypos-10), 20, rl.Black) //text of the point value
}

// function to move the creature
func (c *PlayerCreature) Move(xOffset, yOffset float32) {
	c.Xpos += xOffset //* c.speed * rl.GetFrameTime()
	c.Ypos += yOffset //* c.speed * rl.GetFrameTime()

	//top and bottom boundaries
	if c.Ypos < 0 {
		c.Ypos = 0
	}
	if c.Ypos+c.size > float32(rl.GetScreenHeight()) {
		c.Ypos = float32(rl.GetScreenHeight()) - c.size
	}
	//left and right boundaries
	if c.Xpos < 0 {
		c.Xpos = 0
	}
	if c.Xpos+c.size > float32(rl.GetScreenWidth()) {
		c.Xpos = float32(rl.GetScreenWidth()) - c.size
	}
}

// function to check if two creatures overlap
func checkOverlap(a, b PlayerCreature) bool {
	return a.Xpos < b.Xpos+b.size && a.Xpos+a.size > b.Xpos &&
		a.Ypos < b.Ypos+b.size && a.Ypos+a.size > b.Ypos
}

// function to restart the game state
func resetGame(creatureImage rl.Texture2D, enemySize float32) (PlayerCreature, []PlayerCreature, bool) {
	player := NewCreature(50, 50, 200, 50, 1, rl.Blue, creatureImage)
	enemyCreatures := make([]PlayerCreature, 0)

	for i := 0; i < 5; i++ {
		var newCreature PlayerCreature
		var overlapping bool

		//keep generating until valid position where no overlapping is found
		for {

			//ensuring creatures don't spawn out of bounds
			x := float32(rand.IntN(rl.GetScreenWidth() - int(enemySize)))
			y := float32(rand.IntN(rl.GetScreenHeight() - int(enemySize)))

			newCreature = NewCreature(
				x, y, 0, enemySize, float32(i+1), rl.Red, creatureImage,
			)

			//assuming no overlapping
			overlapping = false

			//check overlap with player
			if checkOverlap(newCreature, player) {
				overlapping = true
			}

			//check previously added creatures
			for _, existing := range enemyCreatures {
				if checkOverlap(newCreature, existing) {
					overlapping = true
					break //stop if overlapping has occurred
				}
			}

			//if no overlapping, break out of loop
			if !overlapping {
				break
			}
		}

		//add new non-overlapping creature to the slice
		enemyCreatures = append(enemyCreatures, newCreature)
	}
	return player, enemyCreatures, false //reset player, enemy creatures, and gameOver state
}

func (c *PlayerCreature) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", " ")
	fmt.Println("Saving...")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return os.WriteFile(filename+".json", data, 0644)
}

func (c *PlayerCreature) Load(filename string) error {
	data, err := os.ReadFile(filename + ".json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return json.Unmarshal(data, c)
}

func main() {
	rl.InitWindow(1980, 1080, "Weekly Project 5 Mini Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// picture for creatures
	creatureImage := rl.LoadTexture("C:/_dev/Go/cs4953/WeeklyProject5/WeeklyProject5cs4953/Textures/Mouse.png")
	defer rl.UnloadTexture(creatureImage)

	//initialize player creature outside of the game loop
	player := NewCreature(50, 50, 200, 50, 1, rl.Blue, creatureImage)

	//slice to hold the enemy creatures
	enemyCreatures := make([]PlayerCreature, 0)
	enemySize := float32(50) //enemy size used in loop

	for i := 0; i < 5; i++ {
		var newCreature PlayerCreature
		var overlapping bool

		//keep generating until valid position where no overlapping is found
		for {

			//ensuring creatures don't spawn out of bounds
			x := float32(rand.IntN(rl.GetScreenWidth() - int(enemySize)))
			y := float32(rand.IntN(rl.GetScreenHeight() - int(enemySize)))

			newCreature = NewCreature(
				x, y, 0, enemySize, float32(i+1), rl.Red, creatureImage,
			)

			//assuming no overlapping
			overlapping = false

			//check overlap with player
			if checkOverlap(newCreature, player) {
				overlapping = true
			}

			//check previously added creatures
			for _, existing := range enemyCreatures {
				if checkOverlap(newCreature, existing) {
					overlapping = true
					break //stop if overlapping has occurred
				}
			}

			//if no overlapping, break out of loop
			if !overlapping {
				break
			}
		}

		//add new non-overlapping creature to the slice
		enemyCreatures = append(enemyCreatures, newCreature)
	}

	var gameOver bool = false //bool used for gameover state
	var gameWon bool = false  //bool used to track if game is won

	//start drawing
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		//game over logic
		if gameOver {
			text := "You lost! Click R to restart."
			textWidth := rl.MeasureText(text, 40)
			rl.DrawText(text, (int32(rl.GetScreenWidth())-(textWidth))/2, int32(rl.GetScreenHeight()/2), 40, rl.Red)

			if rl.IsKeyPressed(rl.KeyR) {
				player, enemyCreatures, gameOver = resetGame(creatureImage, enemySize) //reset game
				gameWon = false
			}
		} else if gameWon {
			//check if player has defeated all creatures
			text := "You Win! Click R to restart."
			textWidth := rl.MeasureText(text, 40)
			rl.DrawText(text, (int32(rl.GetScreenWidth())-(textWidth))/2, int32(rl.GetScreenHeight()/2), 40, rl.Green)

			if rl.IsKeyPressed(rl.KeyR) {
				player, enemyCreatures, gameOver = resetGame(creatureImage, enemySize) //reset game
				gameWon = false
			}
		} else {
			//draw player
			player.DrawCreature()

			//drawing the enemies
			for _, enemy := range enemyCreatures {
				enemy.DrawCreature()
			}

			//player movement
			if rl.IsKeyPressed(rl.KeyW) {
				player.Move(0, -50)
			}
			if rl.IsKeyPressed(rl.KeyS) {
				player.Move(0, 50)
			}
			if rl.IsKeyPressed(rl.KeyA) {
				player.Move(-50, 0)
			}
			if rl.IsKeyPressed(rl.KeyD) {
				player.Move(50, 0)
			}

			//save function?
			if rl.IsKeyPressed(rl.KeyP) {
				player.Save("SaveFile")

			}
			if rl.IsKeyPressed(rl.KeyL) {
				player.Load("SaveFile")
			}

			//check collisions with enemy creatures
			for i := 0; i < len(enemyCreatures); i++ {
				if checkOverlap(player, enemyCreatures[i]) {
					if player.Value >= enemyCreatures[i].Value {
						//player absorbs enemy
						player.Value += enemyCreatures[i].Value

						//remove said enemy from slice
						enemyCreatures = append(enemyCreatures[:i], enemyCreatures[i+1:]...)
						i-- //adjust index after removing element
					} else {
						//player loses
						gameOver = true
					}
				}
				if len(enemyCreatures) == 0 {
					gameWon = true
				}
			}
		}
		rl.EndDrawing()
	}
}
