package uncompressed

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RandomizeLevelItems swapping all items in on a floor of a building
func RandomizeLevelItems(binDir, outDir string, seed int64) error {
	r := rand.New(rand.NewSource(seed))
	const itemBlockSize = 13

	var logBuilder strings.Builder
	logBuilder.WriteString(fmt.Sprintf("🎲 Item-Level Randomizer (Seed: %d)\n", seed))
	logBuilder.WriteString(strings.Repeat("-", 40) + "\n")

	levelMap := make(map[string][]ItemEntry)
	for _, item := range GameItems {
		//fmt.Println(strings.SplitN(item.File, "_", 2)[0])
		levelMap[item.File] = append(levelMap[item.File], item)
	}

	for levelName, items := range levelMap {
		if len(items) < 3 {
			continue
		}

		input := filepath.Join(binDir, levelName)
		output := filepath.Join(outDir, levelName)

		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("Error reading %s: %w", levelName, err)
		}

		permutation, ok := FindValidPermutation(items, r, 1000)
		if !ok {
			logBuilder.WriteString(fmt.Sprintf("⚠️ No permutation found %s \n", levelName))
			continue
		}

		for i, targetIdx := range permutation {
			if i == targetIdx {
				continue
			}

			src := items[i].Offset
			dst := items[targetIdx].Offset

			if src+itemBlockSize > len(data) || dst+itemBlockSize > len(data) {
				continue
			}

			tmp := make([]byte, itemBlockSize)
			copy(tmp, data[src:src+itemBlockSize])
			copy(data[src:src+itemBlockSize], data[dst:dst+itemBlockSize])
			copy(data[dst:dst+itemBlockSize], tmp)

			logBuilder.WriteString(fmt.Sprintf(
				"%s: Swap @%04X ↔ @%04X (%s ↔ %s)\n",
				levelName, src, dst, items[i].ItemName, items[targetIdx].ItemName))
		}

		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(output, data, 0644); err != nil {
			return err
		}
	}

	logPath := filepath.Join(outDir, fmt.Sprintf("randomizer_items_seed%d.txt", seed))
	if err := os.WriteFile(logPath, []byte(logBuilder.String()), 0644); err != nil {
		return err
	}

	fmt.Println("✅ Items in level swapped:", logPath)
	return nil
}

// Function to randomize items in a building and not only on a floor
func RandomizeBuildingItems(binDir, outDir string, seed int64) error {
	r := rand.New(rand.NewSource(seed))
	const itemBlockSize = 13

	var logBuilder strings.Builder
	logBuilder.WriteString(fmt.Sprintf("🏢 Items Building Randomizer (Seed: %d | %s)\n", seed, time.Now().Format("2006-01-02 15:04:05")))
	logBuilder.WriteString(strings.Repeat("-", 60) + "\n")

	buildingMap := make(map[string][]ItemEntry)
	for _, item := range GameItems {
		buildingMap[item.Building] = append(buildingMap[item.Building], item)
	}

	for building, items := range buildingMap {
		if len(items) < 2 {
			continue
		}

		logBuilder.WriteString(fmt.Sprintf("\n--- 🏠 %s (%d Items) ---\n", building, len(items)))

		permutation, ok := FindValidPermutation(items, r, 1000)
		if !ok {
			logBuilder.WriteString(fmt.Sprintf("⚠️ No permutation found %s \n", building))
			continue
		}

		for i, dstIdx := range permutation {
			if i == dstIdx {
				continue
			}

			srcItem := items[i]
			dstItem := items[dstIdx]

			inputSrc := filepath.Join(binDir, srcItem.File)
			inputDst := filepath.Join(binDir, dstItem.File)

			outputSrc := filepath.Join(outDir, srcItem.File)
			outputDst := filepath.Join(outDir, dstItem.File)

			dataSrc, err := os.ReadFile(inputSrc)
			if err != nil {
				return fmt.Errorf("Error reading %s: %w", inputSrc, err)
			}
			dataDst, err := os.ReadFile(inputDst)
			if err != nil {
				return fmt.Errorf("Error reading %s: %w", inputDst, err)
			}

			if srcItem.Offset+itemBlockSize > len(dataSrc) || dstItem.Offset+itemBlockSize > len(dataDst) {
				logBuilder.WriteString(fmt.Sprintf("⚠️  Offset out of bounds %s (%s @%X ↔ %s @%X)\n",
					building, srcItem.File, srcItem.Offset, dstItem.File, dstItem.Offset))
				continue
			}

			tmp := make([]byte, itemBlockSize)
			copy(tmp, dataSrc[srcItem.Offset:srcItem.Offset+itemBlockSize])
			copy(dataSrc[srcItem.Offset:srcItem.Offset+itemBlockSize],
				dataDst[dstItem.Offset:dstItem.Offset+itemBlockSize])
			copy(dataDst[dstItem.Offset:dstItem.Offset+itemBlockSize], tmp)

			if err := os.MkdirAll(filepath.Dir(outputSrc), 0755); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(outputDst), 0755); err != nil {
				return err
			}

			if err := os.WriteFile(outputSrc, dataSrc, 0644); err != nil {
				return err
			}
			if err := os.WriteFile(outputDst, dataDst, 0644); err != nil {
				return err
			}

			logBuilder.WriteString(fmt.Sprintf("%s: %s @%04X ↔ %s @%04X (%s ↔ %s)\n",
				building, srcItem.File, srcItem.Offset, dstItem.File, dstItem.Offset,
				srcItem.ItemName, dstItem.ItemName))
		}
	}

	logPath := filepath.Join(outDir, fmt.Sprintf("randomizer_buildings_seed%d.txt", seed))
	if err := os.WriteFile(logPath, []byte(logBuilder.String()), 0644); err != nil {
		return err
	}

	fmt.Println("✅ Building items randomized:", logPath)
	return nil
}
