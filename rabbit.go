
package main

import (
	"encoding/json"
	"time"
)

// A state the rabbit can be in.
type RabbitState int

const (
	// Initial state.
	Wandering RabbitState = iota
	// Transient state. The rabbit will move to fleeing.
	Spotted
	// The rabbit is in a fleeing state when it's spotted.
	Fleeing
	// If the rabbit was successfully caught.
	Caught
	// If the rabbit died this will be the state. The rabbit only dies
	// if the location it's in no longer exists.
	Dead
)

// The time that elapses before a rabbit wants to moved.
const IdleTime = time.Duration(5) * time.Minute
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
	home		Forest
	// The current location in the forest. May be "", in which
	// case the rabbit is no longer in the forest (dead, caught).
	location	string
	// A tag identifying this specific rabbit.
	tag		string
	// The last location visited. May be "", in which case the
	// rabbit never moved.
	lastLocation	string
	// The last time the rabbit moved to new location.
	lastMoved	time.Time
	// The time the rabbit was spotted last. May be nil.
	lastSpotted	*time.Time
	// State of the rabbit.
	state		RabbitState

	// These are set to the defaults.
	idleTime	time.Duration
	fleeTime	time.Duration
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

// This is called when the rabbit flees, either after
// FleeTime ran up, or a failed catch, or a tag.
func (r *Rabbit) flee() {
	if r.state != Fleeing {
		panic("Tried to call flee() when not fleeing.")
	}
	r.lastMoved = time.Now()
	r.lastLocation = r.location
	// Run far away! We're scared.
	r.location = r.home.FarawayLocation(r.location)
	r.state = Wandering
}

// This is called when the rabbit wanders after IdleTime is
// up.
func (r *Rabbit) wander() {
	if r.state != Wandering {
		panic("Tried to call wander() when not wandering.")
	}
	r.lastMoved = time.Now()
	r.lastLocation = r.location
	r.location = r.home.NearbyLocation(r.location)
}

// This is called before every operation. The rabbit occasionally
// moves.
func (r *Rabbit) wakeup() {
	if !r.CanMove() {
		return
	}
	if !r.home.LocationExists(r.location) {
		r.state = Dead
		r.location = ""
		return
	}

	// Immediately move to the fleeing state. The rabbit
	// hasn't JUST been spotted, it's been sitting here.
	if r.state == Spotted {
		r.state = Fleeing
	}

	if r.state == Fleeing {
		if time.Now().Sub(*r.lastSpotted) >= r.fleeTime {
			r.flee()
		}
	} else {
		if time.Now().Sub(r.lastMoved) >= r.idleTime {
			r.wander()
		}
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
		r.state = Spotted
		t := time.Now()
		r.lastSpotted = &t
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
	catchchance := 1.0 - float64(elapsed) / float64(FleeTime)
	if chance(catchchance) {
		r.state = Caught
		r.location = ""
	} else {
		// Oh-well, better luck next time.
		r.flee()
		return false
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
	r.flee()
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

// Returns true if this rabbit has JUST been spotted. This
// state will immediately move to the fleeing state.
func (r *Rabbit) JustSpotted() bool {
	return r.state == Spotted
}

// Returns the state of the rabbit.
func (r *Rabbit) State() RabbitState {
	return r.state
}

// Returns true if the rabbit can't move, usually because it
// is dead or caught.
func (r *Rabbit) CanMove() bool {
	return r.state != Dead && r.state != Caught
}

type rabbit struct {
	Location	string
	Tag		string
	LastLocation	string
	LastMoved	time.Time
	LastSpotted	*time.Time
	State		RabbitState
	IdleTime	time.Duration
	FleeTime	time.Duration
}

func (r *Rabbit) UnmarshalJSON(b []byte) error {
	data := rabbit{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	r.location = data.Location
	r.tag = data.Tag
	r.lastLocation = data.LastLocation
	r.lastMoved = data.LastMoved
	r.lastSpotted = data.LastSpotted
	r.state = data.State
	r.idleTime = data.IdleTime
	r.fleeTime = data.FleeTime
	return nil
}

func (r *Rabbit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&rabbit{
		Location: r.location,
		Tag: r.tag,
		LastLocation: r.lastLocation,
		LastMoved: r.lastMoved,
		LastSpotted: r.lastSpotted,
		State: r.state,
		IdleTime: r.idleTime,
		FleeTime: r.fleeTime,
	})
}
