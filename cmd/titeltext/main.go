package main

import (
	"flag"
	"fmt"
)

// wandelt Buchstaben in Tile-Index um
// A=00, B=01, ..., Z=19
func charToTile(ch rune) (byte, bool) {
	switch {
	case ch >= 'A' && ch <= 'Z':
		return byte(ch - 'A'), true
	case ch >= 'a' && ch <= 'z':
		return byte(ch - 'a'), true
	case ch >= '0' && ch <= '9':
		// 0 = 0x1B, 9 = 0x24
		return byte(0x1A + (ch - '0')), true
	case ch == '.':
		return 0x24, true
	default:
		return 0, false
	}
}

func main() {
	text := flag.String("text", "Randomizer by Halderim", "Text to convert")
	startX := flag.Int("x", 0x20, "Start X value (hex or dec)")
	startY := flag.Int("y", 0x20, "Start Y value (hex or dec)")
	prop := flag.Int("prop", 0x3C, "Tile property byte")
	step := flag.Int("step", 0x08, "X increment per letter (normally 8)")
	space := flag.Int("space", 0x08, "X increment for space (normally 16)")
	flag.Parse()

	x := *startX
	y := *startY

	fmt.Printf("Text: %q\n", *text)
	fmt.Println("Result (X Y Tile Prop):")

	for _, ch := range *text {
		if ch == ' ' {
			x += *space
			continue
		}
		tile, ok := charToTile(ch)
		if !ok {
			continue
		}
		fmt.Printf("%02X %02X %02X %02X\n", x, y, tile, *prop)
		x += *step
	}
}
