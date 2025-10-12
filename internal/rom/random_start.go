package rom

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"jp_snes_randomizer/internal/tools/flips"
)

func ApplyRandomStartLocation(romPath string, seed int64, logPath string) error {
	patchDir := "internal/patches/start/"
	index := 0
	r := rand.New(rand.NewSource(seed))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// read ips files in patchDir and save names to a map
	ipsFiles := make(map[int]string)
	err = filepath.Walk(patchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ips") {
			ipsFiles[index] = info.Name()
			index++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading IPS files: %w", err)
	}
	randomIndex := r.Intn(len(ipsFiles))
	selectedFile := ipsFiles[randomIndex]
	ipsPath := filepath.Join(patchDir, selectedFile)

	fmt.Fprintln(logFile, "\n===== 🧩 Start Location Randomizer Log =====")
	fmt.Fprintf(logFile, "Applying random start location patch: %s\n", selectedFile)
	fmt.Printf("Applying random start location patch: %s\n", selectedFile)

	if err := flips.PatchIPS(ipsPath, romPath, romPath, logPath); err != nil {
		return fmt.Errorf("error applying IPS patch: %w", err)
	}
	fmt.Printf("Patched ROM saved to: %s\n", romPath)
	return nil
}
