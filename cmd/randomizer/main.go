package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"jp_snes_randomizer/internal/rom"
	"jp_snes_randomizer/internal/tools/rncpropack"

	"jp_snes_randomizer/internal/uncompressed"
)

func main() {
	binDir := flag.String("bin", "internal/uncompressed", "Input BIN dir")
	seed := flag.Int64("seed", 0, "Randomizer seed (0 = use current time)")
	start := flag.Bool("start", false, "Apply random start location")
	flag.Parse()

	finalSeed := *seed
	if finalSeed == 0 {
		finalSeed = time.Now().Unix()
	}

	outDir := filepath.Join("internal", "modbin", fmt.Sprintf("%d", finalSeed))
	unmodifiedRom := "internal/rom/unmodified/jp_usa_rev1_ex.sfc"
	outRomPath := "internal/rom/modified/"
	expandedRom := filepath.Join(outRomPath, fmt.Sprintf("jp_randomized_seed%d.sfc", finalSeed))
	logDir := "internal/logs/"
	logPath := filepath.Join(logDir, "randomizer.log")

	// === Ensure target directory exists ===
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// === Create ROM copy ===
	err := copyFile(unmodifiedRom, expandedRom)
	if err != nil {
		fmt.Println("Error during ROM copy:", err)
		return
	}

	fmt.Printf("📀 ROM Copie created: %s\n", expandedRom)

	fmt.Println("🚀 Running randomizer with seed:", finalSeed)

	// 1) Randomize (user-provided function)
	if err := uncompressed.RandomizeCards(*binDir, outDir, finalSeed, logDir); err != nil {
		log.Fatal("❌ Randomization failed:", err)
	}

	// 2) Repack RNCs into ROM expanded area
	fmt.Println("📀 Packing RNC Files...")
	rncpropack.RepackAll(outDir, logPath)
	fmt.Println("📀 Embedding RNCs into expanded ROM...")
	newStarts, err := rom.AppendRNCsCompact(outDir, expandedRom, 0x300000, logPath, true)
	if err != nil {
		log.Fatal("❌ Embedding RNCs failed:", err)
	}

	// 3) Validate pointers
	fmt.Println("🔍 Validating pointers...")
	if err := rom.ValidatePointers(outDir, 0x400000, logPath); err != nil {
		log.Fatal("❌ Pointer validation failed:", err)
	}

	// 4) Patch pointers in ROM
	fmt.Println("🔗 Patching pointers...")
	if err := rom.PatchPointers(expandedRom, rom.PointerList, newStarts, logPath); err != nil {
		log.Fatal("❌ Pointer patching failed:", err)
	}

	if *start {
		fmt.Println("🔗 Apply random start location...")
		if err := rom.ApplyRandomStartLocation(expandedRom, finalSeed, logPath); err != nil {
			log.Fatal("❌ Applying random start location failed:", err)
		}
	} else {
		fmt.Println("ℹ️  Skipping random start location.")
	}

	fmt.Println("🔗 Patch in QoL patches")
	if err := rom.ApplyQolPatches(expandedRom, finalSeed, logPath); err != nil {
		log.Fatal("❌ Applying QoL patches failed:", err)
	}

	// Done
	fmt.Println("✅ Done. See log:", logPath)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
