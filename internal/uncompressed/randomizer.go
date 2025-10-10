// internal/uncompressed/randomizer.go
package uncompressed

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// RandomizeCards mischt die ID-Karten (außer Hammond) und schreibt Änderungen mit Logausgabe
func RandomizeCards(binDir string, outDir string, seed int64, logPath string) error {
	r := rand.New(rand.NewSource(seed))

	indices := make([]byte, len(CardLocations))
	for i, loc := range CardLocations {
		indices[i] = loc.Index
	}
	r.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	// Logdatei im gleichen Stil wie der ROM-Patcher
	log := filepath.Join(logPath, fmt.Sprintf("randomizer_log_seed%d.log", seed))
	logFile, err := os.OpenFile(log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	fmt.Fprintln(logFile, "\n===== 🧩 Item Randomizer Log =====")
	fmt.Fprintf(logFile, "Seed: %d | %s\n\n", seed, time.Now().Format("2006-01-02 15:04:05"))

	for i, loc := range CardLocations {
		input := filepath.Join(binDir, loc.File)
		output := filepath.Join(outDir, loc.File)

		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("Error reading %s: %v", input, err)
		}

		if loc.Offset >= len(data) {
			return fmt.Errorf("Offset %X out of bounds for file %s", loc.Offset, loc.File)
		}

		oldVal := data[loc.Offset]
		newVal := indices[i]
		data[loc.Offset] = newVal

		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(output, data, 0644); err != nil {
			return err
		}

		oldName := CardNames[oldVal]
		newName := CardNames[newVal]
		line := fmt.Sprintf("%-40s | Offset 0x%X | OLD %02X (%-12s) → NEW %02X (%-12s)\n",
			loc.File, loc.Offset, oldVal, oldName, newVal, newName)
		fmt.Print(line)
		logFile.WriteString(line)
	}

	fmt.Fprintln(logFile, "✅ All cards successfully randomized.\n")
	fmt.Println("📄 Log written to:", logPath)
	return nil
}
