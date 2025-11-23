package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
)

type Entry struct {
	Pointer uint16
	XY      uint16
	Extra   uint16
	Z       uint8
	Xcoord  uint16
	Ycoord  uint16
}

func DecodeX(xy uint16) uint16 {
	x := uint16(xy & 0x00FF)
	x <<= 4
	x += 4
	return x
}

func DecodeY(xy uint16) uint16 {
	y := uint16((xy & 0xFF00))
	y >>= 4
	return y
}

func main() {
	data, err := os.ReadFile("overworld_itemtable.bin")
	if err != nil {
		panic(err)
	}

	entries := []Entry{}

	for i := 0; i+7 <= len(data); i += 7 {
		// Terminator FF FF
		if data[i] == 0xFF && data[i+1] == 0xFF {
			break
		}

		ptr := binary.LittleEndian.Uint16(data[i : i+2])
		xy := binary.LittleEndian.Uint16(data[i+2 : i+4])
		extra := binary.LittleEndian.Uint16(data[i+4 : i+6])
		z := data[i+6]

		entry := Entry{
			Pointer: ptr,
			XY:      xy,
			Extra:   extra,
			Z:       z,
			Xcoord:  DecodeX(xy),
			Ycoord:  DecodeY(xy),
		}

		entries = append(entries, entry)
	}

	// Terminalausgabe
	for i, e := range entries {
		fmt.Printf(
			"%03d  PTR:%04X  XY:%04X  X:%04X  Y:%04X  EXTRA:%04X  Z:%02X\n",
			i, e.Pointer, e.XY, e.Xcoord, e.Ycoord, e.Extra, e.Z,
		)
	}

	// CSV schreiben
	csvFile, err := os.Create("overworld_items.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Header
	writer.Write([]string{
		"Index",
		"PointerHex",
		"XYHex",
		"ExtraHex",
		"ZHex",
		"X_Decoded",
		"Y_Decoded",
	})

	// Daten
	for i, e := range entries {
		writer.Write([]string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%04X", e.Pointer),
			fmt.Sprintf("%04X", e.XY),
			fmt.Sprintf("%04X", e.Extra),
			fmt.Sprintf("%02X", e.Z),
			fmt.Sprintf("%d", e.Xcoord),
			fmt.Sprintf("%d", e.Ycoord),
		})
	}
}
