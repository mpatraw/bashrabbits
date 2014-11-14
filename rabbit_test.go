
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

func (tf TestForest) FarawayLocation(loc string) string {
	return "far"
}

func TestMoving(t *testing.T) {
	tf := TestForest{}
	r := NewRabbit(tf)
	r.setIdleTime(time.Millisecond)
	r.setFleeTime(time.Millisecond)
	if r.Location() != "far" {
		t.Errorf("rabbit did not move (%s!=%s)", r.Location(), "far");
	}

	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt("far1")
	if !r.JustSpotted() {
		t.Errorf("rabbit was not spotted");
	}

	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt("somewhere")
	if r.Location() != "far" {
		t.Errorf("rabbit did not flee somewhere far")
	}
}

func TestDirectoryForest(t *testing.T) {
	f := newDirectoryForest()
	t.Logf("%v\n", f);
}

func TestUtil(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := randRange(1, 6)
		if r < 1 || r > 6 {
			t.Errorf("randRange() returned out of range (%d)", r)
		}
	}

	p1 := "/home/grue"
	p2 := "/home/grue/data"
	if !isAscension(p1, p2) {
		t.Errorf("%s does not ascend from %s\n", p1, p2)
	}
	if !isDescension(p2, p1) {
		t.Errorf("%s does not descend from %s\n", p2, p1)
	}

	p3 := "/home/grue"
	if isAscension(p1, p3) || isDescension(p1, p3) {
		t.Errorf("%s ascends or descends from %s\n", p1, p3)
	}
}
