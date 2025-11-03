package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	bgfparser "github.com/kevung/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: decode_to_json <bgf_file> [output.json]")
		os.Exit(1)
	}

	filename := os.Args[1]
	outputFile := "output.json"
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	}

	match, err := bgfparser.ParseBGF(filename)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		if match == nil {
			os.Exit(1)
		}
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(match.Data, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	var writer io.Writer
	if outputFile == "-" {
		writer = os.Stdout
	} else {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		writer = f
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		fmt.Printf("Error writing JSON: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "-" {
		fmt.Printf("Decoded data written to %s\n", outputFile)
		if match.DecodingWarning != "" {
			fmt.Printf("Warning: %s\n", match.DecodingWarning)
		}
	}
}
