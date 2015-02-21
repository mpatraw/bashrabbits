
package main

import (
	"encoding/json"
	"os"
	"time"
)

// The direction of the tracks indicates where a rabbit went
// from here.
type TrackDirection uint

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

	// How long it takes for tracks to fade. Right
	// now, it's a 1/5 of the time it takes a rabbit
	// to move.
	TrackFadeTime	= IdleTime / 5
)

const (
	// No tracks.
	TrackNone TrackDirection = iota
	// Ascending is going "up" a directory, like cd ..
	TrackAscending
	// Descending is going "down" a directory, like cd ./data
	TrackDescending
)

type track struct {
	Timestamp	time.Time
	Direction	TrackDirection
}

type directoryForest struct {
	// List of rabbits and their locations. Only one
	// rabbit per location.
	rabbits		map[string]*Rabbit
	// Tracks at a given location. Cleared and updated after every move.
	tracks		map[string]track
	// Number of rabbits seen.
	spottedCount	uint
	// Number of rabbits caught.
	caughtCount	uint
	// Number of rabbits killed. :(
	killedCount	uint
}

func newDirectoryForest() directoryForest {
	return directoryForest{
		map[string]*Rabbit{}, map[string]track{}, 0, 0, 0,
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
	added := []string{}
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

		added = append(added, newloc)
	}

	if newloc == loc {
		// Guaranteed to not be the same because you must
		// step twice to get to the same destination.
		steps = 1
		goto tryagain
	}

	pastLoc := loc
	for _, aloc := range added {
		if isAscension(aloc, pastLoc) {
			f.tracks[pastLoc] = track{time.Now(), TrackAscending}
		} else if isDescension(aloc, pastLoc) {
			f.tracks[pastLoc] = track{time.Now(), TrackDescending}
		} else {
			panic("Rabbit didn't move to nearby location.")
		}
		pastLoc = aloc
	}

	return newloc

}

// A random faraway location. Rabbits typically start here
// and run here when they're fleeing. Faraway locations don't
// add tracks.
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

// Returns true if a rabbit is here. Only useful for checking
// before performing an action.
func (f *directoryForest) IsRabbitHere() bool {
	loc, _ := os.Getwd()
	_, ok := f.rabbits[loc]
	return ok
}

// Returns whether tracks
func (f *directoryForest) GetTracksHere() (bool, TrackDirection) {
	loc, _ := os.Getwd()
	t, ok := f.tracks[loc]
	if ok {
		return true, t.Direction
	} else {
		return false, TrackNone
	}
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

		if (r.IsPlaying()) {
			if r.JustSpotted() {
				// It's possible for two rabbits to "wakeup"
				// to the same location in the same update.
				// Right now the most recent on will not be
				// overridden so we spotted that one.
				//
				// XXX: Fix rabbits running into each other?
				spotted = r
				f.spottedCount++
			}
			newrabbits[r.Location()] = r
		} else {
			if r.State() == Dead {
				f.killedCount++
			} else if r.State() == Caught {
				// Update in PerformCatch, otherwise
				// catching a rabbit score won't
				// be update until the next call to this
				// function.
			}
		}
	}

	f.rabbits = newrabbits

	// See if we should repopulate.
	f.repopulate()

	f.fadeTracks()

	return
}

// Attempts to catch a rabbit if it's still where we are.
func (f *directoryForest) PerformCatch() bool {
	loc, _ := os.Getwd()

	f.fadeTracks()

	rab, ok := f.rabbits[loc]
	if ok {
		delete(f.rabbits, rab.Location())
		succ := rab.TryCatch(loc)
		// We must update the table, else we can run into two rabbits.
		f.rabbits[rab.Location()] = rab
		if succ {
			f.caughtCount++
		}
		return succ
	}

	return false
}

// Attempts to tag a rabbit if it's still where we are.
func (f *directoryForest) PerformTag(tag string) bool {
	loc, _ := os.Getwd()

	f.fadeTracks()

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

// Fades the tracks depending on how old they are. Faded
// tracks are removed.
func (f *directoryForest) fadeTracks() {
	list := []string{}
	for loc, track := range f.tracks {
		age := time.Now().Sub(track.Timestamp)
		if age >= TrackFadeTime {
			list = append(list, loc)
		}
	}

	for _, loc := range list {
		delete(f.tracks, loc)
	}
}

// Used for marshalling/unmarshalling.
type forest struct {
	Rabbits		map[string]*Rabbit
	Tracks		map[string]track
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
	f.tracks = data.Tracks
	f.spottedCount = data.SpottedCount
	f.caughtCount = data.CaughtCount
	f.killedCount = data.KilledCount

	// Circular reference. Couldn't marshal their home so
	// we do it here.
	for _, r := range f.rabbits {
		r.ChangeHome(f)
	}
	return nil
}

func (f *directoryForest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forest{
		Rabbits:	f.rabbits,
		Tracks:		f.tracks,
		SpottedCount:	f.spottedCount,
		CaughtCount:	f.caughtCount,
		KilledCount:	f.killedCount,
	})
}
