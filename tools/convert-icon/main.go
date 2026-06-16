package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/fyne-io/image/ico"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./tools/convert-icon <input.png> <output.ico>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open and decode PNG file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding PNG: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Encode as ICO
	err = ico.Encode(outFile, img)
	if err != nil {
		fmt.Printf("Error encoding ICO: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

