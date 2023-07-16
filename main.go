package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

const (
	outputFilename = "spritesheet.png"
)

func main() {
	// Get the folder path from the command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide the folder path as an argument.")
		return
	}

	folderPath := os.Args[1]

	// Scan the folder for PNG files
	filePaths, err := scanFolder(folderPath)
	if err != nil {
		fmt.Printf("Error scanning folder: %v\n", err)
		return
	}

	// Sort the file paths alphabetically
	sort.Strings(filePaths)

	// Open the images to get the dimensions
	images := make([]image.Image, len(filePaths))
	maxWidth := 0

	for i, filePath := range filePaths {
		img, err := openImage(filePath)
		if err != nil {
			fmt.Printf("Error opening image: %v\n", err)
			return
		}

		images[i] = img

		width := img.Bounds().Dx()
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Create the spritesheet canvas
	spritesheet := image.NewRGBA(image.Rect(0, 0, maxWidth, images[0].Bounds().Dy()*len(images)))

	// Append each image vertically to the spritesheet
	y := 0
	for _, img := range images {
		// Draw the image on the spritesheet
		rect := image.Rect(0, y, img.Bounds().Dx(), y+img.Bounds().Dy())
		draw.Draw(spritesheet, rect, img, image.Point{0, 0}, draw.Src)

		// Increment the y position for the next image
		y += img.Bounds().Dy()
	}

	// Save the spritesheet as a PNG file in the input folder
	outputPath := filepath.Join(folderPath, outputFilename)
	err = saveImage(spritesheet, outputPath)
	if err != nil {
		fmt.Printf("Error saving spritesheet: %v\n", err)
		return
	}

	fmt.Printf("Spritesheet created successfully: %s\n", outputPath)
}

// scanFolder scans the given folder path for PNG files and returns their file paths.
func scanFolder(folderPath string) ([]string, error) {
	var filePaths []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			filePaths = append(filePaths, filepath.Join(folderPath, file.Name()))
		}
	}

	return filePaths, nil
}

// openImage opens the image file at the given path and returns the decoded image.
func openImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// saveImage saves the given image as a PNG file with the specified filename.
func saveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
