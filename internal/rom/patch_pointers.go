package rom

import (
	"fmt"
	"os"
)

// SNES → ROM Offset (LoROM)
func snesToRomOffset(addr int) int {
	bank := (addr >> 16) & 0x7F
	offset := addr & 0x7FFF
	return bank*0x8000 + offset
}

// romToSNESAddr converts a ROM-Offset into a 24-bit SNES-address (Bank<<16 | Addr)
// for LoROM: SNES-Bank = 0x80 + bankNum
func romToSNESAddr(romOffset int) int {
	bankNum := romOffset / 0x8000
	snesBank := 0x80 + bankNum
	addr := (romOffset % 0x8000) + 0x8000
	return (snesBank << 16) | addr
}

// PatchPointers updates pointers to the new RNC locations in the rom
func PatchPointers(romPath string, pointerList []PointerEntry, newStarts map[string]int, logPath string) error {
	data, err := os.ReadFile(romPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Lesen des ROMs: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der Logdatei: %w", err)
	}
	defer logFile.Close()

	fmt.Fprintln(logFile, "\n===== 🔗 Pointer Patch Log =====")

	for _, entry := range pointerList {
		newStart, ok := newStarts[entry.Filename]
		if !ok {
			fmt.Fprintf(logFile, "⚠️ No new start found for %s\n", entry.Filename)
			continue
		}

		// SNES-Address to ROM-Offsets
		hiOff := snesToRomOffset(entry.Hi)
		loOff := snesToRomOffset(entry.Lo)
		bankOff := snesToRomOffset(entry.Bank)

		if hiOff >= len(data) || loOff >= len(data) || bankOff >= len(data) {
			fmt.Fprintf(logFile, "❌ Invalid offset ranges for %s (Hi=%X Lo=%X Bank=%X)\n",
				entry.Filename, hiOff, loOff, bankOff)
			continue
		}

		// ROM-Offset -> SNES-Adresse (0x80+ bank)
		snesAddr := romToSNESAddr(newStart)

		// Pointer-Bytes: Lo, Hi, Bank (LoROM little-endian storage)
		lo := byte(snesAddr & 0xFF)
		hi := byte((snesAddr >> 8) & 0xFF)
		bank := byte((snesAddr >> 16) & 0xFF)

		old := []byte{data[loOff], data[hiOff], data[bankOff]}

		data[loOff] = lo
		data[hiOff] = hi
		data[bankOff] = bank

		newPtr := []byte{lo, hi, bank}

		// Log
		line := fmt.Sprintf("%-35s | OLD %02X %02X %02X → NEW %02X %02X %02X (ROM:%06X → SNES:%06X)\n",
			entry.Filename,
			old[0], old[1], old[2],
			newPtr[0], newPtr[1], newPtr[2],
			newStart, snesAddr,
		)
		fmt.Print(line)
		logFile.WriteString(line)
	}

	if err := os.WriteFile(romPath, data, 0644); err != nil {
		return fmt.Errorf("Error during writing ROM: %w", err)
	}

	fmt.Fprintln(logFile, "✅ Pointer successfully updated.\n")
	return nil
}
