package uncompressed

import (
	"math/rand"
	"strings"
)

// FindValidPermutation searches for a random permutation that is valid based on some rules,
// (Battery swap into a dark room).

func FindValidPermutation(items []ItemEntry, r *rand.Rand, maxAttempts int) ([]int, bool) {
	isBattery := func(it ItemEntry) bool {
		return strings.Contains(strings.ToLower(it.ItemName), "battery")
	}

	isDark := func(it ItemEntry) bool {
		return it.DarkRoom
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		perm := r.Perm(len(items))
		valid := true
		for i, dstIdx := range perm {
			src := items[i]
			dst := items[dstIdx]
			// Rule Battery not in dark rooms
			if (isBattery(src) && isDark(dst)) || (isBattery(dst) && isDark(src)) {
				valid = false
				break
			}
		}
		if valid {
			return perm, true
		}
	}
	return nil, false
}
