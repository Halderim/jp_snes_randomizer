----------------------------------------------------
-- SNES Jurassic Park – Player Pixel Mover (Mesen)
-- Numpad 4/6 = X-- / X++
-- Numpad 8/2 = Y-- / Y++
----------------------------------------------------
function moveAround()
	-- HIER die korrekten Player-Koordinaten eintragen!
	local X_LO = 0x7E2476   -- Beispielwerte!
	local X_HI = 0x7E2477
	local Y_LO = 0x7E2C4A
	local Y_HI = 0x7E2C4A
	
	
    -- aktuelle Position lesen
    local x = emu.read16(X_LO, "snesWorkRam", false)
    local y = emu.read16(Y_LO, "snesWorkRam", false)

    -- NUMPAD 4 = links
    if emu.isKeyPressed("Numpad 4") then
        x = (x - 1) & 0xFFFF
    end
    
    -- NUMPAD 6 = rechts
    if emu.isKeyPressed("Numpad 6") then
        x = (x + 1) & 0xFFFF
    end

    -- NUMPAD 8 = hoch
    if emu.isKeyPressed("Numpad 8") then
        y = (y - 1) & 0xFFFF
    end

    -- NUMPAD 2 = runter
    if emu.isKeyPressed("Numpad 2") then
        y = (y + 1) & 0xFFFF
    end

    -- neue Position schreiben
    emu.write16(X_LO,x,"snesWorkRam")
    emu.write16(Y_LO,y,"snesWorkRam")
end

emu.addEventCallback(moveAround, emu.eventType.endFrame);
	
