
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
	bytes, err := json.MarshalIndent(df, "", "\t")
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

	savefile := filepath.Join(os.Getenv("HOME"), ".rabbit")
	df := loadDirectoryForest(savefile)
	defer saveDirectoryForest(savefile, df)

	if flag.NArg() == 0 {
		usage()
		return
	}

	switch flag.Arg(0) {
	case "stats":
		fmt.Printf("Rabbits\n");
		fmt.Printf("...spotted:    %d\n", df.spottedCount)
		fmt.Printf("...caught:     %d\n", df.caughtCount)
		fmt.Printf("...killed:     %d :(\n", df.killedCount)
	case "check":
		spotted := df.PerformCheck()
		if spotted != nil {
			if spotted.Tag() != "" {
				fmt.Printf("You see the %s rabbit!\n", spotted.Tag())
			} else if spotted.SeenBefore() {
				fmt.Printf("You see a familiar rabbit here.\n")
			} else {
				fmt.Printf("A rabbit is here!!\n")
			}
		}
	case "catch":
		if df.PerformCatch() {
			fmt.Printf("You caught the rabbit!\n")
		} else {
			fmt.Printf("The rabbit got away...\n")
		}
	case "tag":
		if flag.NArg() < 2 {
			usage()
			return
		}

		if df.PerformTag(flag.Arg(1)) {
			fmt.Printf("You successfully tagged the rabbit!\n")
		} else {
			fmt.Printf("The rabbit got away...\n")
		}
	default: usage()
	}
}
