// Row Animation, inspiration: https://stonestoryrpg.com/
package ascii_generator

import (
	"aimssh/helper"
	"math/rand"
	"sync"
	"time"
)

// Row represents an animated rowing figure moving across the screen
type Row struct {
	cached       string
	background   [][]helper.Cell
	curAnimation int
	x, y         int
	width        int
	height       int
	mu           sync.RWMutex
}

var rowAnimationCell = [][][]helper.Cell{
	{
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'o', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '|', Color: lightBrownColor}},
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '\\', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '|', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}},
		{{Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '(', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '\\', Color: lightSkinColor}, {Ch: '|', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '`', Color: lightBrownColor}, {Ch: ')', Color: lightBrownColor}}},
	{
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'o', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightSkinColor}},
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '\\', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '(', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '\\', Color: lightSkinColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '`', Color: lightBrownColor}, {Ch: ')', Color: lightBrownColor}},
		{{Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '/', Color: lightBrownColor}, {Ch: '_', Color: navyBlueColor}}},
	{
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'o', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'V', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '(', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '`', Color: lightBrownColor}, {Ch: ')', Color: lightBrownColor}},
		{{Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '~', Color: navyBlueColor}}},
	{
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'o', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '(', Color: lightSkinColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ',', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '(', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '`', Color: lightBrownColor}, {Ch: ')', Color: lightBrownColor}},
		{{Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '~', Color: navyBlueColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}}},
	{
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: 'o', Color: lightSkinColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '<', Color: lightSkinColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: ',', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}},
		{{Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '(', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '/', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: '_', Color: lightBrownColor}, {Ch: ' ', Color: lightBrownColor}, {Ch: '`', Color: lightBrownColor}, {Ch: ')', Color: lightBrownColor}},
		{{Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: '_', Color: navyBlueColor}, {Ch: ' ', Color: navyBlueColor}, {Ch: '/', Color: lightBrownColor}},
	},
}

// NextAndString returns the current cached frame
func (r *Row) NextAndString(percent int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cached
}

// Width returns the width of the animation area
func (r *Row) Width() int {
	return r.width
}

// Height returns the height of the animation area
func (r *Row) Height() int {
	return r.height
}

// nextFrame advances to the next animation frame and returns its index
func (r *Row) nextFrame() int {
	r.curAnimation++
	if r.curAnimation >= len(rowAnimationCell) {
		r.curAnimation = 0
	}
	return r.curAnimation
}

// randomBlue returns a random blue color for water effect
func randomBlue() string {
	x := rand.Float32()
	if x >= 0.9 {
		return navyBlueColor
	} else if x > 0.6 {
		return deepBlueColor
	}
	return skyblueColor
}

// backgroundCreator generates a new random water background
func (r *Row) backgroundCreator() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.background = make([][]helper.Cell, r.height)
	for i := 0; i < r.height; i++ {
		cell := make([]helper.Cell, r.width)
		for j := 0; j < r.width; j++ {
			if rand.Float32() > 0.92 {
				cell[j] = helper.Cell{Ch: '_', Color: randomBlue()}
			} else {
				cell[j] = helper.Cell{Ch: ' '}
			}
		}
		r.background[i] = cell
	}
}

// GenerateRow creates a new Row animation with the given dimensions
func GenerateRow(width, height int) *Row {
	r := &Row{
		x:      -30,
		y:      rand.Intn(height - 4),
		width:  width,
		height: height,
	}

	r.backgroundCreator()

	// Background regeneration goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			r.backgroundCreator()
		}
	}()

	// Animation frame update goroutine
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			r.mu.Lock()
			frame := r.nextFrame()
			r.cached = helper.LayerCanvas(r.background, rowAnimationCell[frame], r.x, r.y, true)
			r.x++
			if r.x > r.width {
				r.x = -30
				r.y = rand.Intn(height - 4)
			}
			r.mu.Unlock()
		}
	}()

	return r
}
