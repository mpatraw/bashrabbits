
package main

import "bytes"
import "testing"

type TestForest struct { }

func (tf TestForest) LocationExists(loc string) bool {
	return true
}

func (tf TestForest) NearbyLocation(loc string) string {
	var buffer bytes.Buffer
	buffer.WriteString(loc)
	buffer.WriteString("a")
	return buffer.String()
}

func (tf TestForest) FarawayLocation() string {
	return "far"
}

func TestMoving(t *testing.T) {
	tf := TestForest{}
	r := NewRabbit(tf)
	t.Logf("%v\n", r)
}