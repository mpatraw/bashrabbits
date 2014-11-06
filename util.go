package main

import (
	"encoding/binary"
	"crypto/rand"
	"io/ioutil"
	"math"
	//"os"
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

func randRange(low, high uint) uint {
	f := randFloat() * float64(high - low + 1)
	return uint(math.Floor(f)) + low
}

func listDirs(path string) []string {
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

func canAscend(path string) bool {
	dirs := listDirs(path)
	return len(dirs) > 0
}

func randAscension(path string) string {
	dirs := listDirs(path)
	if len(dirs) == 0 {
		panic("Tried to ascend when unable")
	}
	return dirs[randRange(0, uint(len(dirs) - 1))]
}