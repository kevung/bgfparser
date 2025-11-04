package bgfparser

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseTXT parses a BGBlitz position text file from disk
// This is a convenience wrapper around ParseTXTFromReader that handles file reading.
func ParseTXT(filename string) (*Position, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, &ParseError{File: filename, Message: err.Error()}
	}
	defer file.Close()

	pos, err := ParseTXTFromReader(file)
	if err != nil {
		// Add filename to error if not already present
		if parseErr, ok := err.(*ParseError); ok && parseErr.File == "" {
			parseErr.File = filename
			return nil, parseErr
		}
		return nil, err
	}

	return pos, nil
}

// parseBoard extracts checker positions from board lines
func parseBoard(pos *Position, lines []string) {
	// Initialize board
	for i := range pos.Board {
		pos.Board[i] = 0
	}

	// Simple board parsing - count X and O on each point
	// This is a simplified version; full implementation would parse the exact positions
	for _, line := range lines {
		// Extract point numbers and checkers
		// This is complex due to the ASCII art format
		// For now, we'll extract basic information
		if strings.Contains(line, "BAR") && (strings.Contains(line, "X") || strings.Contains(line, "O")) {
			// Check for checkers on bar
			parts := strings.Split(line, "BAR")
			if len(parts) > 0 {
				if strings.Count(parts[0], "X") > strings.Count(parts[0], "O") {
					pos.OnBar["X"]++
				} else if strings.Count(parts[0], "O") > strings.Count(parts[0], "X") {
					pos.OnBar["O"]++
				}
			}
		}
	}
}

// parseXGID extracts information from XGID format
func parseXGID(pos *Position, xgid string) {
	// XGID format: board:cubeValue:cubeOwner:onRoll:dice:crawford:score1:score2:matchLength:turn
	parts := strings.Split(xgid, ":")
	if len(parts) >= 5 {
		// Parse cube value
		if val, err := strconv.Atoi(parts[1]); err == nil {
			pos.CubeValue = 1 << val // Cube value is 2^n
		}
		// Parse cube owner
		switch parts[2] {
		case "1":
			pos.CubeOwner = "X"
		case "-1":
			pos.CubeOwner = "O"
		default:
			pos.CubeOwner = ""
		}
		// Parse on roll
		switch parts[3] {
		case "1":
			pos.OnRoll = "X"
		case "-1":
			pos.OnRoll = "O"
		}
	}
}

// parseEvaluation parses a single evaluation line
func parseEvaluation(line string, rank *int) *Evaluation {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "=") {
		return nil
	}

	eval := &Evaluation{}

	// Check if this is marked as best move
	if strings.Contains(line, "*") {
		eval.IsBest = true
		line = strings.ReplaceAll(line, "*", "")
	}

	// Parse rank number at start
	re := regexp.MustCompile(`^(\d+)\)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 2 {
		*rank++
		eval.Rank = *rank
		line = line[len(matches[0]):]
	} else {
		return nil
	}

	// Split into move and numbers
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	// Extract move (everything before the equity value)
	moveEnd := 0
	for i, part := range parts {
		// Look for the equity value (format: 0.123 or -0.123)
		if matched, _ := regexp.MatchString(`^-?\d+\.\d+$`, part); matched {
			moveEnd = i
			break
		}
	}

	if moveEnd > 0 {
		eval.Move = strings.Join(parts[:moveEnd], " ")

		// Parse equity and diff
		if moveEnd < len(parts) {
			eval.Equity, _ = strconv.ParseFloat(parts[moveEnd], 64)
		}
		if moveEnd+2 < len(parts) {
			diffStr := strings.Trim(parts[moveEnd+2], "()")
			eval.Diff, _ = strconv.ParseFloat(diffStr, 64)
		}

		// Parse win percentages (next line typically)
		// Format: 0.443  0.113  0.002  -  0.557  0.179  0.003
		// This would need to be handled in the main parsing loop
	}

	return eval
}

// parseCubeDecision parses a cube decision line
func parseCubeDecision(line string) *CubeDecision {
	decision := &CubeDecision{}

	if strings.Contains(line, "*") {
		decision.IsBest = true
		line = strings.ReplaceAll(line, "*", "")
	}

	// Extract action
	if strings.Contains(line, "Double") && strings.Contains(line, "Take") {
		decision.Action = "Double/Take"
	} else if strings.Contains(line, "Double") && strings.Contains(line, "Pass") {
		decision.Action = "Double/Pass"
	} else if strings.Contains(line, "No Double") {
		decision.Action = "No Double"
	}

	// Parse MWC and EMG values
	re := regexp.MustCompile(`(\d+\.\d+)`)
	matches := re.FindAllString(line, -1)

	if len(matches) >= 2 {
		decision.MWC, _ = strconv.ParseFloat(matches[0], 64)
		decision.EMG, _ = strconv.ParseFloat(matches[1], 64)
	}

	// Parse differences in parentheses
	reDiff := regexp.MustCompile(`\(([+-]?\d+\.\d+)\)`)
	diffMatches := reDiff.FindAllStringSubmatch(line, -1)
	if len(diffMatches) >= 1 {
		decision.MWCDiff, _ = strconv.ParseFloat(diffMatches[0][1], 64)
	}
	if len(diffMatches) >= 2 {
		decision.EMGDiff, _ = strconv.ParseFloat(diffMatches[1][1], 64)
	}

	return decision
}
