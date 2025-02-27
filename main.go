package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerCreature struct {
	xpos  float32
	ypos  float32
	speed float32
	size  float32
	value float32 //point value for creature
	color rl.Color
}

// New Creature containing values from PlayerCreature struct
func NewCreature(xpos, ypos, speed, size, value float32, color rl.Color) PlayerCreature {
	return PlayerCreature{
		xpos:  xpos,
		ypos:  ypos,
		speed: speed,
		size:  size,
		value: value,
		color: color,
	}
}

// function to draw creature from the PlayerCreature struct
func (c PlayerCreature) DrawCreature() {
	rl.DrawRectangle(int32(c.xpos), int32(c.ypos), int32(c.size), int32(c.size), c.color)
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

func main() {
	rl.InitWindow(1980, 1080, "Weekly Project 5 Mini Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	//initialize player creature outside of the game loop
	player := NewCreature(100, 100, 200, 50, 1, rl.Beige)

	//slice to hold the enemy creatures
	enemyCreatures := make([]PlayerCreature, 0)
	//make 5 random creatures
	for i := 0; i < 5; i++ {
		enemyCreatures = append(enemyCreatures, NewCreature(float32(rand.IntN(rl.GetScreenWidth())), float32(rand.IntN(rl.GetScreenHeight())), float32(0), float32(50), float32(1), rl.Red))
	}

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
			fmt.Println("W Key Pressed")
		}
		if rl.IsKeyPressed(rl.KeyS) {
			player.Move(0, 50)
			fmt.Println("S Key Pressed")
		}
		if rl.IsKeyPressed(rl.KeyA) {
			player.Move(-50, 0)
			fmt.Println("A Key Pressed")
		}
		if rl.IsKeyPressed(rl.KeyD) {
			player.Move(50, 0)
			fmt.Println("D Key Pressed")
		}
		rl.EndDrawing()
	}
}
