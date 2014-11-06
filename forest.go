
package main

import (
	"os"
)

const (
	// Default value.
	MaxRabbits	= 30
)

type directoryForest struct {
	// List of rabbits and their locations. Only one
	// rabbit per location.
	rabbits		map[string]Rabbit
	// Number of rabbits currently in the forest.
	numRabbits	uint
	// Max number of rabbits allowed in the forest. Configurable.
	maxRabbits	uint
}

func newDirectoryForest() directoryForest {
	return directoryForest{
		make(map[string]Rabbit), 0, MaxRabbits,
	}
}

func (f *directoryForest) LocationExists(loc string) bool {
	fi, err := os.Stat(loc)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (f *directoryForest) NearbyLocation(loc string) string {
	return ""
}

func (f *directoryForest) FarawayLocation() string {
	return ""
}