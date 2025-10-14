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

	//initialize key-lock matrix: 0 is the key location accessible?, 1 is the lock accessible?, 2 has a key been placed in this location?, 3 has the key to this lock been placed?
	slots := [8][4]bool{
		{false, true, false, false},  //0 Ray Arnold
		{true, false, false, false},  //1 Dennis Nedry
		{false, true, false, false},  //2 Donald Gennaro
		{true, false, false, false},  //3 Dr. Wu
		{false, false, false, false}, //4 Ellie Sattler - we don't care about this door; 1 can be marked false
		{true, false, false, false},  //5 Alan Grant - we don't care about this door; 1 can be marked false
		{true, true, false, false},   //6 Ian Malcom
		{true, false, false, false},  //7 Robert Muldoon - we don't care about this door; 1 can be marked false
	}

	locpool := [8]int{}  //pool of available key locations
	doorpool := [8]int{} //pool if available doors
	loccount := 0        //number of available key locations
	doorcount := 0       //number of available doors
	chooseloc := 0       //chosen location to place ket
	choosedoor := 0      //chosen door to unlock

	for cardsleft := 0; cardsleft < 8; cardsleft++ {

		loccount = 0
		doorcount = 0

		for i := 0; i < 8; i++ {

			if slots[i][0] == true { //place available locations in the pool
				if slots[i][2] == false {
					locpool[loccount] = i
					loccount++
				}
			}

			if slots[i][1] == true { //place available doors in the pool
				if slots[i][3] == false {
					doorpool[doorcount] = i
					doorcount++
				}
			}
		}

		if doorcount == 0 { //if all important doors are unlocked, choose from the remaining doors
			for i := 0; i < 8; i++ {
				if slots[i][3] == false {
					doorpool[doorcount] = i
					doorcount++
				}
			}
		}

		chooseloc = locpool[r.Intn(loccount)] //from the pools, choose a door to unlock and the key placement
		choosedoor = doorpool[r.Intn(doorcount)]

		// place chosen id card into indices based on value
		if choosedoor == 0 {
			indices[chooseloc] = 10
		} // Ray Arnold 0A
		if choosedoor == 1 {
			indices[chooseloc] = 12
		} // Dennis Nedry 0C
		if choosedoor == 2 {
			indices[chooseloc] = 8
		} // Donald Gennaro 08
		if choosedoor == 3 {
			indices[chooseloc] = 14
		} // Dr Wu 14
		if choosedoor == 4 {
			indices[chooseloc] = 2
		} // Ellie Sattler 02
		if choosedoor == 5 {
			indices[chooseloc] = 6
		} // Alan Grant 06
		if choosedoor == 6 {
			indices[chooseloc] = 16
		} // Ian Malcom 16
		if choosedoor == 7 {
			indices[chooseloc] = 4
		} // Robert Muldoon 04

		slots[choosedoor][3] = true
		slots[chooseloc][2] = true

		// update the key-lock matrix
		if slots[2][3] == true {
			slots[0][0] = true
		} //if Gennaro's key has been placed, then Arnold's key location is accessible
		if slots[0][3] == true {
			slots[1][1] = true
		} //if Arnold's key has been placed, then Nedry's lock is accessible
		if slots[6][3] == true {
			slots[2][0] = true
		} //if Malcom's key has been placed, then Gennaro's key location is accessible
		if slots[1][3] == true {
			slots[3][1] = true
		} //if Nedry's key has been placed, then Wu's lock is accessible
		if slots[3][3] == true {
			slots[4][0] = true
		} //if Wu's key has been placed, then Sattler key location is accessible

	}

	// for i, loc := range CardLocations {
	//	indices[i] = loc.Index
	// }
	// r.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	// fmt.Println(indices)

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
