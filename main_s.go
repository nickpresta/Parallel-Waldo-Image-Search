package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Set up command-line flags
var waldoDir = flag.String("waldoDir", "", "The directory containing waldo images")
var targetDir = flag.String("targetDir", "", "The directory containing target images")

func main() {
	startTime := time.Nanoseconds()
	flag.Parse()
	if *waldoDir == "" || *targetDir == "" {
		fmt.Println("You need to specify waldo and target directories!")
		fmt.Println("See", os.Args[0], "--help for more information.")
		return
	}

	runtime.GOMAXPROCS(1)

	// Read Waldo Directory
	waldoImages, err := ReadDirectory(*waldoDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Read Target Directory
	targetImages, err := ReadDirectory(*targetDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(targetImages); i++ {
		targetImages[i].FindImages(waldoImages)
	}

	fmt.Printf("Completed in %f seconds!\n", float64(time.Nanoseconds() - startTime) / 1000000000.0)
}
