// Package helper provides utility functions for the AimSSH application,
// including string manipulation, canvas layering, and terminal commands.
package helper

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ANSI reset color code
const resetColor = "\033[0m"

// Notification title constants
const (
	TimerEndedTitle     = "Timer Ended"
	TimerStartedTitle   = "Timer Started"
	TimerRestartedTitle = "Timer Restarted"
)

// Cell represents a single character cell with optional color
type Cell struct {
	Ch    rune
	Color string
}

// Center centers a string within the given width by padding with spaces
func Center(s string, totalWidth int) string {
	if len(s) >= totalWidth {
		return s
	}
	n := totalWidth - len(s)
	div := n / 2
	return strings.Repeat(" ", div) + s
}

// LayerString overlays a string on top of a background string at position (x, y).
// If hideOverflow is true, parts of the top string that exceed the background bounds are hidden.
func LayerString(background, top string, x, y int, hideOverflow bool) string {
	bgLines := strings.Split(background, "\n")
	topLines := strings.Split(top, "\n")

	if y < 0 {
		trim := -y
		if trim >= len(topLines) {
			return background
		}
		topLines = topLines[trim:]
		y = 0
	}

	trimLeft := 0
	if x < 0 {
		trimLeft = -x
		x = 0
	}
	for i := range topLines {
		if trimLeft > len(topLines[i]) {
			topLines[i] = ""
		} else {
			topLines[i] = topLines[i][trimLeft:]
		}
	}

	for i, topLine := range topLines {
		targetLine := y + i

		for len(bgLines) <= targetLine {
			bgLines = append(bgLines, "")
		}

		original := bgLines[targetLine]
		var start string

		if len(original) >= x {
			start = original[:x]
		} else {
			start = original + strings.Repeat(" ", x-len(original))
		}

		availableSpace := len(original) - x
		if availableSpace < 0 {
			availableSpace = 0
		}

		if hideOverflow {
			if len(topLine) > availableSpace {
				if availableSpace <= 0 {
					topLine = ""
				} else {
					topLine = topLine[:availableSpace]
				}
			}
		}

		endStart := x + len(topLine)
		var end string
		if len(original) > endStart {
			end = original[endStart:]
		}

		bgLines[targetLine] = start + topLine + end
	}

	return strings.Join(bgLines, "\n")
}

// cellsToString converts a 2D grid of Cells to a string with ANSI colors
func cellsToString(grid [][]Cell) string {
	lines := make([]string, len(grid))
	for i, row := range grid {
		var sb strings.Builder
		for _, cell := range row {
			if cell.Ch == 0 {
				sb.WriteRune(' ')
			} else {
				sb.WriteString(cell.Color)
				sb.WriteRune(cell.Ch)
				sb.WriteString(resetColor)
			}
		}
		lines[i] = sb.String()
	}
	return strings.Join(lines, "\n") + "\n"
}

// LayerCanvas overlays a Cell grid on top of a background Cell grid at position (x, y).
// If hideOverflow is true, parts of the top canvas that exceed the background bounds are hidden.
func LayerCanvas(background, topCan [][]Cell, x, y int, hideOverflow bool) string {
	bg := make([][]Cell, len(background))
	for i := range background {
		bg[i] = make([]Cell, len(background[i]))
		copy(bg[i], background[i])
	}

	top := make([][]Cell, len(topCan))
	for i := range topCan {
		top[i] = make([]Cell, len(topCan[i]))
		copy(top[i], topCan[i])
	}

	if y < 0 {
		if trim := -y; trim < len(top) {
			top = top[trim:]
		} else {
			return cellsToString(bg)
		}
		y = 0
	}

	trimLeft := 0
	if x < 0 {
		trimLeft = -x
		x = 0
	}
	for i := range top {
		if trimLeft < len(top[i]) {
			top[i] = top[i][trimLeft:]
		} else {
			top[i] = []Cell{}
		}
	}

	for i, tline := range top {
		targetY := y + i

		for len(bg) <= targetY {
			bg = append(bg, []Cell{})
		}

		bgline := bg[targetY]
		available := len(bgline) - x
		if available < 0 {
			available = 0
		}

		visible := tline
		if hideOverflow && len(tline) > available {
			if available > 0 {
				visible = tline[:available]
			} else {
				visible = []Cell{}
			}
		}

		if !hideOverflow && x+len(visible) > len(bgline) {
			newLen := x + len(visible)
			if cap(bgline) >= newLen {
				bgline = bgline[:newLen]
			} else {
				newBgline := make([]Cell, newLen)
				copy(newBgline, bgline)
				bgline = newBgline
			}
		}

		if x > len(bgline) {
			padding := make([]Cell, x-len(bgline))
			for j := range padding {
				padding[j].Ch = ' '
			}
			bgline = append(bgline, padding...)
		}

		for j, cell := range visible {
			pos := x + j
			if pos < len(bgline) {
				bgline[pos] = cell
			} else if !hideOverflow {
				bgline = append(bgline, cell)
			}
		}

		bg[targetY] = bgline
	}

	return cellsToString(bg)
}

// BeepCmd returns a tea.Cmd that sends a terminal bell character
func BeepCmd() tea.Cmd {
	fmt.Print("\a")
	return nil
}
