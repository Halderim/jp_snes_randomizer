package rom

import (
	"fmt"
	"os"

	"jp_snes_randomizer/internal/tools/flips"
)

// PatchIntro puts new intro into rom and patches the rom to point to new text
func PatchIntro(romPath string, logPath string) error {
	offset := 0x200000
	introBinPath := "internal/patches/intro/intro.bin"
	introIPSPath := "internal/patches/intro/intro.ips"

	data, err := os.ReadFile(romPath)
	if err != nil {
		return fmt.Errorf("Error loading ROMs: %w", err)
	}

	intro, err := os.ReadFile(introBinPath)
	if err != nil {
		return fmt.Errorf("Error loading intro file: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Error opening logfile: %w", err)
	}
	defer logFile.Close()

	fmt.Fprintln(logFile, "\n===== 🔗 Patching Intro =====")

	copy(data[offset:], intro)

	if err := os.WriteFile(romPath, data, 0644); err != nil {
		return fmt.Errorf("Error during writing ROM: %w", err)
	}

	if err := flips.PatchIPS(introIPSPath, romPath, romPath, logPath); err != nil {
		return fmt.Errorf("error applying IPS patch: %w", err)
	}

	fmt.Fprintln(logFile, "✅ Intro Updated updated.\n")
	return nil
}
