package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jp_snes_randomizer/internal/rom"
	"jp_snes_randomizer/internal/tools/rncpropack"
	"jp_snes_randomizer/internal/uncompressed"
)

type RandomizerSettings struct {
	RomPath        string `json:"romPath"`
	StartLocations bool   `json:"startLocations"`
	Seed           string `json:"seed"`
	Difficulty     int    `json:"difficulty"`
	Overworld      bool   `json:"overworld"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/randomize", handleRandomize)
	http.HandleFunc("/download", handleDownload)

	addr := ":" + port
	fmt.Printf("🦖 Jurassic Park Randomizer Web running at http://0.0.0.0%s\n", addr)
	http.ListenAndServe(addr, nil)
}

// --- Upload & ROM-Prüfung ---
func handleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("rom")
	if err != nil {
		http.Error(w, "Fehler beim Upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d.smc", time.Now().UnixNano()))
	out, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	valid, reason := checkRom(tempPath)
	if !valid {
		os.Remove(tempPath)
		http.Error(w, fmt.Sprintf("Ungültiges ROM – %s", reason), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": tempPath,
		"name": header.Filename,
	})
}

// --- ROM-Check ---
func checkRom(path string) (bool, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, "Datei konnte nicht gelesen werden"
	}
	if len(data) < 0x8000 {
		return false, "Datei ist zu klein – kein gültiges SNES-ROM"
	}

	// Prüfen ob 512-Byte-Header existiert
	hasHeader := (len(data) % 0x8000) == 512

	baseOffset := 0x7FC0
	if hasHeader {
		baseOffset += 512
	}

	if len(data) < baseOffset+0x40 {
		return false, "ROM-Datei ist beschädigt oder unvollständig"
	}

	header := data[baseOffset : baseOffset+0x40]
	gameTitle := string(bytes.Trim(header[0x00:0x15], "\x00 "))
	layout := header[0x15]
	country := header[0x19]
	version := header[0x1B]

	// Titelprüfung
	if !strings.Contains(strings.ToUpper(gameTitle), "JURASSIC") {
		return false, fmt.Sprintf("Falsches Spiel (%q gefunden)", strings.TrimSpace(gameTitle))
	}

	// Layoutprüfung (LoROM + Fastrom)
	if layout != 0x30 {
		return false, fmt.Sprintf("Falsches ROM-Layout: 0x%X (erwartet 0x30 = LoROM+FastROM)", layout)
	}

	// Ländercode
	if country != 0x01 {
		return false, fmt.Sprintf("Falscher Ländercode: 0x%X (erwartet 0x01 = USA)", country)
	}

	// Versionsprüfung
	if version != 0x00 {
		return false, fmt.Sprintf("Falsche Version: 0x%X (erwartet 0x00 = v1.0)", version)
	}

	return true, ""
}

// --- Randomizer ---
func handleRandomize(w http.ResponseWriter, r *http.Request) {
	var settings RandomizerSettings
	err := json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		http.Error(w, "Ungültige Einstellungen", http.StatusBadRequest)
		return
	}

	// Validierung
	if settings.RomPath == "" {
		http.Error(w, "ROM-Pfad fehlt", http.StatusBadRequest)
		return
	}

	if settings.Difficulty < 0 || settings.Difficulty > 2 {
		http.Error(w, "Ungültiger Schwierigkeitsgrad (0-2)", http.StatusBadRequest)
		return
	}

	// Seed bestimmen
	var finalSeed int64
	if settings.Seed != "" {
		seedInt, err := strconv.ParseInt(settings.Seed, 10, 64)
		if err != nil {
			http.Error(w, "Ungültiger Seed", http.StatusBadRequest)
			return
		}
		finalSeed = seedInt
	} else {
		finalSeed = time.Now().Unix()
	}

	// Arbeitsverzeichnisse vorbereiten
	srcBinDir := "internal/uncompressed"
	outDir := filepath.Join(os.TempDir(), fmt.Sprintf("randomizer_%d_%d", finalSeed, time.Now().UnixNano()))
	logDir := filepath.Join(outDir, "logs")
	logPath := filepath.Join(logDir, "randomizer.log")

	if err := os.MkdirAll(outDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Erstellen des Arbeitsverzeichnisses: %v", err), http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Erstellen des Log-Verzeichnisses: %v", err), http.StatusInternalServerError)
		return
	}

	// BIN-Dateien kopieren
	if err := copyBinFiles(srcBinDir, outDir); err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Kopieren der BIN-Dateien: %v", err), http.StatusInternalServerError)
		return
	}

	// ROM kopieren
	expandedRom := filepath.Join(outDir, fmt.Sprintf("jp_randomized_seed%d.sfc", finalSeed))
	if err := copyFile(settings.RomPath, expandedRom); err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Kopieren der ROM: %v", err), http.StatusInternalServerError)
		return
	}

	// Randomisierung durchführen
	if err := runRandomization(outDir, expandedRom, finalSeed, settings.Difficulty, settings.StartLocations, settings.Overworld, logPath); err != nil {
		http.Error(w, fmt.Sprintf("Fehler bei der Randomisierung: %v", err), http.StatusInternalServerError)
		return
	}

	// Antwort senden
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"download": "/download?path=" + expandedRom,
		"seed":     fmt.Sprintf("%d", finalSeed),
	})
}

// --- Randomisierungs-Logik ---
func runRandomization(outDir, expandedRom string, seed int64, difficulty int, start, overworld bool, logPath string) error {
	// Randomisierung basierend auf Schwierigkeitsgrad
	if difficulty == 0 {
		if err := uncompressed.RandomizeCards(outDir, outDir, seed, filepath.Dir(logPath)); err != nil {
			return fmt.Errorf("Randomisierung fehlgeschlagen: %w", err)
		}
	} else if difficulty == 1 {
		if err := uncompressed.RandomizeCards(outDir, outDir, seed, filepath.Dir(logPath)); err != nil {
			return fmt.Errorf("Randomisierung fehlgeschlagen: %w", err)
		}
		if err := uncompressed.RandomizeLevelItems(outDir, outDir, seed); err != nil {
			return fmt.Errorf("Randomisierung fehlgeschlagen: %w", err)
		}
	} else if difficulty == 2 {
		if err := uncompressed.RandomizeCards(outDir, outDir, seed, filepath.Dir(logPath)); err != nil {
			return fmt.Errorf("Randomisierung fehlgeschlagen: %w", err)
		}
		if err := uncompressed.RandomizeBuildingItems(outDir, outDir, seed); err != nil {
			return fmt.Errorf("Randomisierung fehlgeschlagen: %w", err)
		}
	}

	// Intro patchen
	if err := rom.PatchIntro(expandedRom, logPath); err != nil {
		fmt.Printf("Fehler beim Patchen des Intros: %v \n", err)
		return fmt.Errorf("Intro-Patching fehlgeschlagen: %w", err)
	}

	// RNC-Dateien packen und einbetten
	rncpropack.RepackAll(outDir, logPath)
	newStarts, err := rom.AppendRNCsCompact(outDir, expandedRom, 0x200500, logPath, true)
	if err != nil {
		return fmt.Errorf("RNC-Einbettung fehlgeschlagen: %w", err)
	}

	// Pointer validieren
	if err := rom.ValidatePointers(outDir, 0x400000, logPath); err != nil {
		return fmt.Errorf("Pointer-Validierung fehlgeschlagen: %w", err)
	}

	// Pointer patchen
	if err := rom.PatchPointers(expandedRom, rom.PointerList, newStarts, logPath); err != nil {
		return fmt.Errorf("Pointer-Patching fehlgeschlagen: %w", err)
	}

	// Overworld randomisieren
	if overworld {
		if err := rom.RandomizeOverworldItems(expandedRom, seed, logPath); err != nil {
			return fmt.Errorf("Overworld-Randomisierung fehlgeschlagen: %w", err)
		}
	}

	// Start-Location randomisieren
	if start {
		if err := rom.ApplyRandomStartLocation(expandedRom, seed, logPath); err != nil {
			return fmt.Errorf("Start-Location-Randomisierung fehlgeschlagen: %w", err)
		}
	}

	// QoL-Patches anwenden
	if err := rom.ApplyQolPatches(expandedRom, seed, logPath); err != nil {
		return fmt.Errorf("QoL-Patches fehlgeschlagen: %w", err)
	}

	return nil
}

// --- Download-Handler ---
func handleDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Pfad fehlt", http.StatusBadRequest)
		return
	}

	// Sicherheitsprüfung: Nur Dateien aus dem Temp-Verzeichnis erlauben
	if !strings.HasPrefix(path, os.TempDir()) {
		http.Error(w, "Ungültiger Pfad", http.StatusForbidden)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Datei nicht gefunden", http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Fehler beim Lesen der Datei", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	io.Copy(w, file)
}

// --- Helper-Funktionen ---
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

		return nil
	})
}
