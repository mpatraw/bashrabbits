
package main

import (
	"os"
)

const (
	// Default value.
	MaxRabbits	= 30
	
	// Chance to ascend deeper (closer to /). The weight
	// has to be fair, because ascending is very limited
	// and you only have one option.
	AscendChance	= 0.35
	
	// Chance to move twice instead of once.
	TwoStepChance	= 0.50
)

type directoryForest struct {
	// List of rabbits and their locations. Only one
	// rabbit per location.
	rabbits		map[string]Rabbit
	// Number of rabbits currently in the forest.
	numRabbits	uint
	// Max number of rabbits allowed in the forest. Configurable.
	maxRabbits	uint
}

func newDirectoryForest() directoryForest {
	return directoryForest{
		make(map[string]Rabbit), 0, MaxRabbits,
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
	steps := 1
	if chance(TwoStepChance) {
		steps = 2
	}
	
	for i := 0; i < steps; i++ {
		// Can't move.
		if !canAscend(loc) && !canDescend(loc) {
			return loc
		} else if chance(AscendChance) && canAscend(loc) {
			loc = ascend(loc)
		} else {
			loc = randDescension(loc)
		}
	}
	
	return loc
	
}

// A random faraway location. Rabbits typically start here
// and run here when they're fleeing.
func (f *directoryForest) FarawayLocation() string {
	loc := baseLocation()
	
	steps := 1
	if chance(TwoStepChance) {
		steps = 2
	}
	
	for i := 0; i < steps; i++ {
		if canDescend(loc) {
			loc = randDescension(loc)
		}
	}
	
	return loc
}