package rncpropack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Automatically selects the correct binary (rnc32 / rnc64)
func getBinaryPath() string {
	base := filepath.Join("internal", "tools", "rncpropack")
	if runtime.GOARCH == "amd64" {
		return filepath.Join(base, "rnc64")
	}
	return filepath.Join(base, "rnc32")
}

// RepackBin packing .bin file into a .rnc file
func RepackBin(binPath, outDir string) (string, error) {
	binFile := filepath.Base(binPath)
	outFile := filepath.Join(outDir, binFile[:len(binFile)-4]+".rnc")

	binary := getBinaryPath()
	cmd := exec.Command(binary, "p", binPath, outFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("rncpropack failed for %s: %w", binPath, err)
	}
	return outFile, nil
}

// RepackAll processes all BIN files in the directory and converts them to RNCs
func RepackAll(binDir string, logPath string) error {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file: %w", err)
	}
	defer logFile.Close()

	logFile.WriteString("\n===== 📦 RNC Repack Log =====\n")
	entries, err := os.ReadDir(binDir)

	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".bin" {
			binPath := filepath.Join(binDir, e.Name())
			outFile, err := RepackBin(binPath, binDir)
			if err != nil {
				return err
			}
			fmt.Printf("Packed %s → %s\n", binPath, outFile)
			logFile.WriteString(fmt.Sprintf("Packed %s → %s\n", binPath, outFile))
		}
	}
	return nil
}
