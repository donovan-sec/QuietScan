package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/fyne-io/image/ico"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./tools/convert-icon-multi <input.png> <output.ico>")
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

	// For now, encode the original size
	// Note: The fyne-io/image/ico package encodes single-size icons
	// For multi-size icons, you'd need a different library or tool
	err = ico.Encode(outFile, img)
	if err != nil {
		fmt.Printf("Error encoding ICO: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
	fmt.Printf("Note: This creates a single-size ICO. For best Windows compatibility,\n")
	fmt.Printf("consider using an online tool to create a multi-size ICO with 16x16, 32x32, 48x48, and 256x256 sizes.\n")
}


