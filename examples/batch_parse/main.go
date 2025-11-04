package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevung/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: batch_parse <directory>")
		fmt.Println("Example: batch_parse ./tmp")
		os.Exit(1)
	}

	directory := os.Args[1]

	// Walk through directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		switch ext {
		case ".txt":
			parseTXTFile(path)
		case ".bgf":
			parseBGFFile(path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}
}

func parseTXTFile(path string) {
	fmt.Printf("\n=== Parsing TXT: %s ===\n", filepath.Base(path))

	position, err := bgfparser.ParseTXT(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Players: %s vs %s\n", position.PlayerX, position.PlayerO)
	fmt.Printf("Score: %d-%d in a %d point match\n", position.ScoreO, position.ScoreX, position.MatchLength)
	fmt.Printf("On Roll: %s", position.OnRoll)
	if position.Dice[0] > 0 {
		fmt.Printf(" with %d-%d", position.Dice[0], position.Dice[1])
	}
	fmt.Println()

	if len(position.Evaluations) > 0 {
		fmt.Printf("Evaluations: %d moves analyzed\n", len(position.Evaluations))
		if position.Evaluations[0].IsBest {
			fmt.Printf("Best move: %s (%.3f)\n", position.Evaluations[0].Move, position.Evaluations[0].Equity)
		}
	}

	if position.CubeDecision != nil {
		fmt.Printf("Cube Decision: %s\n", position.CubeDecision.Action)
	}
}

func parseBGFFile(path string) {
	fmt.Printf("\n=== Parsing BGF: %s ===\n", filepath.Base(path))

	match, err := bgfparser.ParseBGF(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Format: %s v%s\n", match.Format, match.Version)
	fmt.Printf("Compressed: %v, SMILE: %v\n", match.Compress, match.UseSmile)

	if match.UseSmile {
		fmt.Println("Note: SMILE encoding detected (binary JSON format)")
	}

	if match.Data != nil {
		fmt.Printf("Data fields: %d\n", len(match.Data))
	}
}
