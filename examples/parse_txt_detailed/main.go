package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevung/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_txt_detailed <filename.txt>")
		fmt.Println("Example: parse_txt_detailed positions/game1.txt")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse the position file
	position, err := bgfparser.ParseTXT(filename)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("COMPLETE POSITION INFORMATION EXTRACTION")
	fmt.Println("========================================")
	fmt.Println()

	// === PLAYER INFORMATION ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ PLAYER INFORMATION")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ Player X: %s\n", position.PlayerX)
	fmt.Printf("│ Player O: %s\n", position.PlayerO)
	fmt.Printf("│ Score X:  %d\n", position.ScoreX)
	fmt.Printf("│ Score O:  %d\n", position.ScoreO)
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === MATCH INFORMATION ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ MATCH INFORMATION")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ Match Length: %d points\n", position.MatchLength)
	fmt.Printf("│ Crawford:     %v\n", position.Crawford)
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === POSITION IDENTIFIERS ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ POSITION IDENTIFIERS")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ Position-ID: %s\n", position.PositionID)
	fmt.Printf("│ Match-ID:    %s\n", position.MatchID)
	fmt.Printf("│ XGID:        %s\n", position.XGID)
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === CURRENT STATE ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ CURRENT POSITION STATE")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ On Roll:      %s\n", position.OnRoll)
	if position.Dice[0] > 0 || position.Dice[1] > 0 {
		fmt.Printf("│ Dice:         %d-%d\n", position.Dice[0], position.Dice[1])
	} else {
		fmt.Printf("│ Dice:         (not rolled)\n")
	}
	fmt.Printf("│ Cube Value:   %d\n", position.CubeValue)
	if position.CubeOwner != "" {
		fmt.Printf("│ Cube Owner:   %s\n", position.CubeOwner)
	} else {
		fmt.Printf("│ Cube Owner:   (centered)\n")
	}
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === PIP COUNTS ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ PIP COUNTS")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ Player X: %d pips\n", position.PipCount["X"])
	fmt.Printf("│ Player O: %d pips\n", position.PipCount["O"])
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === CHECKERS ON BAR ===
	if position.OnBar["X"] > 0 || position.OnBar["O"] > 0 {
		fmt.Println("┌─────────────────────────────────────")
		fmt.Println("│ CHECKERS ON BAR")
		fmt.Println("├─────────────────────────────────────")
		if position.OnBar["X"] > 0 {
			fmt.Printf("│ Player X: %d checker(s)\n", position.OnBar["X"])
		}
		if position.OnBar["O"] > 0 {
			fmt.Printf("│ Player O: %d checker(s)\n", position.OnBar["O"])
		}
		fmt.Println("└─────────────────────────────────────")
		fmt.Println()
	}

	// === BOARD POSITION ===
	fmt.Println("┌─────────────────────────────────────")
	fmt.Println("│ BOARD ARRAY (26 points)")
	fmt.Println("├─────────────────────────────────────")
	fmt.Printf("│ Board: %v\n", position.Board)
	fmt.Println("│ (Positive = X checkers, Negative = O checkers)")
	fmt.Println("└─────────────────────────────────────")
	fmt.Println()

	// === MOVE EVALUATIONS ===
	if len(position.Evaluations) > 0 {
		fmt.Println("┌─────────────────────────────────────────────────────────────────")
		fmt.Println("│ MOVE EVALUATIONS")
		fmt.Println("├─────────────────────────────────────────────────────────────────")
		fmt.Printf("│ Total Evaluations: %d\n", len(position.Evaluations))
		fmt.Println("├─────────────────────────────────────────────────────────────────")

		for i, eval := range position.Evaluations {
			best := ""
			if eval.IsBest {
				best = " ★ BEST"
			}

			fmt.Printf("│\n")
			fmt.Printf("│ #%d: %s%s\n", eval.Rank, eval.Move, best)
			fmt.Printf("│     ├─ Equity:      %.4f\n", eval.Equity)
			fmt.Printf("│     ├─ Difference:  %.4f\n", eval.Diff)

			if eval.Win > 0 {
				fmt.Printf("│     ├─ Win:         %.3f (%.1f%%)\n", eval.Win, eval.Win*100)
				fmt.Printf("│     ├─ Win Gammon:  %.3f (%.1f%%)\n", eval.WinG, eval.WinG*100)
				fmt.Printf("│     ├─ Win BG:      %.3f (%.1f%%)\n", eval.WinBG, eval.WinBG*100)
				fmt.Printf("│     ├─ Lose:        %.3f (%.1f%%)\n", 1.0-eval.Win, (1.0-eval.Win)*100)
				fmt.Printf("│     ├─ Lose Gammon: %.3f (%.1f%%)\n", eval.LoseG, eval.LoseG*100)
				fmt.Printf("│     └─ Lose BG:     %.3f (%.1f%%)\n", eval.LoseBG, eval.LoseBG*100)
			}

			if i < len(position.Evaluations)-1 {
				fmt.Println("│     ─────────────────────────────────────")
			}
		}
		fmt.Println("└─────────────────────────────────────────────────────────────────")
		fmt.Println()
	}

	// === CUBE DECISION ===
	if len(position.CubeDecisions) > 0 {
		fmt.Println("┌─────────────────────────────────────")
		fmt.Println("│ CUBE DECISION ANALYSIS")
		fmt.Println("├─────────────────────────────────────")

		// Show equity information if available
		if position.CubelessEquity != 0 || position.CubefulEquity != 0 {
			if position.CubelessEquity != 0 {
				fmt.Printf("│ Cubeless Equity: %.4f", position.CubelessEquity)
				if position.EquityStdDev != 0 {
					fmt.Printf("  (Std.Dev.: %.3f)", position.EquityStdDev)
				}
				fmt.Println()
			}
			if position.CubefulEquity != 0 {
				fmt.Printf("│ Cubeful Equity:  %.4f\n", position.CubefulEquity)
			}
			fmt.Println("├─────────────────────────────────────")
		}

		for _, cd := range position.CubeDecisions {
			best := ""
			if cd.IsBest {
				best = " ★ BEST"
			}
			fmt.Printf("│ %s%s\n", cd.Action, best)
			fmt.Printf("│   MWC: %.4f", cd.MWC)
			if cd.MWCDiff != 0 {
				fmt.Printf(" (%+.4f)", cd.MWCDiff)
			}
			fmt.Println()
			fmt.Printf("│   EMG: %.4f", cd.EMG)
			if cd.EMGDiff != 0 {
				fmt.Printf(" (%+.4f)", cd.EMGDiff)
			}
			fmt.Println()
		}
		fmt.Println("└─────────────────────────────────────")
		fmt.Println()
	}

	// === SUMMARY ===
	fmt.Println("========================================")
	fmt.Println("EXTRACTION SUMMARY")
	fmt.Println("========================================")
	fmt.Printf("✓ Player names and scores extracted\n")
	fmt.Printf("✓ Match information extracted\n")
	fmt.Printf("✓ Position identifiers extracted\n")
	fmt.Printf("✓ Current state (dice, cube) extracted\n")
	fmt.Printf("✓ Pip counts extracted\n")
	if len(position.Evaluations) > 0 {
		fmt.Printf("✓ %d move evaluations with probabilities extracted\n", len(position.Evaluations))
	}
	if len(position.CubeDecisions) > 0 {
		fmt.Printf("✓ %d cube decision(s) with equity analysis extracted\n", len(position.CubeDecisions))
	}
	fmt.Println("========================================")
}
