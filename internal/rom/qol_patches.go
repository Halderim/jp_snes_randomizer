package rom

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"jp_snes_randomizer/internal/tools/flips"
)

func ApplyQolPatches(romPath string, seed int64, logPath string) error {
	patchDir := "internal/patches/qol/"

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// read ips files in patchDir and save names to a map
	fmt.Fprintln(logFile, "\n===== 🧩 Apply QoL Patches =====")
	err = filepath.Walk(patchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ips") {
			ipsPath := filepath.Join(patchDir, info.Name())
			fmt.Fprintf(logFile, "Applying QoL patch: %s\n", info.Name())
			fmt.Printf("Applying QoL patch: %s\n", info.Name())
			if err := flips.PatchIPS(ipsPath, romPath, romPath, logPath); err != nil {
				return fmt.Errorf("error applying IPS patch: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading IPS files: %w", err)
	}

	fmt.Printf("Patched ROM saved to: %s\n", romPath)
	return nil
}
