
package main

import (
	"os"
)

const (
	// Default value.
	MinRabbits	= 1
	// Default value. The numbers of rabbits that exist
	// at any given time.
	MaxRabbits	= 15
	// Spawn chance for rabbits.
	SpawnChance	= 0.20
	
	// Chance to ascend deeper (closer to /). The weight
	// has to be fair, because ascending is very limited
	// and you only have one option.
	AscendChance	= 0.30
	
	// Chance to move twice instead of once.
	TwoStepChance	= 0.50
)

type directoryForest struct {
	// List of rabbits and their locations. Only one
	// rabbit per location.
	rabbits		map[string]Rabbit	`json:"rabbits"`
	// Number of rabbits seen.
	spottedCount	uint			`json:"spottedCount"`
	// Number of rabbits caught.
	caughtCount	uint			`json:"caughtCount"`
	// Number of rabbits killed. :(
	killedCount	uint			`json:"killedCount"`
}

func newDirectoryForest() directoryForest {
	return directoryForest{
		map[string]Rabbit{}, 0, 0, 0,
	}
}

// For a location to exist, the directory must exist.
func (f *directoryForest) LocationExists(loc string) bool {
	fi, err := os.Stat(loc)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// Returns a location near the passed location. Nearby is
// found by a small number of random directory changes. Will
// not be the same directory, unless it has to (can't ascend
// or descend).
//
// XXX: Currently this does not check if a rabbit already
// exists at the new location. So we just lose rabbits
// if one encounters another.
func (f *directoryForest) NearbyLocation(loc string) string {
	newloc := loc

	steps := 1
	if chance(TwoStepChance) {
		steps = 2
	}

tryagain:
	for i := 0; i < steps; i++ {
		// Can't move.
		if !canAscend(newloc) && !canDescend(newloc) {
			return newloc
		} else if chance(AscendChance) {
			if canAscend(newloc) {
				newloc = ascend(newloc)
			} else {
				newloc = randDescension(newloc)
			}
		} else {
			if canDescend(newloc) {
				newloc = randDescension(newloc)
			} else {
				newloc = ascend(newloc)
			}
		}
	}
	
	if newloc == loc {
		// Guaranteed to not be the same because you must
		// step twice to get to the same destination.
		steps = 1
		goto tryagain
	}
	
	return newloc
	
}

// A random faraway location. Rabbits typically start here
// and run here when they're fleeing.
func (f *directoryForest) FarawayLocation(loc string) string {
	newloc := baseLocation()
	triedagain := false
	
	steps := 1
	if chance(TwoStepChance) {
		steps = 2
	}

tryagain:	
	for i := 0; i < steps; i++ {
		if canDescend(newloc) {
			newloc = randDescension(newloc)
		}
	}
	
	if newloc == loc && !triedagain {
		// Invert the steps. We can't get to the same
		// location with different steps.
		if steps == 1 {
			steps = 2
		} else if steps == 2 {
			steps = 1
		}
		triedagain = true
		goto tryagain
	}
	
	return newloc
}

// Anytime a location is entered, a check is performed. This
// function updates every rabbit and returns a rabbit if one
// is spotted.
func (f *directoryForest) CheckLocation(loc string) (spotted *Rabbit) {
	spotted = nil
	
	newrabbits := map[string]Rabbit{}

	for _, r := range f.rabbits {
		r.DisturbanceAt(loc)
		
		// Can't move usually means caught or dead.
		if (!r.CantMove()) {
			if r.JustSpotted() {
				spotted = &r
				f.spottedCount++
			}
			newrabbits[r.Location()] = r
		} else {
			if r.State() == Dead {
				f.killedCount++
			} else if r.State() == Caught {
				f.caughtCount++
			}
		}
	}
	
	f.rabbits = newrabbits
	
	// See if we should repopulate.
	f.populate()
	
	return
}

func (f *directoryForest) populate() {
	for len(f.rabbits) < MinRabbits {
		r := NewRabbit(f)
		f.rabbits[r.Location()] = r
	}
	
	if chance(SpawnChance) {
		r := NewRabbit(f)
		f.rabbits[r.Location()] = r
	}
}