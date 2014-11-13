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
        sm.Perform(&s, "Run")
        sm.Perform(&s, "Run")
        sm.Perform(&s, "Run")
}
