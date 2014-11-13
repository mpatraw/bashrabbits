package main

import (
	"fmt"
	"testing"
)

type TestStateful struct {
	state State
}

func (tm *TestStateful) State() State {
	return tm.state
}

func (tm *TestStateful) ShouldTransition(act Action, to State) bool {
	return true
}

func (tm *TestStateful) EnterState(state State) {
	fmt.Printf("%s -> ", tm.state)
	tm.state = state
	fmt.Printf("%s\n", tm.state)
}

func TestState(t *testing.T) {
	s := TestStateful{"Standing"}
	sm := NewMachine()
	sm.AddTransition("Standing", "Run", "Running")
	sm.AddTransition("Standing", "Walk", "Walking")
	sm.AddTransition("Running", "Run", "Tired")
	sm.AddTransition("Walking", "Walk", "Walking")
	sm.AddTransition("Walking", "Run", "Running")
	sm.AddTransition("Tired", "Rest", "Standing")
	sm.AddTransition("Tired", "Run", "Dead")

	sm.Perform(&s, "Walk")
	if s.state != "Walking" {
		t.Errorf("not walking (%s!=%s)", s.state, "Walking");
	}

	sm.Perform(&s, "Run")
	if s.state != "Running" {
		t.Errorf("not running (%s!=%s)", s.state, "Running");
	}

	sm.Perform(&s, "Run")
	if s.state != "Tired" {
		t.Errorf("not tired (%s!=%s)", s.state, "Tired");
	}

	sm.Perform(&s, "Run")
	if s.state != "Dead" {
		t.Errorf("not dead (%s!=%s)", s.state, "Dead");
	}
}
