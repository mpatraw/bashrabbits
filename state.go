package main

// An action to change a state (maybe).
type Action interface{}
// A state of living.
type State interface{}

// Something stateful is something that can enter or exit
// states. It "has" a state.
type Stateful interface {
	State() State
	ShouldTransition(Action, State) bool
	EnterState(State)
}

// A transition maps an event to a state.
type Transition struct {
	State	State
	Action	Action
}

// The state machine keeps track of the current state and the
// transition table.
type Machine struct {
	// Each state has a list of possible transitions.
	Transitions	map[Transition]State
	// These transitions exist when the event doesn't matter.
	// Passed a "*" to AddTransition.
	AnyAction	map[State]State
	// These transitions exist when the state doesn't matter.
	// Passed a "*" to AddTransition.
	AnyState	map[Action]State
}

// Returns a new state machine with a given state as the start.
func NewMachine() Machine {
	return Machine {
		Transitions: map[Transition]State{},
		AnyAction: map[State]State{},
		AnyState: map[Action]State{},
	}
}

// Adds a transition to the transition table.
func (machine *Machine) AddTransition(from State, ev Action, to State) {
	if from == "*" {
		machine.AnyState[ev] = to
	} else if ev == "*" {
		machine.AnyAction[from] = to
	} else {
		machine.Transitions[Transition{from, ev}] = to
	}
}

// Deletes a transition from the transition table.
// Note, doesn't delete every from/ev in the case of "*".
func (machine *Machine) DelTransition(from State, ev Action) {
	if from == "*" {
		delete(machine.AnyState, ev)
	} else if ev == "*" {
		delete(machine.AnyAction, from)
	} else {
		delete(machine.Transitions, Transition{from, ev})
	}
}

// Performs an action on a stateful object using the state machine.
func (machine *Machine) Perform(ful Stateful, ev Action) bool {
	next, exists := machine.Transitions[Transition{ful.State(), ev}]
	if exists && ful.ShouldTransition(ev, next) {
		ful.EnterState(next)
		return true
	}
	return false
}
