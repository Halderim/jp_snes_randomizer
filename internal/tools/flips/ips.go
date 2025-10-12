package flips

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getBinaryPath() string {
	base := filepath.Join("internal", "tools", "flips")

	return filepath.Join(base, "flips")
}

func PatchIPS(ipsPath, romPath string, outPath string, logPath string) error {
	binary := getBinaryPath()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file: %w", err)
	}
	defer logFile.Close()

	logFile.WriteString("\n===== 📦 IPS Patching =====\n")
	logFile.WriteString(fmt.Sprintf("\n===== applying %s to %s =====\n", ipsPath, romPath))

	cmd := exec.Command(binary, "-a", ipsPath, romPath, outPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("flips failed for %s: %w", romPath, err)
	}
	return nil
}
