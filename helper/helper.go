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
