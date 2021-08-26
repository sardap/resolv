package main

import (
	"github.com/SolarLune/resolv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldCompound struct {
	Space  *resolv.Space
	Player *resolv.Space
}

func (w *WorldCompound) Create() {

	w.Space = resolv.NewSpace()

	w.Space.Clear()
	w.Space.Add(resolv.NewRectangle(0, 0, screenWidth, cell))
	w.Space.Add(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	w.Space.Add(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	w.Space.Add(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))
	w.Space.AddTags("solid")

	// Add the "solid" tag to all Shapes within the Space
	for i := 0; i < 200; i++ {
		square := NewSquare(w.Space)
		square.Rect.AddTags("sticky")
		w.Space.Add(square.Rect)
	}

	w.Player = resolv.NewSpace()
	w.Player.Add(NewSquare(w.Space).Rect) // Gonna be lazy here and use a new Square, but just its Rect.
	w.Player.AddTags("player")
	w.Space.Add(w.Player)

}

func (w *WorldCompound) Update() {

	player := w.Player
	moveSpd := float64(1)
	dx, dy := float64(0), float64(0)

	if rl.IsKeyDown(rl.KeyRight) {
		dx = moveSpd
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		dx = -moveSpd
	}

	if rl.IsKeyDown(rl.KeyUp) {
		dy = -moveSpd
	}
	if rl.IsKeyDown(rl.KeyDown) {
		dy = moveSpd
	}

	solids := w.Space.FilterByTags("solid")

	var stickTo resolv.Shape

	if res := solids.Resolve(player, dx, 0); res.Colliding() {
		player.Move(res.ResolveX, 0)

		if res.ShapeB.HasTags("sticky") {
			stickTo = res.ShapeB
		}

	} else {
		player.Move(dx, 0)
	}

	if res := solids.Resolve(player, 0, dy); res.Colliding() {
		player.Move(0, res.ResolveY)

		if res.ShapeB.HasTags("sticky") {
			stickTo = res.ShapeB
		}
	} else {
		player.Move(0, dy)
	}

	if stickTo != nil {
		stickTo.AddTags("player")
		player.Add(stickTo)
		w.Space.Remove(stickTo)
	}

	if rl.IsKeyDown(rl.KeyX) {
		if player.Length() > 1 {
			detach := player.Get(player.Length() - 1)
			detach.RemoveTags("player")
			w.Space.Add(detach)
			player.Remove(detach)
		}
	}

}

func (w *WorldCompound) DrawObject(other resolv.Shape) {

	switch shape := other.(type) {

	case *resolv.Rectangle:

		if !shape.HasTags("square") {

			rl.DrawRectangleLines(int32(shape.X), int32(shape.Y), int32(shape.W), int32(shape.H), rl.LightGray)

		} else {

			squareData := shape.GetData().(*Square)

			color := rl.Color{0, 0, 255, 255}

			if shape.HasTags("player") {
				color = rl.Color{0, 255, 0, 255}
			}

			rl.DrawRectangleLines(int32(squareData.Rect.X), int32(squareData.Rect.Y), int32(squareData.Rect.W), int32(squareData.Rect.H), color)

		}

	case *resolv.Space:

		for _, obj := range *shape {
			w.DrawObject(obj)
		}

	}

}

func (w *WorldCompound) Draw() {

	for _, other := range *w.Space {

		w.DrawObject(other)

	}

	if drawHelpText {
		DrawText(32, 16,
			"-Compound test-",
			"You're the green square. Use",
			"the arrow keys to move.",
			"Touch a blue square to absorb it.",
			"Press the X key to detach absorbed squares.",
		)
	}
}

func (w *WorldCompound) Destroy() {
	w.Player = nil
	w.Space.Clear()
}
