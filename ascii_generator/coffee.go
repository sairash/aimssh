// Coffee animation, inspiration from: https://www.youtube.com/watch?v=I5Q03I-ybXQ
package ascii_generator

import (
	"aimssh/helper"
	"math"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Coffee represents an animated coffee cup that fills over time
type Coffee struct {
	width        int
	height       int
	curFrame     int
	prevPercent  int
	cachedTop    string
	cachedBottom string
	mu           sync.RWMutex
}

// fillAmount calculates the fill pattern for the given percentage
func (c *Coffee) fillAmount(percent, width, height int) [][]rune {
	totalCells := width * height
	fillCells := int(math.Ceil(float64(percent*totalCells) / 100.0))

	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	remaining := fillCells
	for i := height - 1; i >= 0 && remaining > 0; i-- {
		for j := 0; j < width && remaining > 0; j++ {
			grid[i][j] = '#'
			remaining--
		}
	}

	return grid
}

var (
	cupColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f78dbb"))
	coffeeColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#967259"))

	steamAnimation = [][]string{
		{
			"       {", "    {   }",
			"                }" + cupColor.Render("_") + "{ " + cupColor.Render("__") + "{",
			cupColor.Render("             .-") + "{   }   }" + cupColor.Render("-."),
			cupColor.Render("            (") + "   }     {   " + cupColor.Render(")"),
		},
		{
			"        }", "   {   } ",
			"               { " + cupColor.Render("_") + "} " + cupColor.Render("__") + "{",
			cupColor.Render("             .-") + "}   {   {" + cupColor.Render("-."),
			cupColor.Render("            (") + "   {     }   " + cupColor.Render(")"),
		},
		{
			"     {", "  {   }",
			"                }" + cupColor.Render("_") + "{ " + cupColor.Render("__") + "{",
			cupColor.Render("             .-") + "{   }   }" + cupColor.Render("-."),
			cupColor.Render("            (") + "   }     {   " + cupColor.Render(")"),
		},
		{
			"          }", "     {   } ",
			"               { " + cupColor.Render("_") + "} " + cupColor.Render("__") + "{",
			cupColor.Render("             .-") + "}   {   {" + cupColor.Render("-."),
			cupColor.Render("            (") + "   {     }   " + cupColor.Render(")"),
		},
	}

	cupHead = cupColor.Render(helper.Center("|'-.._____..-'|", 40))

	cupBody = [][]string{
		{cupColor.Render("            | "), cupColor.Render(" ;--.")},
		{cupColor.Render("            | "), cupColor.Render(" (__  \\")},
		{cupColor.Render("            | "), cupColor.Render(" |  )  )")},
		{cupColor.Render("            | "), cupColor.Render(" | /  /")},
		{cupColor.Render("            | "), cupColor.Render(" |/  / ")},
		{cupColor.Render("            | "), cupColor.Render(" (  / ")},
		{cupColor.Render("            \\ "), cupColor.Render(" y` ")},
	}
	cupBottom = cupColor.Render("             `-.._____..-'")
)

// Width returns the width of the coffee fill area
func (c *Coffee) Width() int {
	return c.width
}

// Height returns the height of the coffee fill area
func (c *Coffee) Height() int {
	return c.height
}

// NextAndString updates the animation and returns the current frame as string
func (c *Coffee) NextAndString(percent int) string {
	c.Next(percent)
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cachedTop + c.cachedBottom
}

// Next updates the coffee fill level based on percentage
func (c *Coffee) Next(percent int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.prevPercent == percent {
		return false
	}

	grid := c.fillAmount(percent, c.width, c.height)

	totalBody := cupHead + "\n"
	for k, v := range grid {
		totalBody += cupBody[k][0] + coffeeColor.Render(string(v)) + cupBody[k][1] + "\n"
	}
	totalBody += cupBottom
	c.cachedBottom = totalBody + "\n\n\n"
	c.prevPercent = percent

	return true
}

// steam advances the steam animation and updates cachedTop
func (c *Coffee) steam() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.curFrame++
	if c.curFrame >= len(steamAnimation) {
		c.curFrame = 0
	}

	top := ""
	for _, v := range steamAnimation[c.curFrame] {
		top += helper.Center(v, 40) + "\n"
	}
	c.cachedTop = "\n\n\n" + top

	return c.curFrame
}

// GenerateCoffee creates a new Coffee animation
func GenerateCoffee() *Coffee {
	c := &Coffee{
		width:        11,
		height:       7,
		curFrame:     0,
		prevPercent:  -1,
		cachedTop:    "",
		cachedBottom: "",
	}

	c.Next(0)

	// Steam animation goroutine
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			c.steam()
			<-t.C
		}
	}()

	return c
}
