package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevung/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_txt <filename.txt>")
		fmt.Println("Example: parse_txt positions/game1.txt")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse the position file
	position, err := bgfparser.ParseTXT(filename)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	// Display position information
	fmt.Println("=== Position Information ===")
	fmt.Printf("Players: %s (X) vs %s (O)\n", position.PlayerX, position.PlayerO)
	fmt.Printf("Score: %s %d - %d %s\n", position.PlayerO, position.ScoreO, position.ScoreX, position.PlayerX)
	fmt.Printf("Match Length: %d points\n", position.MatchLength)
	fmt.Printf("On Roll: %s\n", position.OnRoll)
	if position.Dice[0] > 0 || position.Dice[1] > 0 {
		fmt.Printf("Dice: %d-%d\n", position.Dice[0], position.Dice[1])
	}
	if position.CubeValue > 0 {
		fmt.Printf("Cube: %d", position.CubeValue)
		if position.CubeOwner != "" {
			fmt.Printf(" (owned by %s)", position.CubeOwner)
		}
		fmt.Println()
	}
	fmt.Printf("Pip Count: X=%d, O=%d\n", position.PipCount["X"], position.PipCount["O"])

	fmt.Println("\n=== Identifiers ===")
	fmt.Printf("Position-ID: %s\n", position.PositionID)
	fmt.Printf("Match-ID: %s\n", position.MatchID)
	fmt.Printf("XGID: %s\n", position.XGID)

	// Display evaluations
	if len(position.Evaluations) > 0 {
		fmt.Println("\n=== Move Evaluations ===")
		for _, eval := range position.Evaluations {
			best := ""
			if eval.IsBest {
				best = " *"
			}
			fmt.Printf("%d) %s%s\n", eval.Rank, eval.Move, best)
			fmt.Printf("   Equity: %.3f (%.3f)\n", eval.Equity, eval.Diff)
			if eval.Win > 0 {
				fmt.Printf("   Win: %.3f  WinG: %.3f  WinBG: %.3f\n", eval.Win, eval.WinG, eval.WinBG)
				if eval.LoseG > 0 || eval.LoseBG > 0 {
					fmt.Printf("   Lose: %.3f  LoseG: %.3f  LoseBG: %.3f\n", 1.0-eval.Win, eval.LoseG, eval.LoseBG)
				}
			}
		}
	}

	// Display cube decision
	if position.CubeDecision != nil {
		fmt.Println("\n=== Cube Decision ===")
		cd := position.CubeDecision
		best := ""
		if cd.IsBest {
			best = " *"
		}
		fmt.Printf("%s%s\n", cd.Action, best)
		fmt.Printf("MWC: %.3f", cd.MWC)
		if cd.MWCDiff != 0 {
			fmt.Printf(" (%+.3f)", cd.MWCDiff)
		}
		fmt.Printf("  EMG: %.3f", cd.EMG)
		if cd.EMGDiff != 0 {
			fmt.Printf(" (%+.3f)", cd.EMGDiff)
		}
		fmt.Println()
	}
}
