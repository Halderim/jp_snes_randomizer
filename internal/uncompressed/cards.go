package uncompressed

import "fmt"

// CardNames mappt das Index-Byte (Byte 12 im Item-Block) auf den Kartennamen.
var CardNames = map[byte]string{
	0x02: "Ellie Sattler",
	0x04: "Robert Muldoon",
	0x06: "Alan Grant",
	0x08: "Donald Gennaro",
	0x0A: "Ray Arnold",
	0x0C: "Dennis Nedry",
	0x0E: "Dr. Wu",
	0x10: "Ian Malcolm",

	// Optional: Hammond (wird derzeit nicht mitgeshufflet, aber als Referenz)
	0x53: "John Hammond",
}

// GetCardName liefert den Namen zur Index-Byte oder einen lesbaren Fallback,
// falls der Index nicht in der Map vorhanden ist.
func GetCardName(idx byte) string {
	if n, ok := CardNames[idx]; ok {
		return n
	}
	return fmt.Sprintf("Unknown (0x%02X)", idx)
}
