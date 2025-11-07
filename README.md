# Jurassic Park SNES Randomizer

## Randomizer for Jurassic Park SNES 
A randomizer for the Game Jurassic Park on the SNES written in go

## How to run
You need to own a 1.0 rom of Jurassic Park USA (Jurassic Park Classic Game Collection does work)
I did not test Version 1.1 or different language versions

The rom needs to be extended to 4 MB for example with Lunar Expand \
The decompiled and unpacked binary files for the interior maps, the disassambly provided by Yoshifanatic extracts the rnc compressed data which can be unpacked by the rnc propack 
https://github.com/Yoshifanatic1/Jurassic-Park-1-SNES-disassembly \
https://github.com/lab313ru/rnc_propack_source/releases

### run default randomizer
go run ./cmd/randomizer [--seed] [--start] [--difficulty]

### options
--seed Int64 sets the seed that will be used to randomize the game for example --seed 12345 \
--start boolean if the option is true a randomized location will be used to start the game default: false\
--difficulty 0,1,2 switches between the different randomizer levels of difficulty, 0 only changes id cards, 1 changes items per floor, 2 changes items per building default: 0

## Future plans if they can be realized 
- adding more locations to start the game
- ~~randomize locations of ID cards inside buildings not just swap the cards in their default location~~
- ~~randomize ammo types in buildings~~
- ~~randomize health items in buildings~~
- randomize items in the overworld
- randomize egg locations on the overworld
- make eggs appear in buildings
- randomize computer terminals
- randomize building floors
- difficulty settings
- web frontend for the randomizer

### Version 0.2 alpha
The randomizer now can swap items inside a floor of a building or across a building when the difficulty is set to 1 or 2 
There is no logic so far beside preventing from batteries spawning inside a dark room 
Added:
- Randomize items per floor
- randomize items per building


### Version 0.1.1 alpha
QoL patches added to give infinite lives add the save feature from Yoshifanatic\
Thanks to coconutED for writing a basic logic that makes most if not all seeds beatable 
Added:
- QoL patches 
- More starting locations
- logic to make seeds beatable, preventing cards to be placed at unreachable locations

### Version 0.1 alpha
Initial version of the rando with no logic, you may be softlocked as a ID card may not be reachable\
Start locations can be randomized 4 new locations + default location the locations are currently placed in a way that you are not softlocked by the start location\
Features:
- Randomizing ID cards
- Randomizing start locations
