package rom

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
)

const (
	OverworldStart = 0xBA6F
	OverworldEnd   = 0xC639
	OWRecordSize   = 7
)

// item types defined by ptr high byte
var overworldItemTypes = map[uint16]string{
	0x0000: "Shotgun",
	0x000C: "Gas",
	0x0018: "Tranquilizer",
	0x0024: "Bola",
	0x0030: "Rockets",
	0x003C: "Egg",
	0x0054: "Medikit",
	0x0060: "HammondID",
}

// structure used internally
type owItem struct {
	Index    int
	Offset   int
	Ptr      uint16
	XY       uint16
	Extra    uint16
	Z        byte
	ItemType string
	DecodedX uint16
	DecodedY uint16
}

func decodeX(xy uint16) uint16 {
	lo := uint16(xy & 0x00FF)
	lo <<= 4
	lo += 4
	return lo
}

func decodeY(xy uint16) uint16 {
	hi := uint16((xy & 0xFF00))
	hi >>= 4
	return hi
}

// -----------------------------------------------------------------------------
// PUBLIC API
// -----------------------------------------------------------------------------

// RandomizeOverworldItems shuffles XY/Extra/Z of all overworld items EXCEPT Hammond's ID.
// The ROM is modified in-place. Logging output is appended to logSlice.
// The same seed always produces the same output.
func RandomizeOverworldItems(romPath string, seed int64, logPath string) error {

	rom, err := os.ReadFile(romPath)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Error opening logfile: %w", err)
	}
	defer logFile.Close()

	if len(rom) < OverworldEnd {
		return fmt.Errorf("ROM too small")
	}

	// load all items
	items, err := loadOverworldItems(rom)
	if err != nil {
		return err
	}

	var shuffleIndices []int
	for i, it := range items {
		// Only shuffle real items, never unknown objects, never HammondID
		if it.ItemType != "Unknown" && it.ItemType != "HammondID" {
			shuffleIndices = append(shuffleIndices, i)
		}
	}

	// deterministic RNG
	rng := rand.New(rand.NewSource(seed))

	// extract all (XY, Extra, Z)
	type pos struct {
		XY    uint16
		Extra uint16
		Z     byte
	}
	positions := make([]pos, len(shuffleIndices))
	for i, idx := range shuffleIndices {
		positions[i] = pos{
			XY:    items[idx].XY,
			Extra: items[idx].Extra,
			Z:     items[idx].Z,
		}
	}

	// shuffle positions
	rng.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	// write back shuffled data
	for i, idx := range shuffleIndices {
		items[idx].XY = positions[i].XY
		//items[idx].Extra = positions[i].Extra
		items[idx].Z = positions[i].Z
		items[idx].DecodedX = decodeX(positions[i].XY)
		items[idx].DecodedY = decodeY(positions[i].XY)
	}

	// apply to ROM
	rom, err = writeOverworldItems(rom, items)
	if err != nil {
		return err
	}

	if err := os.WriteFile(romPath, rom, 0644); err != nil {
		return fmt.Errorf("failed to write modified ROM: %w", err)
	}
	// logging

	for _, it := range items {
		fmt.Fprintf(logFile,
			"[Overworld] idx=%d type=%s ptr=%04X XY=%04X Extra=%04X Z=%02X X=%d Y=%d \n",
			it.Index, it.ItemType, it.Ptr, it.XY, it.Extra, it.Z, it.DecodedX, it.DecodedY,
		)
	}

	return nil
}

// -----------------------------------------------------------------------------
// INTERNALS
// -----------------------------------------------------------------------------

func loadOverworldItems(rom []byte) ([]owItem, error) {
	items := []owItem{}
	offset := OverworldStart
	index := 0

	for offset+OWRecordSize <= OverworldEnd {
		// terminator: FF FF
		if rom[offset] == 0xFF && rom[offset+1] == 0xFF {
			break
		}

		ptr := binary.LittleEndian.Uint16(rom[offset : offset+2])
		xy := binary.LittleEndian.Uint16(rom[offset+2 : offset+4])
		extra := binary.LittleEndian.Uint16(rom[offset+4 : offset+6])
		z := rom[offset+6]

		tp := overworldItemTypes[ptr&0xFFFF]
		if tp == "" {
			tp = "Unknown"
		}

		item := owItem{
			Index:    index,
			Offset:   offset,
			Ptr:      ptr,
			XY:       xy,
			Extra:    extra,
			Z:        z,
			ItemType: tp,
			DecodedX: decodeX(xy),
			DecodedY: decodeY(xy),
		}
		fmt.Println(item)
		items = append(items, item)
		offset += OWRecordSize
		index++
	}

	return items, nil
}

func writeOverworldItems(rom []byte, items []owItem) ([]byte, error) {
	for _, it := range items {
		offset := it.Offset
		if offset+OWRecordSize > len(rom) {
			return nil, fmt.Errorf("ROM offset overflow at item %d", it.Index)
		}

		binary.LittleEndian.PutUint16(rom[offset:offset+2], it.Ptr)
		binary.LittleEndian.PutUint16(rom[offset+2:offset+4], it.XY)
		binary.LittleEndian.PutUint16(rom[offset+4:offset+6], it.Extra)
		rom[offset+6] = it.Z
	}
	return rom, nil
}
