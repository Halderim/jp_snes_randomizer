----------------------------------------------------
-- SNES Jurassic Park – Koordinaten-Navigation
-- Numpad 7 = zurück
-- Numpad 9 = weiter
----------------------------------------------------
prev7 = false
prev9 = false

X_LO = 0x7E2476   -- Beispielwerte!
X_HI = 0x7E2477
Y_LO = 0x7E2C4A
Y_HI = 0x7E2C4A
Z_LO = 0x7E72BE

positions = {
	{ x = 0x0204, y = 0x0590, z = 0x0800 }, --Shotgun ammo
	{ x = 0x00E4, y = 0x07A0, z = 0x0800 },
	{ x = 0x0C14, y = 0x0C90, z = 0x0800 },
	{ x = 0x0A74, y = 0x0D90, z = 0x0800 },
	{ x = 0x0494, y = 0x04A0, z = 0x0800 },
	{ x = 0x0AF4, y = 0x0500, z = 0x0800 },
	{ x = 0x0E44, y = 0x0250, z = 0x0800 },
	{ x = 0x03C4, y = 0x09F0, z = 0x0800 },
	{ x = 0x0294, y = 0x0F10, z = 0x0800 },
	{ x = 0x0594, y = 0x0D60, z = 0x5000 },
	{ x = 0x07C4, y = 0x0420, z = 0x4000 }, -- GAS GAS GAS
	{ x = 0x0104, y = 0x0AC0, z = 0x0800 },
	{ x = 0x06A4, y = 0x0B60, z = 0x0800 },
	{ x = 0x0924, y = 0x0750, z = 0x0800 },
	{ x = 0x0DD4, y = 0x08D0, z = 0x4000 },
	{ x = 0x0044, y = 0x04E0, z = 0x0800 },
	{ x = 0x0B34, y = 0x00A0, z = 0x0800 },
	{ x = 0x0B34, y = 0x0280, z = 0x0800 },
	{ x = 0x0DE4, y = 0x0700, z = 0x5000 },
	{ x = 0x0414, y = 0x0CA0, z = 0x0800 },
	{ x = 0x00A4, y = 0x0150, z = 0x0800 }, -- Darts
	{ x = 0x0234, y = 0x06C0, z = 0x0800 },
	{ x = 0x00D4, y = 0x0E10, z = 0x0800 },
	{ x = 0x0614, y = 0x09E0, z = 0x0800 },
	{ x = 0x09A4, y = 0x0F00, z = 0x0800 },
	{ x = 0x0CC4, y = 0x0DE0, z = 0x0800 },
	{ x = 0x0684, y = 0x0120, z = 0x0800 },
	{ x = 0x0884, y = 0x0070, z = 0x0800 },
	{ x = 0x0DD4, y = 0x0800, z = 0x4000 },
	{ x = 0x0E24, y = 0x09C0, z = 0x0800 },
	{ x = 0x0274, y = 0x0770, z = 0x0800 },
	{ x = 0x0E34, y = 0x0310, z = 0x0800 }, -- BOLA
	{ x = 0x07A4, y = 0x0650, z = 0x4000 },
	{ x = 0x0314, y = 0x0750, z = 0x0800 },
	{ x = 0x0BF4, y = 0x06B0, z = 0x0800 },
	{ x = 0x0494, y = 0x0090, z = 0x0800 },
	{ x = 0x04E4, y = 0x04F0, z = 0x0800 },
	{ x = 0x0564, y = 0x0860, z = 0x0800 },
	{ x = 0x0724, y = 0x0940, z = 0x0800 },
	{ x = 0x0664, y = 0x0F70, z = 0x0800 },
	{ x = 0x0214, y = 0x0300, z = 0x0800 }, -- ROCKET
	{ x = 0x0794, y = 0x0120, z = 0x5000 },
	{ x = 0x0F74, y = 0x0350, z = 0x0800 },
	{ x = 0x0544, y = 0x0A60, z = 0x0800 },
	{ x = 0x0934, y = 0x0AE0, z = 0x0800 },
	{ x = 0x0F94, y = 0x0AE0, z = 0x4000 },
	{ x = 0x0714, y = 0x02F0, z = 0x0800 },
	{ x = 0x0A54, y = 0x02D0, z = 0x4000 },
	{ x = 0x0564, y = 0x0610, z = 0x4000 },
	{ x = 0x0A84, y = 0x0B90, z = 0x4000 },
	{ x = 0x0934, y = 0x0510, z = 0x0800 }, -- 1UP
	{ x = 0x0FD4, y = 0x0720, z = 0x0800 },
	{ x = 0x0F74, y = 0x0F50, z = 0x0800 },
	{ x = 0x00B4, y = 0x0B40, z = 0x0800 },
	{ x = 0x00C4, y = 0x02E0, z = 0x0800 }, -- Medikit
	{ x = 0x07A4, y = 0x0B30, z = 0x0800 },
	{ x = 0x0524, y = 0x0420, z = 0x4000 }, -- Hammond
}

index = 0
	


local function applyPosition()
    local pos = positions[index]
    emu.write16(X_LO, pos.x, "snesWorkRam")
    emu.write16(Y_LO, pos.y, "snesWorkRam")
    emu.write16(Z_LO, pos.z, "snesWorkRam")
end

function teleport()
	if emu.isKeyPressed("Numpad 7") and not prev7 then
		index = index - 1
		if index < 1 then index = #positions end
        applyPosition()
        prev7 = true
        
    end
    
    if emu.isKeyPressed("Numpad 9") and not prev9 then
		index = index + 1
		if index > #positions then index = 1 end
        applyPosition()
        
        prev9 = true
    end
    
    if emu.isKeyPressed("Numpad 0") then
    	prev7 = false
    	prev9 = false
	end
	emu.drawString(12, 20, "Index: " .. string.format("%d", index), 0xFFFFFF, 0xFF000000)
end

emu.addEventCallback(teleport, emu.eventType.endFrame);