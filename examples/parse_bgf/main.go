package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unger/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_bgf <filename.bgf>")
		fmt.Println("Example: parse_bgf matches/game1.bgf")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse the BGF file
	match, err := bgfparser.ParseBGF(filename)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	// Display match information
	fmt.Println("=== BGF Match File ===")
	fmt.Println(match.String())

	fmt.Println("\n=== Header Information ===")
	fmt.Printf("Format: %s\n", match.Format)
	fmt.Printf("Version: %s\n", match.Version)
	fmt.Printf("Compressed: %v\n", match.Compress)
	fmt.Printf("Uses SMILE encoding: %v\n", match.UseSmile)

	// Display any decoding warnings
	if match.DecodingWarning != "" {
		fmt.Printf("\n⚠️  Warning: %s\n", match.DecodingWarning)
		fmt.Println("\nNote: Full SMILE decoding requires additional libraries.")
		fmt.Println("However, we've extracted some information from the binary data:")
	}

	// Display extracted information
	if len(match.Data) > 0 {
		fmt.Println("\n=== Match Data ===")
		info := match.GetMatchInfo()
		for key, value := range info {
			if key != "format" && key != "version" && key != "compress" && key != "useSmile" {
				fmt.Printf("%s: %v\n", key, value)
			}
		}

		// If SMILE encoded, show extracted strings
		if match.UseSmile {
			if strings, ok := match.Data["_extractedStrings"].([]string); ok && len(strings) > 0 {
				fmt.Println("\n=== Extracted Strings (from binary data) ===")
				for i, s := range strings {
					if len(s) > 50 {
						fmt.Printf("%d: %s...\n", i+1, s[:50])
					} else {
						fmt.Printf("%d: %s\n", i+1, s)
					}
				}
			}
			if size, ok := match.Data["_dataSize"].(int); ok {
				fmt.Printf("\nBinary data size: %d bytes\n", size)
			}
		}
	}
}
