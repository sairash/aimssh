package helper

import "strings"

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
