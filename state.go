package main

// An event is an action string like, "Walk", "Run".
type Event string
// A state is a string like, "Walking", "Standing".
type State string

// Something stateful is something that can enter or exit
// states. It "has" a state.
type Stateful interface {
        State() State
        EnterState(State)
}

// A transition maps an event to a state.
type Transition struct {
        State State
        Event Event
}

// The state machine keeps track of the current state and the
// transition table.
type Machine struct {
	// Each state has a list of possible transitions.
	Transitions      map[Transition]State
}

// Returns a new state machine with a given state as the start.
func NewMachine() Machine {
        return Machine {
                Transitions: map[Transition]State{},
        }
}

// Adds a transition to the transition table.
func (machine *Machine) AddTransition(from State, ev Event, to State) {
	machine.Transitions[Transition{from, ev}] = to
}

// Deletes a transition from the transition table.
func (machine *Machine) DelTransition(from State, ev Event) {
	delete(machine.Transitions, Transition{from, ev})
}

// Performs an action on a stateful object using the state machine.
func (machine *Machine) Perform(ful Stateful, ev Event) {
	next, exists := machine.Transitions[Transition{ful.State(), ev}]
	if exists {
		ful.EnterState(next)
	}
}
