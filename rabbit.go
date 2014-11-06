
package main

import (
	"encoding/binary"
	"crypto/rand"
	"math"
	"time"
)

func randFloat() float64 {
	b := make([]byte, 8)
	rand.Read(b)
	bits := binary.BigEndian.Uint64(b)
	return math.MaxFloat64 / math.Float64frombits(bits)
}

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
	// Returns a faraway location, this could be anywhere.
	FarawayLocation() string
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
}

// Creates a new rabbit and moves it to a faraway location.
func NewRabbit(f Forest) *Rabbit {
	r := &Rabbit{
		f, "", "", "", time.Now(), nil, Wandering,
	}
	r.location = f.FarawayLocation()
	return r
}

// Changes the home of the rabbit.
func (r *Rabbit) ChangeHome(f Forest) {
	r.home = f
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
	
	if r.state == Fleeing {
		if elapsed >= FleeTime {
			shouldMove = true
		}
	} else {
		if elapsed >= IdleTime {
			shouldMove = true
		}
	}
	
	if shouldMove {
		r.lastMoved = now
		r.lastLocation = r.location
		r.location = r.home.NearbyLocation(r.lastLocation)
		// Stop fleeing, or whatever we were doing.
		r.state = Wandering
	}
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

func (r *Rabbit) TryTag(loc, tag string) bool {
	r.wakeup()

	if r.location != loc {
		return false
	}
	
	r.tag = tag
	return true
}

func (r *Rabbit) Location() string {
	return r.location
}

func (r *Rabbit) Tag() string {
	return r.tag
}

func (r *Rabbit) SeenBefore() bool {
	return r.lastSpotted != nil
}

func (r *Rabbit) JustSpotted() bool {
	return r.lastSpotted != nil && r.state == Fleeing
}

func (r *Rabbit) State() RabbitState {
	return r.state
}

func (r *Rabbit) CantMove() bool {
	return r.state == Dead || r.state == Caught
}