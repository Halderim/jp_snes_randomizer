package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Item struct {
	Offset   string `json:"offset"`
	Bytes    string `json:"bytes"`
	ItemID   string `json:"item_id"`
	SpriteID string `json:"sprite_id"`
	CardID   string `json:"card_id"`
	ItemName string `json:"item_name,omitempty"`
	Building string `json:"building"`
	DarkRoom bool   `json:"dark_room"`
}

var itemMap = map[byte]struct {
	Name       string
	Visibility byte
	SpriteLo   byte
	SpriteHi   byte
}{
	0x02: {"ID Card", 0x00, 0xE8, 0x00},
	0x04: {"Food", 0x00, 0xE9, 0x00},
	0x08: {"Medikit", 0x00, 0xC9, 0x00},
	0x3C: {"Gas Grenade Ammo", 0x00, 0xEA, 0x00},
	0x3E: {"Shotgun Ammo", 0x00, 0xEB, 0x00},
	0x40: {"Dart Ammo", 0x00, 0xEC, 0x00},
	0x42: {"Rocket Ammo", 0x00, 0xED, 0x00},
	0x44: {"Battery", 0x00, 0xEE, 0x00},
	0x46: {"Nerv Gas Canister", 0x00, 0xEF, 0x00},
	0x48: {"Bola Ammo", 0x01, 0x00, 0x01},
	0x54: {"Extra Life", 0x01, 0x1F, 0x01},
}

func main() {
	inputDir := flag.String("dir", "../../internal/uncompressed", "Ordner mit Level-Dateien")
	jsonOut := flag.String("json", "items.json", "JSON-Ausgabe")
	textOut := flag.String("txt", "items.txt", "Text-Ausgabe")
	flag.Parse()

	results := make(map[string][]Item)
	var textLines []string

	err := filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".bin" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var items []Item
		for i := 0; i+12 <= len(data); i++ {
			id := data[i+1]
			mapping, known := itemMap[id]
			if !known {
				continue
			}

			visibility := data[i+0]
			spriteLo := data[i+9]
			spriteHi := data[i+10]
			cardID := data[i+11]

			if visibility != mapping.Visibility || spriteLo != mapping.SpriteLo || spriteHi != mapping.SpriteHi {
				continue
			}

			valid := true
			for _, idx := range []int{2, 3, 4, 5, 6, 7, 8} {
				if data[i+idx] != 0x00 {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}

			block := data[i : i+12]
			bytesStr := fmt.Sprintf("% X", block)
			item := Item{
				Offset:   fmt.Sprintf("0x%08X", i),
				Bytes:    bytesStr,
				ItemID:   fmt.Sprintf("0x%02X", id),
				SpriteID: fmt.Sprintf("0x%02X%02X", spriteHi, spriteLo),
				CardID:   fmt.Sprintf("0x%02X", cardID),
				ItemName: mapping.Name,
				Building: strings.SplitN(filepath.Base(path), "_", 2)[0],
				DarkRoom: false,
			}
			items = append(items, item)

			textLine := fmt.Sprintf("{\"%s\", %s, 0x%02X, 0x%02X%02X, 0x%02X, \"%s\",\"%s\",\"%t\"},",
				filepath.Base(path),
				item.Offset,
				id,
				spriteHi, spriteLo,
				cardID,
				mapping.Name,
				strings.SplitN(filepath.Base(path), "_", 2)[0],
				false,
			)
			textLines = append(textLines, textLine)
		}

		if len(items) > 0 {
			results[filepath.Base(path)] = items
		}
		return nil
	})

	if err != nil {
		fmt.Println("Fehler:", err)
		return
	}

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = ioutil.WriteFile(*jsonOut, jsonData, 0644)

	output := strings.Join(textLines, "\n")
	_ = ioutil.WriteFile(*textOut, []byte(output), 0644)

	fmt.Println("✅ Fertig!")
	fmt.Println("JSON:", *jsonOut)
	fmt.Println("TXT :", *textOut)
}
