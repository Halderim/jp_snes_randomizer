package rom

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RNCBlock presents a .rnc file with its data and size.
type RNCBlock struct {
	Name string
	Data []byte
	Size int
}

// AppendRNCsCompact reads all .rnc files, optionally sorts them by size
// and writes them compactly & bank-aligned into the ROM.
// Returns a map with Filename -> ROM-Offset.
func AppendRNCsCompact(rncDir string, romPath string, startOffset int, logPath string, sortBySize bool) (map[string]int, error) {
	data, err := os.ReadFile(romPath)
	if err != nil {
		return nil, fmt.Errorf("Error during reading ROM: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error during creating log file: %w", err)
	}
	defer logFile.Close()

	fmt.Fprintf(logFile, "\n\n===== 🦖 RNC Compact Embed Log — %s =====\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(logFile, "ROM: %s\n", romPath)
	fmt.Fprintf(logFile, "Startoffset: 0x%06X\n\n", startOffset)

	// === 1️⃣ Alle .rnc-Dateien einlesen ===
	var blocks []RNCBlock
	err = filepath.WalkDir(rncDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".rnc" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error during reading %s: %w", path, err)
		}
		blocks = append(blocks, RNCBlock{
			Name: filepath.Base(path),
			Data: content,
			Size: len(content),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("No RNC files found in %s", rncDir)
	}

	// === Sort by size ===
	if sortBySize {
		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].Size < blocks[j].Size // smaller first
		})
	} else {
		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].Name < blocks[j].Name // alphabetical
		})
	}

	// === Write into rom ===
	newStarts := make(map[string]int)
	currentOffset := startOffset

	for _, blk := range blocks {
		// Bank boundary alignment (LoROM)
		if currentOffset%0x8000 != 0 {
			pad := 0x8000 - (currentOffset % 0x8000)
			currentOffset += pad
			fmt.Fprintf(logFile, "🔧 Padding till next bank boundary: +0x%X\n", pad)
		}

		neededSize := currentOffset + blk.Size
		if len(data) < neededSize {
			padding := make([]byte, neededSize-len(data))
			data = append(data, padding...)
		}

		copy(data[currentOffset:], blk.Data)
		newStarts[blk.Name] = currentOffset

		fmt.Fprintf(logFile, "📦 %s @ 0x%06X (Size: %d Bytes)\n", blk.Name, currentOffset, blk.Size)
		currentOffset += blk.Size
	}

	if err := os.WriteFile(romPath, data, 0644); err != nil {
		return nil, fmt.Errorf("Error during writing ROM: %w", err)
	}

	fmt.Fprintf(logFile, "\n✅ %d RNCs successfully embedded.\nNew ROM: %s\n", len(blocks), romPath)
	fmt.Printf("✅ %d RNC files successfully embedded.\n", len(blocks))

	return newStarts, nil
}
