package ascii_generator

import (
	"math"
	"math/rand"
)

const (
	branchColor = "\033[38;5;94m" // Brown-ish color for branches.
	leafColor   = "\033[32m"      // Green color for leaves.
	resetColor  = "\033[0m"       // Reset color.
)

func initCanvas(width, height int) AsciiArt {
	tree := AsciiArt{
		Height: height,
		Width:  width,
	}
	tree.Canvas = make([][]Cell, height)
	for i := 0; i < height; i++ {
		tree.Canvas[i] = make([]Cell, width)
		for j := 0; j < width; j++ {
			tree.Canvas[i][j] = Cell{' ', ""}
		}
	}
	return tree
}

func (t AsciiArt) drawLine(x0, y0, x1, y1 float64, ch rune, col string) {
	dx := x1 - x0
	dy := y1 - y0
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	for i := 0; i <= steps; i++ {
		x := x0 + dx*float64(i)/float64(steps)
		y := y0 + dy*float64(i)/float64(steps)
		ix := int(math.Round(x))
		iy := int(math.Round(y))
		if ix >= 0 && ix < t.Width && iy >= 0 && iy < t.Height {
			t.Canvas[iy][ix] = Cell{ch, col}
		}
	}
}

func (t AsciiArt) drawBranch(x, y, angle, length float64, depth int) {
	if depth == 0 || length < 1 {
		ix, iy := int(math.Round(x)), int(math.Round(y))
		if ix >= 0 && ix < t.Width && iy >= 0 && iy < t.Height {
			t.Canvas[iy][ix] = Cell{'#', leafColor}
		}
		return
	}

	xEnd := x + length*math.Cos(angle)
	yEnd := y - length*math.Sin(angle)

	var ch rune
	diff := math.Abs(angle - math.Pi/2)
	if diff < 0.2 {
		ch = '|'
	} else if angle < math.Pi/2 {
		ch = '/'
	} else {
		ch = '\\'
	}

	t.drawLine(x, y, xEnd, yEnd, ch, branchColor)

	newLength := length * (0.7 + 0.1*rand.Float64())
	leftAngle := angle + (0.3 + 0.2*rand.Float64())
	rightAngle := angle - (0.3 + 0.2*rand.Float64())

	t.drawBranch(xEnd, yEnd, leftAngle, newLength, depth-1)
	t.drawBranch(xEnd, yEnd, rightAngle, newLength, depth-1)
}

func (t AsciiArt) StringPrint() string {
	return_string := ""

	for i := 0; i < t.Height; i++ {
		for j := 0; j < t.Width; j++ {
			cell := t.Canvas[i][j]
			if cell.Color != "" {
				return_string += cell.Color + string(cell.Ch) + resetColor
			} else {
				return_string += string(cell.Ch)
			}
		}
		return_string += "\n"
	}

	return return_string
}

func (t AsciiArt) StringArray() []string {
	ret_string_arr := []string{}

	for i := 0; i < t.Height; i++ {
		str := ""
		for j := 0; j < t.Width; j++ {
			cell := t.Canvas[i][j]
			if cell.Color != "" {
				str += cell.Color + string(cell.Ch) + resetColor
			} else {
				str += string(cell.Ch)
			}
		}
		ret_string_arr = append(ret_string_arr, str)
	}
	return ret_string_arr
}

// func printCanvas() {
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			cell := canvas[i][j]
// 			if cell.color != "" {
// 				fmt.Print(cell.color, string(cell.ch), resetColor)
// 			} else {
// 				fmt.Print(string(cell.ch))
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

func GenerateAsciiArt(width, height int) AsciiArt {
	t := initCanvas(width, height)

	startX := float64(width / 2)
	startY := float64(height - 1)
	initialLength := float64(height) / 4

	t.drawBranch(startX, startY, math.Pi/2, initialLength, 10)

	return t
}
