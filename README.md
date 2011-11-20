# Parallel Waldo Image Search

This code performs a simple form of 2D template matching using strings.

The source image, waldo, is represented as a 2D image of 1s and 0s.
The target image, is represented as a 2D image, equal to, or larger than, the waldo image. The target image contains zero or more waldo images (but only 1 unique waldo, per target) in any rotation of 0, 90, 180, or 270 degrees.

This code is attempted in parallel using Go's channel message passing architecture.

## Input

Input is expected in the form of two directories - waldoDir, and targetDir - containing waldo images and target images respectively.

From [Effective Go](http://golang.org/doc/effective_go.html#concurrency):

> The current implementation of gc (6g, etc.) will not parallelize this code by default.
> It dedicates only a single core to user-level processing.
> An arbitrary number of goroutines can be blocked in system calls, but by default only one can be executing user-level code at any time.
> It should be smarter and one day it will be smarter, but until it is if you want CPU parallelism you must tell the run-time how many goroutines you want executing code simultaneously.
> There are two related ways to do this. Either run your job with environment variable GOMAXPROCS set to the number of cores to use (default 1); or import the runtime package and call runtime.GOMAXPROCS(NCPU).
> Again, this requirement is expected to be retired as the scheduling and run-time improve.

As such, this code also takes in a *numProcs* command-line option. This sets GOMAXPROCS at runtime.

## Output

Each match will be printed to STDOUT, in the form:

    $parfile imfile (y,x,r)

Where parfile and imfile are the bare filenames (no directory path) of the parallaldo and its matching image, respectively, and (y,x,r) are bracketed integers printed without extra spaces or leading zeroes.

For (y,x,r) where (y,x) denotes the location in the image of the (unrotated) template's upper left corner, and r being the degrees of clockwise rotation. The convention of the first coordinate being the row (here numbered starting from 1) is normal for graphics and matrix data.

# How to Run

## Compile
Since this includes the file *kmp.go*, you need to compile it first: `6g kmp.go`.

The rest of this code uses the included tool `gomake`. To compile:

    gomake

This will create a binary named wp.

## To Run
To run `wp`:

    wp -waldoDir=/path/to/waldo -targetDir=/path/to/target -numProcs=8

This will spit the results to stdout.
