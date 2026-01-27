// Tree animation, inspiration: https://gitlab.com/jallbrit/cbonsai
package ascii_generator

import (
	"aimssh/helper"
	"math"
	"math/rand"
	"strings"
	"sync"
)

// ANSI color codes for tree rendering
const (
	branchColor     = "\033[38;5;94m"
	orangeColor     = "\033[38;5;208m"
	lightSkinColor  = "\033[38;5;223m"
	whiteColor      = "\033[38;5;15m"
	darkRedColor    = "\033[38;5;124m"
	lightBrownColor = "\033[38;5;137m"
	skyblueColor    = "\033[38;5;69m"
	navyBlueColor   = "\033[38;5;4m"
	deepBlueColor   = "\033[38;5;39m"
	leafColor       = "\033[32m"
	resetColor      = "\033[0m"
	skipHeight      = 2
)

// branchParams holds parameters for drawing a branch
type branchParams struct {
	x, y, angle, length float64
	depth               int
}

// Tree represents a procedurally generated ASCII tree
type Tree struct {
	Canvas        [][]helper.Cell
	width         int
	height        int
	branchesQueue []branchParams
	currentLevel  int
	totalLevels   int
	cachedOutput  string
	dirty         bool
	lastPercent   int
	mu            sync.Mutex
}

// initCanvas creates a new Tree with an initialized canvas
func initCanvas(width, height int) Tree {
	t := Tree{
		height: height,
		width:  width,
		dirty:  true,
	}
	t.Canvas = make([][]helper.Cell, height)
	for i := 0; i < height; i++ {
		t.Canvas[i] = make([]helper.Cell, width)
		appendChar := ' '
		color := ""
		if i == t.height-skipHeight {
			appendChar = '#'
			color = branchColor
		}
		for j := 0; j < width; j++ {
			t.Canvas[i][j] = helper.Cell{Ch: appendChar, Color: color}
		}
	}

	return t
}

// drawLine draws a line on the canvas using Bresenham's algorithm
func (t *Tree) drawLine(x0, y0, x1, y1 float64, ch rune, col string) {
	dx := x1 - x0
	dy := y1 - y0
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))

	for i := 0; i <= steps; i++ {
		x := x0 + dx*float64(i)/float64(steps)
		y := y0 + dy*float64(i)/float64(steps)
		ix := int(math.Round(x))
		iy := int(math.Round(y))

		if ix >= 0 && ix < t.width && iy >= 0 && iy < t.height-skipHeight {
			// Only mark dirty if cell actually changes
			if t.Canvas[iy][ix].Ch != ch || t.Canvas[iy][ix].Color != col {
				t.Canvas[iy][ix] = helper.Cell{Ch: ch, Color: col}
				t.dirty = true
			}
		}
	}
}

// Next processes the tree growth up to the given percentage
func (t *Tree) Next(percentage int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Return cached result if percentage hasn't changed
	if percentage == t.lastPercent && !t.dirty {
		return len(t.branchesQueue) > 0
	}

	t.lastPercent = percentage
	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}

	desiredLevels := (percentage * t.totalLevels) / 100
	levelsToProcess := desiredLevels - t.currentLevel
	if levelsToProcess <= 0 {
		return len(t.branchesQueue) > 0
	}

	processedAny := false
	for i := 0; i < levelsToProcess; i++ {
		if len(t.branchesQueue) == 0 {
			break
		}

		nextQueue := make([]branchParams, 0, len(t.branchesQueue)*2)
		for _, bp := range t.branchesQueue {
			if bp.depth == 0 || bp.length < 1 {
				ix, iy := int(math.Round(bp.x)), int(math.Round(bp.y))
				if ix >= 0 && ix < t.width && iy >= 0 && iy < t.height {
					// Only mark dirty if cell changes
					if t.Canvas[iy][ix].Ch != '#' || t.Canvas[iy][ix].Color != leafColor {
						t.Canvas[iy][ix] = helper.Cell{Ch: '#', Color: leafColor}
						t.dirty = true
					}
				}
				continue
			}

			xEnd := bp.x + bp.length*math.Cos(bp.angle)
			yEnd := bp.y - bp.length*math.Sin(bp.angle)

			var ch rune
			diff := math.Abs(bp.angle - math.Pi/2)
			if diff < 0.2 {
				ch = '|'
			} else if bp.angle < math.Pi/2 {
				ch = '/'
			} else {
				ch = '\\'
			}

			t.drawLine(bp.x, bp.y, xEnd, yEnd, ch, branchColor)

			newLength := bp.length * (0.7 + 0.1*rand.Float64())
			leftAngle := bp.angle + (0.3 + 0.2*rand.Float64())
			rightAngle := bp.angle - (0.3 + 0.2*rand.Float64())

			nextDepth := bp.depth - 1
			nextQueue = append(nextQueue, branchParams{
				x:      xEnd,
				y:      yEnd,
				angle:  leftAngle,
				length: newLength,
				depth:  nextDepth,
			}, branchParams{
				x:      xEnd,
				y:      yEnd,
				angle:  rightAngle,
				length: newLength,
				depth:  nextDepth,
			})
			processedAny = true
		}

		t.branchesQueue = nextQueue
		t.currentLevel++
	}

	if processedAny {
		t.dirty = true
	}
	return len(t.branchesQueue) > 0
}

// StringPrint returns the tree as a string
func (t *Tree) StringPrint() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.dirty && t.cachedOutput != "" {
		return t.cachedOutput
	}

	var builder strings.Builder
	builder.Grow(t.height * (t.width + 1)) // Pre-allocate memory

	for i := 0; i < t.height; i++ {
		line := make([]byte, 0, t.width)
		for j := 0; j < t.width; j++ {
			cell := t.Canvas[i][j]
			if cell.Color != "" {
				line = append(line, cell.Color...)
				line = append(line, byte(cell.Ch))
				line = append(line, resetColor...)
			} else {
				line = append(line, byte(cell.Ch))
			}
		}
		builder.Write(line)
		builder.WriteByte('\n')
	}

	t.cachedOutput = builder.String()
	t.dirty = false
	return t.cachedOutput
}

// NextAndString updates the tree and returns it as a string
func (t *Tree) NextAndString(percent int) string {
	t.Next(percent)
	return t.StringPrint() + "\n"
}

// Width returns the canvas width
func (t *Tree) Width() int {
	return t.width
}

// Height returns the canvas height
func (t *Tree) Height() int {
	return t.height
}

// GenerateTree creates a new procedural tree with the given dimensions
func GenerateTree(width, height int) *Tree {
	t := initCanvas(width, height)

	startX := float64(width / 2)
	startY := float64(height - 1)
	initialLength := float64(height) / 4

	t.totalLevels = 10 // Corresponds to the initial depth of 10
	t.branchesQueue = []branchParams{
		{
			x:      startX,
			y:      startY,
			angle:  math.Pi / 2,
			length: initialLength,
			depth:  10,
		},
	}

	return &t
}
