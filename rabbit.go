
package main

import (
	"time"
)

// A state the rabbit can be in.
type RabbitState int

const (
	// Initial state.
	Wandering RabbitState = iota
	// The rabbit is in a fleeing state when it's spotted.
	Fleeing
	// If the rabbit was successfully caught.
	Caught
	// If the rabbit died this will be the state. The rabbit only dies
	// if the location it's in no longer exists.
	Dead
)

// The time that elapses before a rabbit wants to moved.
const IdleTime = time.Duration(10) * time.Minute
// The time that elapses before a rabbit moves after being spotted.
const FleeTime = time.Duration(5) * time.Second

// A forest is a place that can be traversed. Locations in a forest
// are simple strings.
type Forest interface {
	// Returns true if passed location exists.
	LocationExists(loc string) bool
	// Returns a location fairly close to the one provided.
	NearbyLocation(loc string) string
	// Returns a faraway location, this could be anywhere
	// except the location passed (unless it's the only location).
	FarawayLocation(loc string) string
}

// A rabbit is a simple creature that likes to move around a forest. You can
// spot it, try to catch it, tag it, or accidentally kill it. :(
type Rabbit struct {
	// The forest the rabbit lives in.
	home		Forest		`json:"-"`
	// The current location in the forest. May be "", in which
	// case the rabbit is no longer in the forest (dead, caught).
	location	string		`json:"location"`
	// A tag identifying this specific rabbit.
	tag		string		`json:"tag"`
	// The last location visited. May be "", in which case the
	// rabbit never moved.
	lastLocation	string		`json:"lastLocation"`
	// The last time the rabbit moved to new location.
	lastMoved	time.Time	`json:"lastMoved"`
	// The time the rabbit was spotted last. May be nil.
	lastSpotted	*time.Time	`json:"lastSpotted"`
	// State of the rabbit.
	state		RabbitState	`json:"state"`
	
	// These are set to the defaults.
	idleTime	time.Duration	`json:"idleTime"`
	fleeTime	time.Duration	`json:"fleeTime"`
}

// Creates a new rabbit and moves it to a faraway location.
func NewRabbit(f Forest) Rabbit {
	r := Rabbit{
		f, "", "", "", time.Now(), nil, Wandering,
		IdleTime, FleeTime,
	}
	r.location = f.FarawayLocation("")
	return r
}

// This is called before every operation. The rabbit occasionally
// moves.
func (r *Rabbit) wakeup() {
	if r.CantMove() {
		return
	}
	if !r.home.LocationExists(r.location) {
		r.state = Dead
		r.location = ""
		return
	}
	
	shouldMove := false
	now := time.Now()
	elapsed := now.Sub(r.lastMoved)
	var to string
	
	if r.state == Fleeing {
		if elapsed >= r.fleeTime {
			shouldMove = true
			// After fleeing the rabbit goes FAR away.
			to = r.home.FarawayLocation(r.location)
		}
	} else {
		if elapsed >= r.idleTime {
			shouldMove = true
			to = r.home.NearbyLocation(r.location)
		}
	}
	
	if shouldMove {
		r.lastMoved = now
		r.lastLocation = r.location
		r.location = to
		// Stop fleeing, or whatever we were doing.
		r.state = Wandering
	}
}

// Used mostly for testing. The default is preferred.
func (r *Rabbit) setIdleTime(d time.Duration) {
	r.idleTime = d
}

// Used mostly for testing. The default is preferred.
func (r *Rabbit) setFleeTime(d time.Duration) {
	r.fleeTime = d
}

// Changes the home of the rabbit.
func (r *Rabbit) ChangeHome(f Forest) {
	r.home = f
}

// A place in the forest was disturbed. Possibly move, or
// if the place is here, the rabbit is spotted.
func (r *Rabbit) DisturbanceAt(loc string) {
	r.wakeup()

	// Uh-oh!
	if r.location == loc {
		t := time.Now()
		r.lastSpotted = &t
		r.state = Fleeing	
	}
}

// Attempts to catch the rabbit. The rabbit first check if
// it already moved with wakeup(). The chance to catch the
// rabbit is the inverse of the time is has left before moving.
func (r *Rabbit) TryCatch(loc string) bool {
	r.wakeup()

	if r.location != loc {
		return false
	}

	elapsed := time.Now().Sub(*r.lastSpotted)
	chance := 1.0 - float64(elapsed) / float64(FleeTime)
	if randFloat() < chance {
		r.state = Caught
		r.location = ""
	}
	
	return true
}

// Attempts to tag the rabbit. Right now it's a 100% chance.
func (r *Rabbit) TryTag(loc, tag string) bool {
	r.wakeup()

	if r.location != loc {
		return false
	}
	
	r.tag = tag
	return true
}

// Returns the current location of the rabbit.
func (r *Rabbit) Location() string {
	return r.location
}

// Returns the current tag of the rabbit, "" is none.
func (r *Rabbit) Tag() string {
	return r.tag
}

// Returns true if this rabbit has been seen before.
func (r *Rabbit) SeenBefore() bool {
	return r.lastSpotted != nil
}

// Returns true if this rabbit has JUST been spotted, it should
// be fleeing.
func (r *Rabbit) JustSpotted() bool {
	return r.lastSpotted != nil && r.state == Fleeing
}

// Returns the state of the rabbit.
func (r *Rabbit) State() RabbitState {
	return r.state
}

// Returns true if the rabbit can't move, usually because it
// is dead or caught.
func (r *Rabbit) CantMove() bool {
	return r.state == Dead || r.state == Caught
}