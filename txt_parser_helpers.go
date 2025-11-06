package bgfparser

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

// parseBoardLine checks if a line is part of the board display
func parseBoardLine(line string, boardLines *[]string) bool {
	if !strings.Contains(line, "|") {
		return false
	}

	if strings.Contains(line, "+") {
		// Board boundary lines - skip
		return true
	}

	if strings.Contains(line, "BAR") || strings.Contains(line, "X") || strings.Contains(line, "O") {
		*boardLines = append(*boardLines, line)
		return true
	}

	return false
}

// parsePlayerInfo extracts player names and pip counts
func parsePlayerInfo(line string, pos *Position) {
	// Look for either "O:" or "X:" in the line
	if !strings.Contains(line, "O:") && !strings.Contains(line, "X:") {
		return
	}

	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "O:" && i+1 < len(parts) {
			pos.PlayerO = parts[i+1]
			if i+2 < len(parts) {
				if score, err := strconv.Atoi(parts[i+2]); err == nil {
					pos.PipCount["O"] = score
				}
			}
		}
		if part == "X:" && i+1 < len(parts) {
			pos.PlayerX = parts[i+1]
			if i+2 < len(parts) {
				if score, err := strconv.Atoi(parts[i+2]); err == nil {
					pos.PipCount["X"] = score
				}
			}
		}
	}
}

// parsePositionID extracts Position-ID and Match-ID
func parsePositionID(line string, pos *Position) {
	if !strings.Contains(line, "Position-ID:") {
		return
	}

	re := regexp.MustCompile(`Position-ID:\s*(\S+)\s+Match-ID:\s*(\S+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		pos.PositionID = matches[1]
		pos.MatchID = matches[2]
	}
}

// parseXGIDLine extracts and parses XGID
func parseXGIDLine(line string, pos *Position) {
	if !strings.Contains(line, "XGID=") {
		return
	}

	re := regexp.MustCompile(`XGID=(\S+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 2 {
		pos.XGID = matches[1]
		parseXGID(pos, matches[1])
	}
}

// parseMatchScore extracts match length and scores
func parseMatchScore(line string, pos *Position) {
	if !strings.Contains(line, "point match") {
		return
	}

	re := regexp.MustCompile(`(\S+)\s*-\s*(\d+)\s+(\S+)\s*-\s*(\d+)\s+in a\s+(\d+)\s+point match`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 6 {
		pos.ScoreO, _ = strconv.Atoi(matches[2])
		pos.ScoreX, _ = strconv.Atoi(matches[4])
		pos.MatchLength, _ = strconv.Atoi(matches[5])
	}
}

// parseCurrentPlayer extracts current player and dice
func parseCurrentPlayer(line string, pos *Position) {
	if !strings.Contains(line, "to move") {
		return
	}

	if strings.Contains(line, pos.PlayerX) {
		pos.OnRoll = "X"
	} else if strings.Contains(line, pos.PlayerO) {
		pos.OnRoll = "O"
	}

	// Parse dice
	re := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		pos.Dice[0], _ = strconv.Atoi(matches[1])
		pos.Dice[1], _ = strconv.Atoi(matches[2])
	}
}

// parseCubeValue extracts cube value from display
func parseCubeValue(line string, scanner *bufio.Scanner, pos *Position) bool {
	if !strings.Contains(line, "+--+") {
		return false
	}

	if !scanner.Scan() {
		return true
	}

	cubeLine := scanner.Text()
	if !strings.Contains(cubeLine, "|") {
		return true
	}

	re := regexp.MustCompile(`\|\s*(\d+)\s*\|`)
	matches := re.FindStringSubmatch(cubeLine)
	if len(matches) == 2 {
		pos.CubeValue, _ = strconv.Atoi(matches[1])
	}

	return true
}

// handleEvaluationSection manages evaluation and cube decision section state
func handleEvaluationSection(line string, inEvaluation, inCubeDecision *bool, evalRank *int) bool {
	// Detect evaluation section - support multiple languages
	// English: "Evaluation", French: "Évaluation", German: "Bewertung", Japanese: "評価"
	if strings.Contains(line, "Evaluation") ||
		strings.Contains(line, "Évaluation") ||
		strings.Contains(line, "Bewertung") ||
		strings.Contains(line, "評価") {
		*inEvaluation = true
		*inCubeDecision = false
		*evalRank = 0
		return true
	}

	// Skip separator lines
	if *inEvaluation && strings.TrimSpace(line) == "==========" {
		return true
	}

	// Detect cube decision section - look for multilingual "Cube Action:" headers or MWC/EMG column headers
	// English: "Cube Action:", German: "Würfelaktion:", French: "Videau:", Japanese: "キューブアクション:"
	if strings.Contains(line, "Cube Action") ||
		strings.Contains(line, "Würfelaktion") ||
		strings.Contains(line, "Videau") ||
		strings.Contains(line, "キューブアクション") ||
		(strings.Contains(line, "MWC") && strings.Contains(line, "EMG")) {
		*inCubeDecision = true
		*inEvaluation = false
		return true
	}

	return false
}
