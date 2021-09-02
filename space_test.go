package resolv_test

import (
	"testing"

	"github.com/SolarLune/resolv"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		shapes   []resolv.Shape
		resolves []struct {
			x, y      float64
			colliding bool
			note      string
		}
	}{
		{
			name: "all rectangles",
			shapes: []resolv.Shape{
				resolv.NewRectangle(0, 0, 10, 10),
				resolv.NewRectangle(10, 0, 10, 10),
				resolv.NewRectangle(0, 100, 10, 10),
			},
			// Note this always moves the first shape
			resolves: []struct {
				x, y      float64
				colliding bool
				note      string
			}{
				{0, 0, false, "starting state no collisions"},
				{1, 0, true, "right rectangle basic"},
				{1, 1, true, "diagonally moving"},
				{11, 0, true, "try to jump past right rectangle"},
				{0, 100, true, "futher move down"},
				{0, 89, false, "far move down"},
				{0, 10, false, "move down"},
				{10000, 0, false, "extreme move"},
			},
		},
		{
			name: "rectnagles and lines",
			shapes: []resolv.Shape{
				resolv.NewRectangle(0, 0, 3, 3),
				resolv.NewLine(10, 0, 10, 5),
				resolv.NewLine(0, 10, 5, 10),
			},
			resolves: []struct {
				x, y      float64
				colliding bool
				note      string
			}{
				{0, 0, false, "starting state no collisions"},
				{8, 0, true, "right line basic"},
				{0, 8, true, "bottom line basic"},
				{100, 100, false, "diagonally moving move between gap between lines"},
			},
		},
		{
			name: "rectnagles and circle",
			shapes: []resolv.Shape{
				resolv.NewRectangle(0, 0, 5, 5),
				resolv.NewCircle(15, 15, 5),
			},
			resolves: []struct {
				x, y      float64
				colliding bool
				note      string
			}{
				{0, 0, false, "starting state no collisions"},
				{10, 10, true, "circle diagonally large"},
				{5, 5, false, "circle diagonally small"},
			},
		},
		{
			name: "rectnagles and space",
			shapes: []resolv.Shape{
				resolv.NewRectangle(0, 0, 5, 5),
				func() (result *resolv.Space) {
					result = resolv.NewSpace()
					result.Add(
						resolv.NewRectangle(10, 0, 10, 10),
						resolv.NewRectangle(0, 10, 10, 10),
					)
					return
				}(),
			},
			resolves: []struct {
				x, y      float64
				colliding bool
				note      string
			}{
				{0, 0, false, "starting state no collisions"},
				{10, 0, true, "space right"},
				{0, 10, true, "space bottom"},
			},
		},
	}

	for _, test := range testCases {
		for _, resolve := range test.resolves {
			space := resolv.NewSpace()
			space.Add(test.shapes...)

			collsion := space.Resolve(test.shapes[0], resolve.x, resolve.y)
			assert.Equal(t, resolve.colliding, collsion.Colliding(), resolve.note)

			// Make sure resolveX and ResolveY is is not causing another collsion
			collsion = space.Resolve(test.shapes[0], collsion.ResolveX, collsion.ResolveY)
			assert.False(t, collsion.Colliding(), resolve.note)
		}
	}
}

func createTestShapeWithTags(s resolv.Shape, tags ...string) resolv.Shape {
	s.AddTags(tags...)

	return s
}

func TestSpaceFilterTags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		shapes []resolv.Shape
		tests  []struct {
			token string
			count int
		}
	}{
		{
			shapes: []resolv.Shape{
				createTestShapeWithTags(resolv.NewRectangle(0, 0, 10, 10), "player"),
				createTestShapeWithTags(resolv.NewRectangle(0, 10, 10000, 10), "ground"),
				createTestShapeWithTags(resolv.NewLine(100, 0, 100, 10), "ground"),
				createTestShapeWithTags(resolv.NewCircle(1000, 0, 5), "token"),
			},
			tests: []struct {
				token string
				count int
			}{
				{"player", 1},
				{"ground", 2},
				{"token", 1},
			},
		},
	}

	for _, testCase := range testCases {
		space := resolv.NewSpace()
		space.Add(testCase.shapes...)
		for _, test := range testCase.tests {
			assert.Equal(t, test.count, space.FilterByTags(test.token).Length())
		}
	}
}

func TestSpaceAdd(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		spaceA := resolv.NewSpace()
		spaceA.Add(spaceA)
	})

	rectA := resolv.NewRectangle(0, 0, 0, 0)
	rectB := resolv.NewRectangle(0, 0, 0, 0)
	rectC := resolv.NewRectangle(0, 0, 0, 0)

	spaceA := resolv.NewSpace()
	spaceA.Add(rectA)
	spaceA.Add(rectB)
	spaceB := resolv.NewSpace()
	spaceB.Add(rectB)
	spaceB.Add(rectC)

	assert.True(t, spaceA.Contains(rectA))
	assert.True(t, spaceA.Contains(rectB))
	assert.False(t, spaceA.Contains(rectC))

	assert.False(t, spaceB.Contains(rectA))
	assert.True(t, spaceB.Contains(rectB))
	assert.True(t, spaceB.Contains(rectC))
}

func TestCollidingShapes(t *testing.T) {
	t.Parallel()

	rectB := resolv.NewRectangle(11, 0, 10, 10)
	rectC := resolv.NewRectangle(0, 11, 10, 10)

	spaceA := resolv.NewSpace()
	spaceA.Add(rectB, rectC)

	target := resolv.NewRectangle(0, 0, 10, 10)
	assert.Equal(t, 0, spaceA.GetCollidingShapes(target).Length())

	target.SetXY(2, 2)
	assert.Equal(t, 2, spaceA.GetCollidingShapes(target).Length())
}

func TestCollisions(t *testing.T) {
	t.Parallel()

	rectA := resolv.NewRectangle(0, 11, 10, 10)
	rectB := resolv.NewRectangle(0, 11, 10, 10)
	rectC := resolv.NewRectangle(0, 11, 10, 10)

	space := resolv.NewSpace()
	space.Add(rectA, rectB, rectC)

	assert.Equal(t, 2, len(space.Collisions(rectA)))
}
