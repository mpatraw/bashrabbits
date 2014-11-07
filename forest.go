
package main

import (
	"encoding/json"
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
	rabbits		map[string]*Rabbit
	// Number of rabbits seen.
	spottedCount	uint
	// Number of rabbits caught.
	caughtCount	uint
	// Number of rabbits killed. :(
	killedCount	uint
}

func newDirectoryForest() directoryForest {
	return directoryForest{
		map[string]*Rabbit{}, 0, 0, 0,
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
func (f *directoryForest) PerformCheck() (spotted *Rabbit) {
	spotted = nil

	// We always check our current directory.
	loc, _ := os.Getwd()

	newrabbits := map[string]*Rabbit{}

	for _, r := range f.rabbits {
		r.DisturbanceAt(loc)

		if (r.CanMove()) {
			if r.JustSpotted() {
				if spotted != nil {
					// XXX: Shouldn't happen.
					//panic("Spotted two rabbits. Impossible.")
				}
				spotted = r
				f.spottedCount++
			}
			newrabbits[r.Location()] = r
		} else {
			if r.State() == Dead {
				f.killedCount++
			} else if r.State() == Caught {
				// Being caught is updated in the TryCatch.
			}
		}
	}

	f.rabbits = newrabbits

	// See if we should repopulate.
	f.repopulate()

	return
}

// Attempts to catch a rabbit if it's still where we are.
func (f *directoryForest) PerformCatch() bool {
	loc, _ := os.Getwd()

	rab, ok := f.rabbits[loc]

	if ok {
		delete(f.rabbits, rab.Location())
		succ := rab.TryCatch(loc)
		// We must update the table, else we can run into two rabbits.
		f.rabbits[rab.Location()] = rab
		f.caughtCount++
		return succ
	}

	return false
}

// Attempts to tag a rabbit if it's still where we are.
func (f *directoryForest) PerformTag(tag string) bool {
	loc, _ := os.Getwd()

	rab, ok := f.rabbits[loc]

	if ok {
		delete(f.rabbits, rab.Location())
		succ := rab.TryTag(loc, tag)
		// We must update the table, else we can run into two rabbits.
		f.rabbits[rab.Location()] = rab
		return succ
	}

	return false
}

// Repopulated the forest if under the minimum number of rabbits
// we want. Otherwise, chance a rabbit will spawn.
func (f *directoryForest) repopulate() {
	for len(f.rabbits) < MinRabbits {
		r := NewRabbit(f)
		f.rabbits[r.Location()] = &r
	}

	if chance(SpawnChance) {
		r := NewRabbit(f)
		f.rabbits[r.Location()] = &r
	}
}

type forest struct {
	Rabbits		map[string]*Rabbit
	SpottedCount	uint
	CaughtCount	uint
	KilledCount	uint
}

// These are implemented because we can't encode private fields.
func (f *directoryForest) UnmarshalJSON(b []byte) error {
	data := forest{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	f.rabbits = data.Rabbits
	f.spottedCount = data.SpottedCount
	f.caughtCount = data.CaughtCount
	f.killedCount = data.KilledCount

	// Circular reference. Couldn't marshal their home so
	// we do it here.
	for _, r := range f.rabbits {
		(&r).ChangeHome(f)
	}
	return nil
}

func (f *directoryForest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forest{
		Rabbits:	f.rabbits,
		SpottedCount:	f.spottedCount,
		CaughtCount:	f.caughtCount,
		KilledCount:	f.killedCount,
	})
}
