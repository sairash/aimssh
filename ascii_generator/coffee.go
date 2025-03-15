// refrence from: https://www.youtube.com/watch?v=I5Q03I-ybXQ

package ascii_generator

import (
	"math"
	"time"
	"zencli/helper"

	"github.com/charmbracelet/lipgloss"
)

type Coffee struct {
	width, height, curFrame, prevPercent int
	cachedTop                            string
	cachedBottom                         string
}

func (c *Coffee) fill_amount(percent, width, height int) [][]rune {
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
	cup_color    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f78dbb"))
	coffee_color = lipgloss.NewStyle().Foreground(lipgloss.Color("#967259"))

	steam_animation = [][]string{
		{
			"       {", "    {   }",
			"                }" + cup_color.Render("_") + "{ " + cup_color.Render("__") + "{",
			cup_color.Render("             .-") + "{   }   }" + cup_color.Render("-."),
			cup_color.Render("            (") + "   }     {   " + cup_color.Render(")"),
		},
		{
			"        }", "   {   } ",
			"               { " + cup_color.Render("_") + "} " + cup_color.Render("__") + "{",
			cup_color.Render("             .-") + "}   {   {" + cup_color.Render("-."),
			cup_color.Render("            (") + "   {     }   " + cup_color.Render(")"),
		},
		{
			"     {", "  {   }",
			"                }" + cup_color.Render("_") + "{ " + cup_color.Render("__") + "{",
			cup_color.Render("             .-") + "{   }   }" + cup_color.Render("-."),
			cup_color.Render("            (") + "   }     {   " + cup_color.Render(")"),
		},
		{
			"          }", "     {   } ",
			"               { " + cup_color.Render("_") + "} " + cup_color.Render("__") + "{",
			cup_color.Render("             .-") + "}   {   {" + cup_color.Render("-."),
			cup_color.Render("            (") + "   {     }   " + cup_color.Render(")"),
		},
	}

	cup_head = cup_color.Render(helper.Center("|'-.._____..-'|", 40))

	cup_body = [][]string{
		{cup_color.Render("            | "), cup_color.Render(" ;--.")},
		{cup_color.Render("            | "), cup_color.Render(" (__  \\")},
		{cup_color.Render("            | "), cup_color.Render(" |  )  )")},
		{cup_color.Render("            | "), cup_color.Render(" | /  /")},
		{cup_color.Render("            | "), cup_color.Render(" |/  / ")},
		{cup_color.Render("            | "), cup_color.Render(" (  / ")},
		{cup_color.Render("            \\ "), cup_color.Render(" y` ")},
	}
	cup_bottom = cup_color.Render("             `-.._____..-'")
)

func (t *Coffee) Width() int {
	return t.width
}

func (t *Coffee) Height() int {
	return t.height
}

func (c *Coffee) NextAndString(percent int) string {
	c.Next(percent)
	return c.cachedTop + c.cachedBottom
}

func (c *Coffee) Next(percent int) bool {
	if c.prevPercent == percent {
		return false
	}

	grid := c.fill_amount(percent, c.width, c.height)

	total_body := cup_head + "\n"
	for k, v := range grid {
		total_body += cup_body[k][0] + coffee_color.Render(string(v)) + cup_body[k][1] + "\n"
	}
	total_body += cup_bottom
	c.cachedBottom = total_body + "\n\n\n\n"
	c.prevPercent = percent

	return true
}

func (c *Coffee) steam() int {
	c.curFrame++

	if c.curFrame >= len(steam_animation) {
		c.curFrame = 0
	}

	top := ""
	for _, v := range steam_animation[c.curFrame] {
		top += helper.Center(v, 40) + "\n"
	}
	c.cachedTop = "\n\n\n\n" + top

	return c.curFrame
}

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
	go func() {
		for {
			t := time.NewTicker(2 * time.Second)
			c.steam()
			<-t.C
		}
	}()

	return c
}
