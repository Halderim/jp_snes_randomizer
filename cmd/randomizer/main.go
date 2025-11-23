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

	seed := flag.Int64("seed", 0, "Randomizer seed (0 = use current time)")
	start := flag.Bool("start", false, "Apply random start location")
	difficulty := flag.Int64("difficulty", 0, "Set game difficulty")
	overworld := flag.Bool("overworld", true, "Disable overworld randomization")
	flag.Parse()

	finalSeed := *seed
	if finalSeed == 0 {
		finalSeed = time.Now().Unix()
	}

	// === 2️⃣ Prepare bin files===
	srcBinDir := "internal/uncompressed"
	outDir := fmt.Sprintf("internal/modbin/%d", *seed)

	fmt.Println("📂 Kopiere BIN-Dateien für Seed-Arbeitsverzeichnis...")
	if err := copyBinFiles(srcBinDir, outDir); err != nil {
		fmt.Println("❌ Fehler beim Kopieren der BIN-Dateien:", err)
		return
	}
	fmt.Println("✅ BIN-Dateien kopiert nach:", outDir)

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

	// Randomize
	if *difficulty < 0 || *difficulty > 2 {
		log.Fatal("❌ Invalid difficulty level. Must be 0 (Easy), 1 (Normal), 2 (Hard), or 3 (Extreme).")
	}
	fmt.Printf("🎲 Randomizing items with difficulty level %d...\n", *difficulty)
	if *difficulty == 0 {
		if err := uncompressed.RandomizeCards(outDir, outDir, finalSeed, logDir); err != nil {
			log.Fatal("❌ Randomization failed:", err)
		}
	} else if *difficulty == 1 {
		if err := uncompressed.RandomizeCards(outDir, outDir, finalSeed, logDir); err != nil {
			log.Fatal("❌ Randomization failed:", err)
		}
		if err := uncompressed.RandomizeLevelItems(outDir, outDir, finalSeed); err != nil {
			log.Fatal(err)
		}
	} else if *difficulty == 2 {
		if err := uncompressed.RandomizeCards(outDir, outDir, finalSeed, logDir); err != nil {
			log.Fatal("❌ Randomization failed:", err)
		}
		if err := uncompressed.RandomizeBuildingItems(outDir, outDir, finalSeed); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("🔗 patch intro...")
	if err := rom.PatchIntro(expandedRom, logPath); err != nil {
		log.Fatal("❌ Intro patching failed:", err)
	}

	// Repack RNCs into ROM expanded area
	fmt.Println("📀 Packing RNC Files...")
	rncpropack.RepackAll(outDir, logPath)
	fmt.Println("📀 Embedding RNCs into expanded ROM...")
	newStarts, err := rom.AppendRNCsCompact(outDir, expandedRom, 0x200500, logPath, true)
	if err != nil {
		log.Fatal("❌ Embedding RNCs failed:", err)
	}

	// Validate pointers
	fmt.Println("🔍 Validating pointers...")
	if err := rom.ValidatePointers(outDir, 0x400000, logPath); err != nil {
		log.Fatal("❌ Pointer validation failed:", err)
	}

	// Patch pointers in ROM
	fmt.Println("🔗 Patching pointers...")
	if err := rom.PatchPointers(expandedRom, rom.PointerList, newStarts, logPath); err != nil {
		log.Fatal("❌ Pointer patching failed:", err)
	}

	if *overworld {
		fmt.Println("🔗 random overworld items...")
		if err := rom.RandomizeOverworldItems(expandedRom, finalSeed, logPath); err != nil {
			log.Fatal("❌ Applying random overworld items failed:", err)
		}
	} else {
		fmt.Println("ℹ️  Skipping random overworld items.")
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

func copyBinFiles(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".bin" {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		fmt.Printf("→ %s\n", dstPath)
		return nil
	})
}
