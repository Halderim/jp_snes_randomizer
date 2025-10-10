package uncompressed

type CardLocation struct {
	File   string
	Index  byte
	Offset int // absoluter Offset im BIN-File, wo das Kartenindex-Byte sitzt
}

var CardLocations = []CardLocation{
	{"BeachUtilityShed_SubLevel_Entities.bin", 0x0A, 0x36E},    // Ray Arnold
	{"BeachUtilityShed_GroundLevel_Entities.bin", 0x0C, 0x2BC}, // Dennis Nedry
	{"NublarUtilityShed_SubLevel_Entities.bin", 0x08, 0x706},   // Donald Gennaro
	{"Ship_SubLevel1_Entities.bin", 0x0E, 0x922},               // Dr. Wu
	{"Ship_SubLevel3_Entities_JP_U1.bin", 0x02, 0x1F2A},        // Ellie Sattler
	{"VisitorCenter_GroundLevel_Entities.bin", 0x06, 0xA9E},    // Alan Grant
	{"RaptorPen_UpperLevel_Entities_JP_U1.bin", 0x10, 0x670},   // Ian Malcolm
	{"RaptorPen_SubLevel1_Entities.bin", 0x04, 0x2AF4},         // Robert Muldoon
}
