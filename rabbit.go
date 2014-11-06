
package main

import "time"

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

type RabbitState int

const (
	// Initial state.
	Wandering RabbitState = iota
	// If the rabbit was successfully caught.
	Caught
	// If the rabbit died this will be the state. The rabbit only dies
	// if the location it's in no longer exists.
	Dead
)

// A rabbit is a simple creature that likes to move around a forest. You can
// spot it, try to catch it, tag it, or accidentally kill it. :(
type Rabbit struct {
	// The forest the rabbit lives in.
	home		Forest
	// The current location in the forest.
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
}

// Creates a new rabbit and moves it to a faraway location.
func NewRabbit(f Forest) *Rabbit {
	r := &Rabbit{
		f, "", "", "", time.Now(), nil, Wandering,
	}
	r.location = f.FarawayLocation()
	return r
}

// The rabbit wakes up and and decides if it wants to move.
func (r *Rabbit) Wakeup() {
	if !r.IsWandering() {
		return
	}
	if !r.LocationExists(r.location) {
		r.state = Dead
		return
	}
}

// The rabbit was spotted, make it flee.
func (r *Rabbit) Spot() {
	t := time.Now()
	r.lastSpotted = &t

}

func (r *Rabbit) TryCatch() {

}

func (r *Rabbit) TryTag() {

}

func (r *Rabbit) Location() string {
	return r.location
}

func (r *Rabbit) Tag() string {
	return r.tag
}

func (r *Rabbit) WasSpotted() bool {
	return r.lastSpotted != nil
}

func (r *Rabbit) State() RabbitState {
	return r.state
}

func (r *Rabbit) IsWandering() bool {
	return r.state == Wandering
}

func (r *Rabbit) IsCaught() bool {
	return r.state == Caught
}

func (r *Rabbit) IsDead() bool {
	return r.state == Dead
}