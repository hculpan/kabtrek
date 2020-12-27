# KabTrek
This is a rewrite of the classic text-based Star Trek game.  It is written in go using the TCell (https://github.com/gdamore/tcell) library for terminal output.

I have tried to stick with the spirit of the original game, and the rules generally follow the same as for most of the various versions.

# The Game
You are the commander of the starship U.S.S. Enterprise.  Federation space has been invaded by Klingons, and since yours is the only ship in the quadrant,
it is up to you to turn back the tide.  The galaxy is split up into 64 quadrants (8x8), and the Klingons are spread throughout.  Also within this space
are 5 starbases.  Destroy the Klingons before they destroy the starbases and don't let your ship be destroyed, and you win the war!

# Building
This is written in Go v1.15.6, and uses v2 of the TCell library.

To build, you should be able to copy the main branch locally and use `go build` or `go install`.
