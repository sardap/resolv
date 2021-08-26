package main

import (
	"math/rand"

	"github.com/SolarLune/resolv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldInterface interface {
	Create()
	Update()
	Draw()
	Destroy()
}

type Square struct {
	Rect        *resolv.Rectangle
	SpeedX      float64
	SpeedY      float64
	BounceFrame float64
}

func NewSquare(space *resolv.Space) *Square {

	square := &Square{Rect: resolv.NewRectangle(cell*2+float64(rand.Int31n(screenWidth-cell*4)), cell*2+float64(rand.Int31n(screenHeight-cell*4)), cell, cell),
		SpeedX: (0.5 - rand.Float64()) * 8,
		SpeedY: (0.5 - rand.Float64()) * 8}

	// Attempt to not spawn a Square in an occupied location
	for i := 0; i < 100; i++ {

		if space.IsColliding(square.Rect) {

			square.Rect.X = cell*2 + float64(rand.Int31n(screenWidth-cell*4))
			square.Rect.Y = cell*2 + float64(rand.Int31n(screenHeight-cell*4))

		}

	}

	square.Rect.AddTags("square", "solid")

	// We set a pointer to the square on the Rect itself so if another Shape has a collision with it, we can check the data pointer to see
	// what the struct is.
	square.Rect.SetData(square)

	return square

}

func DrawText(x, y float64, textLines ...string) {

	for i, line := range textLines {

		// length := rl.MeasureText(line, 8)
		// rl.DrawRectangle(x-2, y+(float64(i)*10), length+2, 8, rl.Black)
		rl.DrawText(line, int32(x), int32(y)+int32(i*10), 8, rl.Blue)
		rl.DrawText(line, int32(x), int32(y)-1+int32(i*10), 8, rl.White)

	}

}
