
package main

import (
	"encoding/json"
	"time"
)

// A state the rabbit can be in.
type RabbitState uint

// An event that can be performed on a rabbit.
type RabbitAction uint

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

const (
	// Default action. This is performed on every rabbit.
	Wait RabbitAction = iota
	// A rabbit is spotted.
	Spot
	// Force a fee.
	Flee
	// When a catch attempt succeeds.
	Catch
	// When a rabbit dies. :(
	Kill
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

var rMachine Machine

func init() {
	// Create the rabbit state machine.
	rMachine = NewMachine()

	rMachine.AddTransition(State(Wandering), Action(Wait), State(Wandering))
	rMachine.AddTransition(State(Wandering), Action(Spot), State(Spotted))
	rMachine.AddTransition(State(Wandering), Action(Catch), State(Caught))
	rMachine.AddTransition(State(Wandering), Action(Kill), State(Dead))

	rMachine.AddTransition(State(Spotted), Action(Wait), State(Fleeing))
	// Can't spot an already spotted rabbit.
	rMachine.AddTransition(State(Spotted), Action(Flee), State(Fleeing))
	rMachine.AddTransition(State(Spotted), Action(Catch), State(Caught))
	rMachine.AddTransition(State(Spotted), Action(Kill), State(Dead))

	rMachine.AddTransition(State(Fleeing), Action(Wait), State(Wandering))
	// Can't spot or catch a fleeing rabbit.
	rMachine.AddTransition(State(Fleeing), Action(Kill), State(Dead))
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

// Step 1 for becoming Stateful.
func (r *Rabbit) State() State {
	return State(r.state)
}

// Step 2 for becoming Stateful.
func (r *Rabbit) ShouldTransition(act Action, to State) bool {
	ract := act.(RabbitAction)
	rstate := to.(RabbitState)

	switch ract {
	case Wait:
		if rstate == Wandering {
			return time.Now().Sub(r.lastMoved) >= r.idleTime
		} else if rstate == Fleeing {
			return time.Now().Sub(*r.lastSpotted) >= r.fleeTime
		} else {
			panic("Waiting when not wandering or fleeing.")
		}
	case Catch:
		elapsed := time.Now().Sub(*r.lastSpotted)
		catchchance := 1.0 - float64(elapsed) / float64(FleeTime)
		return chance(catchchance)
	default:
		return true
	}
}

// Step 3 for becoming Stateful.
func (r *Rabbit) EnterState(state State) {
	rstate := state.(RabbitState)

	switch rstate {
	case Wandering:
		r.lastMoved = time.Now()
		r.lastLocation = r.location
		r.location = r.home.NearbyLocation(r.location)
		r.state = rstate
	case Spotted:
		// Uh-oh!
		r.state = rstate
		t := time.Now()
		r.lastSpotted = &t
		// Will start to flee the next update.
	case Fleeing:
		r.lastMoved = time.Now()
		r.lastLocation = r.location
		r.location = r.home.FarawayLocation(r.location)
		r.state = rstate
	case Caught:
		r.location = ""
		r.state = rstate
	case Dead:
		r.location = ""
		r.state = rstate
	default:
	}
}

// This is called before every operation. Returns true if the rabbit
// awake and ready.
func (r *Rabbit) wakeup() bool {
	if !r.IsPlaying() {
		return false
	}
	rMachine.Perform(r, Wait)
	return true
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
	if !r.wakeup() {
		return
	}
	if r.location == loc {
		rMachine.Perform(r, Spot)
	}
}

// Attempts to catch the rabbit. The rabbit first checks if
// it already moved with wakeup(). The chance to catch the
// rabbit is the inverse of the time is has left before moving.
func (r *Rabbit) TryCatch(loc string) bool {
	if !r.wakeup() {
		return false
	}

	if r.location != loc {
		return false
	}

	if !rMachine.Perform(r, Catch) {
		// Oh-well, better luck next time.
		rMachine.Perform(r, Flee)
		return false
	}
	return true
}

// Attempts to tag the rabbit. Right now it's a 100% chance.
func (r *Rabbit) TryTag(loc, tag string) bool {
	if !r.wakeup() {
		return false
	}

	if r.location != loc {
		return false
	}

	r.tag = tag
	rMachine.Perform(r, Flee)
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

// Returns true if this rabbit has JUST been spotted. This
// state will immediately move to the fleeing state.
func (r *Rabbit) JustSpotted() bool {
	return r.state == Spotted
}

// Returns true if the rabbit is apart of the game. The rabbit
// is no longer playing if it's caught/dead/etc. Or if the location
// the rabbit is no longer exists.
func (r *Rabbit) IsPlaying() bool {
	if !r.home.LocationExists(r.location) {
		rMachine.Perform(r, Kill)
		return false
	}
	return r.state != Dead && r.state != Caught
}

// Used for marshalling/unmarshalling.
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
