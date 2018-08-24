package main

import (
	"fmt"
)

func timeString(seconds int32) string {
	var d int32
	var h int32
	var m int32

	d = seconds / 86400
	seconds = seconds % 86400

	h = seconds / 3600
	seconds = seconds % 3600

	m = seconds / 60

	return fmt.Sprintf("%dd %02d:%02d", d, h, m)
}
