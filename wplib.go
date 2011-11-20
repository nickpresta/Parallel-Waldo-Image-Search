package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type imageLine struct {
	position int
	line     string
}

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
	ch := make(chan imageLine, img.height)
	for i := 0; i < img.height; i++ {
		rotateLine := func(lineNo int, imageWidth int, data []string, outchan chan imageLine) {
			var line []uint8 = make([]uint8, imageWidth)
			for j := 0; j < imageWidth; j++ {
				line[j] = data[imageWidth-j-1][lineNo]
			}
			var out imageLine
			out.position = lineNo
			out.line = string(line)
			outchan <- out
		}
		go rotateLine(i, img.width, this.data, ch)
	}

	// put lines where they belong
	for i := 0; i < img.height; i++ {
		data := <-ch
		img.data[data.position] = data.line
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

	// Channel to pass along Images (used in reading file contents)
	ch := make(chan *Image)

	var numImage2Ds int
	for index, file := range file {
		if file.IsRegular() {
			// Create named function for goroutine
			processFile := func(file os.FileInfo) {
				path, _ := filepath.Abs(filepath.Join(directory, file.Name))
				image := ReadFile(path)
				ch <- image
			}
			go processFile(file)
		}
		numImage2Ds = index
	}

	// Collect images and store them
	for i := 0; i <= numImage2Ds; i++ {
		image := <-ch
		if image != nil {
			images = append(images, image)
		}
	}

	return
}

func (this *Image) FindImages(images []*Image, done chan bool) {
	// this is the target image
	rotations := []int{0, 90, 180, 270}
	ch := make(chan bool, len(images))
	// For each waldo
	for i := 0; i < len(images); i++ {
		searchImage := func(index int, images []*Image, ch chan bool) {
			waldo := images[index]
			// For each rotation
			for j := 0; j < 4; j++ {
				waldo.rotation = rotations[j]
				found := this.FindImage(waldo)
				if found {
					break
				} else {
					waldo = waldo.Rotate()
				}
			}
			ch <- true
		}
		go searchImage(i, images, ch)
	}

	for i := 0; i < len(images); i++ {
		<-ch
	}

	done <- true
}

func (this *Image) FindImage(image *Image) bool {
	needle := image.data[0]
	for i := 0; i < this.height; i++ {
		haystack := this.data[i]
		foundCol := strings.Index(haystack, needle)
		for foundCol != -1 {
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
			// Move forward past first found waldo
			// Find waldo in this slice, add what we skipped over previously
			previous := foundCol
			foundCol = strings.Index(haystack[previous+image.width:], needle)
			if foundCol >= 0 {
				foundCol = foundCol + previous + image.width
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
