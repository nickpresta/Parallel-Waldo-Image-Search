package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// Set up command-line flags
var waldoDir = flag.String("waldoDir", "", "The directory containing waldo images")
var targetDir = flag.String("targetDir", "", "The directory containing target images")
var numProcs = flag.Int("numProcs", 16, "The number of processors to use (defaults to 16)")

// Profiling flags
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *waldoDir == "" || *targetDir == "" {
		fmt.Println("You need to specify waldo and target directories!")
		fmt.Println("See", os.Args[0], "--help for more information.")
		return
	}

	// More Profiling boilerplate
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Waldo Dir has value:", *waldoDir)
	fmt.Println("Target Dir has value:", *targetDir)
	fmt.Println("Current number of processors:", runtime.GOMAXPROCS(0))
	fmt.Println("Number of processors requested:", *numProcs)
	runtime.GOMAXPROCS(*numProcs)
	fmt.Println("New number of processors:", runtime.GOMAXPROCS(0))

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

	fmt.Println("Waldos:", len(waldoImages), "Targets:", len(targetImages))

	// Spawn worker threads with directory data
	// This should be using goroutines
	done := make(chan bool, len(targetImages))
	for i := 0; i < len(targetImages); i++ {
		go targetImages[i].FindImages(waldoImages, done)
	}

	// Drain channel
	for i := 0; i < len(targetImages); i++ {
		<-done
	}
}
