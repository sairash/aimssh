package helper

import "strings"

var (
	resetColor = "\033[0m"
)

const (
	TimerEndedTitle     = "Timer Ended"
	TimerStartedTitle   = "Timer Started"
	TimerRestartedTitle = "Timer Restarted"
)

type Cell struct {
	Ch    rune
	Color string
}

func Center(s string, total_width int) string {
	if len(s) >= total_width {
		return s
	}
	n := total_width - len(s)
	div := n / 2
	return strings.Repeat(" ", div) + s
}

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
