package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerCreature struct {
	xpos    float32
	ypos    float32
	speed   float32
	size    float32
	value   float32 //point value for creature
	color   rl.Color
	texture rl.Texture2D
}

// New Creature containing values from PlayerCreature struct
func NewCreature(xpos, ypos, speed, size, value float32, color rl.Color, texture rl.Texture2D) PlayerCreature {
	return PlayerCreature{
		xpos:    xpos,
		ypos:    ypos,
		speed:   speed,
		size:    size,
		value:   value,
		color:   color,
		texture: texture,
	}
}

// function to draw creature from the PlayerCreature struct
func (c PlayerCreature) DrawCreature() {
	scale := c.size / float32(c.texture.Width) //scale factor based on size
	rl.DrawTextureEx(c.texture, rl.Vector2{X: c.xpos, Y: c.ypos}, 0, scale, rl.White)
	rl.DrawText(fmt.Sprintf("%d", int(c.value)), int32(c.xpos-10), int32(c.ypos-10), 20, rl.Black) //text of the point value
}

// function to move the creature
func (c *PlayerCreature) Move(xOffset, yOffset float32) {
	c.xpos += xOffset //* c.speed * rl.GetFrameTime()
	c.ypos += yOffset //* c.speed * rl.GetFrameTime()

	//top and bottom boundaries
	if c.ypos < 0 {
		c.ypos = 0
	}
	if c.ypos+c.size > float32(rl.GetScreenHeight()) {
		c.ypos = float32(rl.GetScreenHeight()) - c.size
	}
	//left and right boundaries
	if c.xpos < 0 {
		c.xpos = 0
	}
	if c.xpos+c.size > float32(rl.GetScreenWidth()) {
		c.xpos = float32(rl.GetScreenWidth()) - c.size
	}
}

// function to check if two creatures overlap
// TODO: add in functionality to detect when the player collides with the enemy
func checkOverlap(a, b PlayerCreature) bool {
	return a.xpos < b.xpos+b.size && a.xpos+a.size > b.xpos &&
		a.ypos < b.ypos+b.size && a.ypos+a.size > b.ypos
}

func main() {
	rl.InitWindow(1980, 1080, "Weekly Project 5 Mini Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// picture for creatures
	creatureImage := rl.LoadTexture("C:/_dev/Go/cs4953/WeeklyProject5/WeeklyProject5cs4953/Textures/Mouse.png")
	defer rl.UnloadTexture(creatureImage)

	//initialize player creature outside of the game loop
	player := NewCreature(50, 50, 200, 50, 1, rl.Red, creatureImage)

	//slice to hold the enemy creatures
	enemyCreatures := make([]PlayerCreature, 0)
	enemySize := float32(50)

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

	//start drawing
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

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
		rl.EndDrawing()
	}
}
