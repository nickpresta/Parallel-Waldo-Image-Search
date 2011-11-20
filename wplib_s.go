package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"./kmp"
)

type Image struct {
	height   int
	width    int
	rotation int
	fileName string
	data     []string
}

type Image2D interface {
	Rotate() *Image2D
	findImage(image *Image) bool
	findImages(images []*Image)
}

// Rotate rotates a Image by 90 degrees to the right
// Returns a new Image
func (this *Image) Rotate() *Image {
	var img = new(Image)
	img.height = this.width
	img.width = this.height
	img.fileName = this.fileName
	img.data = make([]string, img.height)

	for i := 0; i < img.height; i++ {
			var line []uint8 = make([]uint8, img.width)
			for j := 0; j < img.width; j++ {
				line[j] = this.data[img.width-j-1][i]
			}
			img.data[i] = string(line)
	}

	return img
}

func Read(file *os.File) (img *Image) {
	reader, err := bufio.NewReaderSize(file, 6*1024)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	img = new(Image)

	img.fileName = filepath.Base(file.Name())

	line, isPrefix, err := reader.ReadLine()
	// Get first line, for the dimensions
	dimensions := strings.Split(string(line), " ")
	img.height, _ = strconv.Atoi(dimensions[0])
	img.width, _ = strconv.Atoi(dimensions[1])

	line, isPrefix, err = reader.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		img.data = append(img.data, s)
		line, isPrefix, err = reader.ReadLine()
	}
	if isPrefix {
		fmt.Println("Buffer was declared to be too small for file (", file.Name(), ")")
		return nil
	}
	if err != os.EOF {
		fmt.Println(err)
		return nil
	}

	return
}

func ReadFile(filePath string) *Image {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	return Read(file)
}

func ReadDirectory(directory string) (images []*Image, err os.Error) {
	dirContents, err := os.Open(directory)
	if err != nil {
		return images, err
	}
	defer dirContents.Close()

	file, err := dirContents.Readdir(-1)
	if err != nil {
		return images, err
	}

	for _, file := range file {
		if file.IsRegular() {
			// Create named function for goroutine
			path, _ := filepath.Abs(filepath.Join(directory, file.Name))
			image := ReadFile(path)
			if image != nil {
				images = append(images, image)
			}
		}
	}

	return
}

func (this *Image) FindImages(images []*Image) {
	// this is the target image
	rotations := []int{0, 90, 180, 270}
	// For each waldo
	for i := 0; i < len(images); i++ {
		waldo := images[i]
		// For each rotation
		for _, rotation := range rotations {
			waldo.rotation = rotation
			found := this.FindImage(waldo)
			if found {
				break
			} else {
				waldo = waldo.Rotate()
			}
		}
	}
}

func (this *Image) FindImage(image *Image) bool {
	needle, _ := kmp.NewKMP(image.data[0])
	for i := 0; i < this.height; i++ {
		haystack := this.data[i]
		foundCols := needle.FindAllStringIndex(haystack)
		for _, foundCol := range foundCols {
			// Start descent through image
			numRows := 1
			found := false
			for j := 1; j < image.height && j < this.height &&
				j+i < this.height; j, numRows = j+1, numRows+1 {
				// Check each starting position for the substring match
				substr := this.data[i+j][foundCol : foundCol+image.width]
				if substr == image.data[j] {
					found = true
				} else {
					break
				}
			}
			if found && numRows == image.height {
				y, x := formatCoords(i, foundCol, image)
				fmt.Printf("$%s %s (%d,%d,%d)\n", image.fileName, this.fileName, y, x, image.rotation)
				return true
			}
		}
	}
	return false
}

func formatCoords(y int, x int, image *Image) (int, int) {
	switch image.rotation {
	case 90:
		return y, x + image.width
	case 180:
		return y + image.height, x + image.width
	case 270:
		return y + image.height, x
	default: // For a 0deg rotation
		return y + 1, x + 1
	}
	return y, x
}
