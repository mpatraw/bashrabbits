
package main

import "time"

type forest interface {

}

type rabbit struct {
	home		*forest
	location	string
	tag		string
	lastLocation	string
	lastSpotted	*time.Time
}

func newRabbit(f *forest) *rabbit {
	return &rabbit{
		f, "", "", "", nil,
	}
}

func (r *rabbit) spot() {
	t := time.Now()
	r.lastSpotted = &t
}