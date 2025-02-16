package tui

import (
	"strings"
)

func DrawLine(width int) string {
	var msg strings.Builder

	for range width {
		msg.WriteString("┄")
	}

	return msg.String()
}
