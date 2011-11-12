package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Set up command-line flags
var waldoDir = flag.String("waldoDir", "", "The directory containing waldo images")
var targetDir = flag.String("targetDir", "", "The directory containing target images")
var numProcs = flag.Int("numProcs", 16, "The number of processors to use (defaults to 16)")

func main() {
	flag.Parse()
	if len(*waldoDir) == 0 || len(*targetDir) == 0 {
		fmt.Println("You need to specify waldo and target directories!")
		fmt.Println("See", os.Args[0], "--help for more information.")
		return
	}
	fmt.Println("Waldo Dir has value:", *waldoDir)
	fmt.Println("Target Dir has value:", *targetDir)
	fmt.Println("Current number of processors:", runtime.GOMAXPROCS(0))
	fmt.Println("Number of processors requested:", *numProcs)

	runtime.GOMAXPROCS(*numProcs)

	fmt.Println("New number of processors:", runtime.GOMAXPROCS(0))

	// Read Waldo Directory
	waldoImages := TraverseDirectory(*waldoDir, ReadFileContents)
	for _, image := range waldoImages {
		fmt.Println("Image:", image)
	}

	// Read Target Directory
	targetImages := TraverseDirectory(*targetDir, ReadFileContents)
	for _, image := range targetImages {
		fmt.Println("Image:", image)
	}
}
