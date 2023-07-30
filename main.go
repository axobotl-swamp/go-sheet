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

type Mode int

const (
	None Mode = iota
	Action
	Individual
)

func main() {
	// Get the folder path and mode from the command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Please provide the folder path and mode as arguments.")
		return
	}

	modeStr := os.Args[1]
	folderPath := os.Args[2]
	mode := None
	if modeStr == "a" || modeStr == "action" {
		mode = Action
	}
	if modeStr == "i" || modeStr == "individual" {
		mode = Individual
	}
	switch mode {
	case Action:
		processActionMode(folderPath)
	case Individual:
		processIndividualMode(folderPath)
	default:
		fmt.Println("Invalid mode. Available modes: action, individual")
		return
	}

	switch modeStr {
	case "a":
	case "action":
		processActionMode(folderPath)
	case "i":
	case "individual":
		processIndividualMode(folderPath)
	default:
		fmt.Println("Invalid mode. Available modes: action, individual")
		return
	}
}

func processActionMode(folderPath string) {
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
	fmt.Printf("Rows: %d\n", len(filePaths))
}

func processIndividualMode(folderPath string) {
	// Get the list of subfolders in the main folder
	subfolders, err := getSubfolders(folderPath)
	if err != nil {
		fmt.Printf("Error getting subfolders: %v\n", err)
		return
	}

	// Sort subfolders alphabetically
	sort.Strings(subfolders)

	// Calculate total width and height for the spritesheet
	maxWidth := 0
	maxColumns := 0
	maxHeight := 0
	for _, subfolder := range subfolders {
		// Scan the subfolder for PNG files
		subfolderPath := filepath.Join(folderPath, subfolder)
		filePaths, err := scanFolder(subfolderPath)
		if err != nil {
			fmt.Printf("Error scanning subfolder %s: %v\n", subfolder, err)
			continue
		}

		// Sort the file paths alphabetically
		sort.Strings(filePaths)

		// Open the images to get the dimensions
		images := make([]image.Image, len(filePaths))

		columns := len(filePaths)
		if columns > maxColumns {
			maxColumns = columns
		}
		for i, filePath := range filePaths {
			img, err := openImage(filePath)
			if err != nil {
				fmt.Printf("Error opening image: %v\n", err)
				return
			}

			images[i] = img

			height := img.Bounds().Dy()
			if height > maxHeight {
				maxHeight = height
			}

			width := img.Bounds().Dx()
			if width > maxWidth {
				maxWidth = height
			}
		}
	}

	// Create the spritesheet canvas
	spritesheet := image.NewRGBA(image.Rect(0, 0, maxHeight*maxColumns, maxHeight*len(subfolders)))

	// Append each folder horizontally to the spritesheet
	y := 0
	for _, subfolder := range subfolders {
		// Scan the subfolder for PNG files
		subfolderPath := filepath.Join(folderPath, subfolder)
		filePaths, err := scanFolder(subfolderPath)
		if err != nil {
			fmt.Printf("Error scanning subfolder %s: %v\n", subfolder, err)
			continue
		}

		// Sort the file paths alphabetically
		sort.Strings(filePaths)

		x := 0
		for _, filePath := range filePaths {
			img, err := openImage(filePath)
			if err != nil {
				fmt.Printf("Error opening image: %v\n", err)
				return
			}

			// Draw the image on the spritesheet
			rect := image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy())
			draw.Draw(spritesheet, rect, img, image.Point{0, 0}, draw.Src)

			// Increment the x position for the next image
			x += img.Bounds().Dx()
		}

		// Increment the y position for the next folder
		y += maxHeight
	}

	// Save the spritesheet as a PNG file in the input folder
	outputPath := filepath.Join(folderPath, outputFilename)
	err = saveImage(spritesheet, outputPath)
	if err != nil {
		fmt.Printf("Error saving spritesheet: %v\n", err)
		return
	}

	fmt.Printf("Spritesheet created successfully: %s\n", outputPath)
	fmt.Printf("Columns: %d | Rows: %d\n", maxColumns, len(subfolders))
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

// getSubfolders returns a list of subfolders in the given folder path.
func getSubfolders(folderPath string) ([]string, error) {
	var subfolders []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subfolders = append(subfolders, file.Name())
		}
	}

	return subfolders, nil
}
