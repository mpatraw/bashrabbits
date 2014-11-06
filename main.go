
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func init() {
	
}

func main() {
	flag.Parse()
	
	for i, arg := range flag.Args() {
		fmt.Printf("Arg[%d]=%s\n", i, arg)
	}
	
	dir, _ := os.Getwd()
	fmt.Printf("%s\n", dir)
	fi, _ := os.Stat(dir)
	fmt.Printf("%v %t\n", fi.Mode(), fi.IsDir())
	
	f := newDirectoryForest()
	r := NewRabbit(&f)
	r.setIdleTime(time.Millisecond)
	r.setFleeTime(time.Millisecond)
	
	fmt.Printf("rabbit is at %s\n", r.Location())
	time.Sleep(time.Duration(2) * time.Millisecond)
	r.DisturbanceAt(dir)
	fmt.Printf("rabbit is at %s\n", r.Location())
}