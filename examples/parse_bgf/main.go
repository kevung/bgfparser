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

	if match.UseSmile {
		fmt.Println("\nNote: SMILE encoding is detected but not yet fully supported.")
		fmt.Println("The file contains binary-encoded JSON that requires a SMILE decoder.")
		fmt.Println("Consider using a SMILE library like github.com/stdiopt/smile for full parsing.")
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
	}
}
