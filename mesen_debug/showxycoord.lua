--This is an example Lua (https://www.lua.org) script to give a general idea of how to build scripts
--Press F5 or click the Run button to execute it
--Type "emu." to show a list of all available API function

function printCoord()
 
  bgColor = 0x302060FF
  fgColor = 0x30FF4040

  
  --Draw some rectangles and print some text
  emu.drawRectangle(8, 16, 64, 40, bgColor, true, 1)
  emu.drawRectangle(8, 16, 64, 40, fgColor, false, 1)
  xposlabel = emu.getLabelAddress("Exterior_Player_XPos")
  yposlabel = emu.getLabelAddress("Exterior_Player_YPos")
  zposlabel = emu.getLabelAddress("Exterior_Player_ZPos")

  xpos = emu.read16(xposlabel["address"],xposlabel["memType"])
  ypos = emu.read16(yposlabel["address"],yposlabel["memType"])
  zpos = emu.read16(zposlabel["address"],zposlabel["memType"])
  emu.drawString(12, 20, "Xpos: " .. string.format("%x", xpos), 0xFFFFFF, 0xFF000000)
  emu.drawString(12, 30, "Ypos: " .. string.format("%x", ypos), 0xFFFFFF, 0xFF000000)
  emu.drawString(12, 40, "Zpos: " .. zpos, 0xFFFFFF, 0xFF000000)
end

--Register some code (printInfo function) that will be run at the end of each frame
emu.addEventCallback(printCoord, emu.eventType.endFrame);

--Display a startup message
emu.displayMessage("Script", "X Y Coordinates loaded")