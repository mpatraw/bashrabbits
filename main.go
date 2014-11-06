
package main

import (
	"flag"
	"fmt"
	"os"
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
}