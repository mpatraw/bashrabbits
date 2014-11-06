
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	
}

func usage() {
	fmt.Printf("usage: rabbit [stats|check|catch|tag string]\n")
	os.Exit(0)
}

func loadDirectoryForest(filename string) *directoryForest {
	
	file, err := os.Open(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		} else {
			df := newDirectoryForest()
			return &df
		}
	}
	var df directoryForest
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &df)
	if err != nil {
		log.Fatal(err)
	}
	return &df
}

func saveDirectoryForest(filename string, df *directoryForest) {
	bytes, err := json.Marshal(df)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	
	if flag.NArg() == 0 {
		usage()
	}
	
	for i, arg := range flag.Args() {
		fmt.Printf("Arg[%d]=%s\n", i, arg)
	}
	
	savefile := filepath.Join(os.Getenv("HOME"), ".rabbit")
	df := loadDirectoryForest(savefile)
	defer saveDirectoryForest(savefile, df)
	
	switch flag.Arg(0) {
	case "stats":
	case "check":
	case "catch":
	case "tag":
		if flag.NArg() < 2 {
			usage()
		}
	default: usage()
	}
}