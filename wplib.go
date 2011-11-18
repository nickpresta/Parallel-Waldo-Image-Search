package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type waldoImage struct {
	height   int
	width    int
	rotation int
	fileName string
	data     []string
}

type Image interface {
	Rotate() Image
}

// Rotate rotates a waldoImage by 90 degrees to the right
// Returns a new waldoImage
func (this *waldoImage) Rotate() (img waldoImage) {
	img.height = this.width
	img.width = this.height
	img.fileName = this.fileName
	img.data = make([]string, img.height)
	for i := 0; i < img.height; i++ {
		var line []uint8 = make([]uint8, img.width)
		for j := 0; j < img.width; j++ {
			line[j] = this.data[img.width - j - 1][i]
		}
		img.data[i] = string(line)
	}
	return
}

func Read(file *os.File) (img *waldoImage) {
	reader, err := bufio.NewReaderSize(file, 6*1024)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	img = new(waldoImage)

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

func ReadFile(filePath string) *waldoImage {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	return Read(file)
}

func ReadDirectory(directory string) (images []*waldoImage) {
	dirContents, err := os.Open(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer dirContents.Close()

	file, err := dirContents.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Channel to pass along waldoImages (used in reading file contents)
	ch := make(chan *waldoImage)

	var numImages int
	for index, file := range file {
		if file.IsRegular() {
			// Create named function for goroutine
			processFile := func (file os.FileInfo) {
				path, _ := filepath.Abs(filepath.Join(directory, file.Name))
				image := ReadFile(path)
				ch <- image
			}
			go processFile(file)
		}
		numImages = index
	}

	// Collect images and store them
	for i := 0; i <= numImages; i++ {
		image := <-ch
		if image != nil {
			images = append(images, image)
		}
	}

	return
}
