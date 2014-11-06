
package main

import (
	"bytes"
	"testing"
	"time"
)

type TestForest struct { }

func (tf TestForest) LocationExists(loc string) bool {
	return true
}

func (tf TestForest) NearbyLocation(loc string) string {
	var buffer bytes.Buffer
	buffer.WriteString(loc)
	buffer.WriteString("1")
	return buffer.String()
}

func (tf TestForest) FarawayLocation() string {
	return "far"
}

func TestMoving(t *testing.T) {
	tf := TestForest{}
	r := NewRabbit(tf)
	r.setIdleTime(time.Millisecond)
	r.setFleeTime(time.Millisecond)
	
	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt("yo")
	if r.Location() != "far1" {
		t.Errorf("rabbit did not move (%s==%s)", r.Location(), "far1");
	}
	
	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt("far11")
	if !r.JustSpotted() {
		t.Errorf("rabbit was not spotted");
	}
	
	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt("somewhere")
	if r.Location() != "far" {
		t.Errorf("rabbit did not flee somewhere far")
	}
}