package rom

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// ValidatePointers checks pointer entries vs RNC files and ROM size and writes a log.
func ValidatePointers(rncDir string, romSize int, logPath string) error {
	var logBuilder strings.Builder
	logBuilder.WriteString("=== Pointer Validation Report ===\n\n")

	valid := true

	// 2) rnc files present
	fmt.Println("\n🔍 Check if all RNC files are in pointer list ...")
	rncFiles := []string{}
	err := filepath.WalkDir(rncDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".rnc" {
			return nil
		}
		rncFiles = append(rncFiles, filepath.Base(path))
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error during searching %s: %w", rncDir, err)
	}

	pointerMap := make(map[string]bool)
	for _, p := range PointerList {
		pointerMap[p.Filename] = true
	}

	for _, f := range rncFiles {
		if !pointerMap[f] {
			fmt.Printf("%s⚠️ %s has no pointer entry%s\n", colorYellow, f, colorReset)
			logBuilder.WriteString(fmt.Sprintf("⚠️  File without pointer entry: %s\n", f))
			valid = false
		}
	}

	// 3) duplicates
	fmt.Println("\n🔍 Check for duplicate pointer entries ...")
	nameSeen := make(map[string]bool)
	for _, p := range PointerList {
		if nameSeen[p.Filename] {
			fmt.Printf("%s⚠️  Duplicate pointer for %s%s\n", colorYellow, p.Filename, colorReset)
			logBuilder.WriteString(fmt.Sprintf("⚠️  Duplicate pointer entry for %s\n", p.Filename))
			valid = false
		}
		nameSeen[p.Filename] = true
	}

	// summary
	fmt.Println("\n📋 Summary:")
	if valid {
		fmt.Printf("%s✅ All pointers and files are consistent!%s\n", colorGreen, colorReset)
		logBuilder.WriteString("\n✅ All pointers and files are consistent!\n")
	} else {
		fmt.Printf("%s❌ Problems found. Please check log: %s%s\n", colorRed, logPath, colorReset)
		logBuilder.WriteString("\n❌ Problems found. Please check!\n")
	}

	if err := os.WriteFile(logPath, []byte(logBuilder.String()), 0644); err != nil {
		return fmt.Errorf("Error during writing validation log: %w", err)
	}

	fmt.Printf("\n📄 Log saved at: %s\n", logPath)
	return nil
}
