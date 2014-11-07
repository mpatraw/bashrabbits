
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ascii bool

func init() {
	flag.BoolVar(&ascii, "a", false, "use ascii art instead of words")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rabbit [-a] [stats|check|catch|tag string]\n")
	flag.PrintDefaults()
}

// Loads the directory forest from the file passed, if it doesn't exist
// a fresh directory forest is returned.
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
	defer file.Close()

	fz, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fz.Close()

	bytes, err := ioutil.ReadAll(fz)
	if err != nil {
		log.Fatal(err)
	}

	var df directoryForest
	err = json.Unmarshal(bytes, &df)
	if err != nil {
		log.Fatal(err)
	}

	return &df
}

// Saves the directory forest to a file.
func saveDirectoryForest(filename string, df *directoryForest) {
	bs, err := json.Marshal(df)
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(bs)
	w.Close()
	err = ioutil.WriteFile(filename, b.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Returns flavor for the number of spotted rabbits.
func spottedFlavor(count uint) string {
	switch {
	case count < 20: return ""
	case count < 50: return ":)"
	default: return ":D"
	}
}

// Returns flavor for the number of caught rabbits.
func caughtFlavor(count uint) string {
	switch {
	case count < 5: return ""
	case count < 20: return ":)"
	default: return ":D"
	}
}

// Returns flavor for the number of killed rabbits. The higher, the
// more dramatic the flavor.
func killedFlavor(count uint) string {
	switch {
	case count < 5: return ""
	case count < 20: return ":("
	case count < 50: return ";("
	default: return "MONSTER!!"
	}
}

// Prints the stats. Number of rabbits seen, caught, killed, etc.
func printStats(df *directoryForest) {
	sflavor := spottedFlavor(df.spottedCount)
	cflavor := caughtFlavor(df.caughtCount)
	kflavor := killedFlavor(df.killedCount)
	fmt.Printf("Rabbits\n");
	fmt.Printf("...spotted:    %d %s\n", df.spottedCount, sflavor)
	fmt.Printf("...caught:     %d %s\n", df.caughtCount, cflavor)
	fmt.Printf("...killed:     %d %s\n", df.killedCount, kflavor)
}

func printRabbit(state RabbitState) {
	switch state {
	case Wandering:
		fmt.Printf(" ()_()\n")
		fmt.Printf(" (-.-)\n")
		fmt.Printf("'(\"|\")'\n")
	case Spotted:
		fmt.Printf("(_/  _#\n")
		fmt.Printf("'.'_( )\n")
		//fmt.Printf("/)/)\n")
		//fmt.Printf("(o.o)\n")
		//fmt.Printf("c(")(")\n")
	case Fleeing:
		fmt.Printf("  o __(\\\\\n")
		fmt.Printf("   ) _ --\n")
		fmt.Printf(" //    \\\\\n")
	case Caught:
		fmt.Printf("_________\n")
		fmt.Printf("| ()|() |\n")
		fmt.Printf("+---+---+\n")
		fmt.Printf("|(\")|(\")|\n")
		fmt.Printf("---------\n")
	case Dead:
		fmt.Printf("(\\ /)\n")
		fmt.Printf("(x.x)\n")
		fmt.Printf("(> <)\n")
	}
}

// Check the current directory for rabbits.
func check(df *directoryForest) {
	spotted := df.PerformCheck()
	if spotted != nil {
		if spotted.Tag() != "" {
			fmt.Printf("You see the %s rabbit!\n", spotted.Tag())
			if ascii {
				printRabbit(Spotted)
			}
		} else {
			fmt.Printf("A rabbit is here!!\n")
			if ascii {
				printRabbit(Spotted)
			}
		}
	}
}

// Try to catch a rabbit.
func catch(df *directoryForest) {
	if df.IsRabbitHere() {
		if df.PerformCatch() {
			fmt.Printf("You caught the rabbit!\n")
			if ascii {
				printRabbit(Caught)
			}
		} else {
			fmt.Printf("The rabbit got away...\n")
			if ascii {
				printRabbit(Fleeing)
			}
		}
	} else {
		fmt.Printf("There are no rabbits here.\n")
	}
}

// Try to tag a rabbit.
func tag(df *directoryForest, tag string) {
	if df.IsRabbitHere() {
		if df.PerformTag(tag) {
			fmt.Printf("You successfully tagged the rabbit!\n")
			if ascii {
				printRabbit(Wandering)
			}
		} else {
			fmt.Printf("The rabbit got away...\n")
			if ascii {
				printRabbit(Caught)
			}
		}
	} else {
		fmt.Printf("There are no rabbits here.\n")
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
		printStats(df)
	case "check":
		check(df)
	case "catch":
		catch(df)
	case "tag":
		if flag.NArg() < 2 {
			usage()
			return
		}
		tag(df, flag.Arg(1))
	default: usage()
	}
}
