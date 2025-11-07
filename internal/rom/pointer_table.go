package rom

type PointerEntry struct {
	Filename string
	Lo       int
	Hi       int
	Bank     int
}

// PointerList contains all pointer entries to be patched in the ROM.
var PointerList = []PointerEntry{
	{"BeachUtilityShed_GroundLevel_Entities.rnc", 0xA583C2, 0xA583C3, 0xA583C4},
	{"BeachUtilityShed_SubLevel_Entities.rnc", 0xA58389, 0xA5838A, 0xA5838B},
	{"NorthUtilityShed_GroundLevel_Entities.rnc", 0xA583F1, 0xA583F2, 0xA583F3},
	{"NorthUtilityShed_SubLevel_Entities_JP_U1.rnc", 0xA58425, 0xA58426, 0xA58427},
	{"NublarUtilityShed_GroundLevel_Entities.rnc", 0xA58341, 0xA58342, 0xA58343},
	{"NublarUtilityShed_SubLevel_Entities.rnc", 0xA58308, 0xA58309, 0xA5830A},
	{"RaptorNest_Entities.rnc", 0xA5805F, 0xA58060, 0xA58061},
	{"RaptorPen_EntryLevel_Entities.rnc", 0xA58454, 0xA58455, 0xA58456},
	{"RaptorPen_UpperLevel_Entities_JP_U1.rnc", 0xA5847E, 0xA5847F, 0xA58480},
	{"RaptorPen_GroundLevel_Entities.rnc", 0xA584BC, 0xA584BD, 0xA584BE},
	{"RaptorPen_SubLevel1_Entities.rnc", 0xA5850F, 0xA58510, 0xA58511},
	{"RaptorPen_SubLevel2_Entities.rnc", 0xA5855D, 0xA5855E, 0xA5855F},
	{"SecretLevel_Entities.rnc", 0xA58037, 0xA58038, 0xA58039},
	{"Ship_EntryLevel_Entities.rnc", 0xA58073, 0xA58074, 0xA58075},
	{"Ship_SubLevel1_Entities.rnc", 0xA5809D, 0xA5809E, 0xA5809F},
	{"Ship_SubLevel2_Entities_JP_U1.rnc", 0xA580D1, 0xA580D2, 0xA580D3},
	{"Ship_SubLevel3_Entities_JP_U1.rnc", 0xA5813D, 0xA5813E, 0xA5813F},
	{"Ship_SubLevel4_Entities.rnc", 0xA58176, 0xA58177, 0xA58178},
	{"VisitorCenter_GroundLevel_Entities.rnc", 0xA581C8, 0xA581C9, 0xA581CA},
	{"VisitorCenter_Floor1_Entities.rnc", 0xA58234, 0xA58235, 0xA58236},
	{"VisitorCenter_RoofLevel_Entities.rnc", 0xA58286, 0xA58287, 0xA58288},
	{"VisitorCenter_SubLevel_Entities.rnc", 0xA582B0, 0xA582B1, 0xA582B2},
}
