package main

import (
	"encoding/binary"
	"crypto/rand"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Uses /dev/urandom to generate random numbers. We don't
// need to recreate generated numbers, so we don't save
// a RNG state.
func randUint() uint64 {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	for err != nil {
		_, err = rand.Read(b)
	}
	return binary.BigEndian.Uint64(b)
}

func randFloat() float64 {
	return float64(randUint()) / float64(math.MaxUint64)
}

// The range returned is inclusive.
func randRange(low, high uint) uint {
	f := randFloat() * float64(high - low + 1)
	return uint(math.Floor(f)) + low
}

// Returns true when under the specified chance.
func chance(f float64) bool {
	return randFloat() < f
}

// Returns the directory listing as full path names. The passed path
// must be absolute.
func listDirs(path string) []string {
	if !filepath.IsAbs(path) {
		panic("cannot list dirs on non-absolute path")
	}

	dirs := []string{}
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		isPrivate := strings.HasPrefix(file.Name(), ".")
		if file.IsDir() && !isPrivate {
			dirs = append(dirs, filepath.Join(path, file.Name()))
		}
	}
	return dirs
}

// Returns true if you can descend from this path, descending is going
// down a directory, as opposed to up (`cd ..` is up). The passed path
// must be absolute
func canDescend(path string) bool {
	dirs := listDirs(path)
	return len(dirs) > 0
}

// Returns a random path to desend. The passed path must be absolute.
func randDescension(path string) string {
	dirs := listDirs(path)
	if len(dirs) == 0 {
		panic("Tried to descend when unable")
	}
	return dirs[randRange(0, uint(len(dirs) - 1))]
}

// Returns true if you can ascend from this path. No ascending
// below the home directory. The passed path must be absolute.
func canAscend(path string) bool {
	home := os.Getenv("HOME")
	return strings.HasPrefix(filepath.Dir(path), home)
}

// No need to be random. You can only ascend in one direction.
func ascend(path string) string {
	return filepath.Dir(path)
}

// This is the furthest we can ascend.
func baseLocation() string {
	home := os.Getenv("HOME")
	return home
}
