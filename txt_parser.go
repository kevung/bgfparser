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
	originalLine := line
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "=") {
		return nil
	}

	// Skip lines that are just probabilities (second line of each evaluation)
	// These lines start with a decimal number (e.g., "0.254  0.000...")
	// after trimming, not with a rank marker (e.g., "1." or "1)")
	// Check the original untrimmed line for the rank marker
	if !regexp.MustCompile(`^\s*\d+[.)]`).MatchString(originalLine) {
		// No rank marker at start of original line, so this is not an evaluation line
		return nil
	}

	// Also skip if the trimmed line starts with a decimal number
	// (probability lines like "0.254  0.000  0.000  -  0.746...")
	if regexp.MustCompile(`^\d+\.\d+\s`).MatchString(line) {
		return nil
	}

	eval := &Evaluation{}

	// Check if this is marked as best move
	if strings.Contains(line, "*") {
		eval.IsBest = true
		line = strings.ReplaceAll(line, "*", "")
	}

	// Parse rank number at start - support both formats: "1)" and "1."
	// Format 1: "1) 13-11 24-23                0.473 / -0.289"
	// Format 2: "1.   0.124 mwp /  -0.492            19/18, 14/12"
	re := regexp.MustCompile(`^(\d+)[.)]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 2 {
		rankNum, _ := strconv.Atoi(matches[1])
		eval.Rank = rankNum
		*rank = rankNum
		line = line[len(matches[0]):]
	} else {
		return nil
	}

	// Trim whitespace after rank
	line = strings.TrimSpace(line)

	// Parse the rest of the line
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	// Check for "mwp" format (new style)
	// Format: "0.124 mwp /  -0.492            19/18, 14/12"
	if len(parts) >= 2 && parts[1] == "mwp" {
		// New format with mwp
		eval.Equity, _ = strconv.ParseFloat(parts[0], 64)

		// Find where the move starts (after the "/" and equity values)
		moveStartIdx := -1
		for i := 0; i < len(parts); i++ {
			if parts[i] == "/" && i+1 < len(parts) {
				// The next element should be the equity value (negative)
				// Move starts after the equity and optional diff
				moveStartIdx = i + 2 // Skip "/" and equity
				// Check if there's a diff in parentheses
				if moveStartIdx < len(parts) && strings.HasPrefix(parts[moveStartIdx], "(") {
					// Parse diff
					diffStr := strings.Trim(parts[moveStartIdx], "()")
					eval.Diff, _ = strconv.ParseFloat(diffStr, 64)
					moveStartIdx++
				}
				break
			}
		}

		// Extract move
		if moveStartIdx > 0 && moveStartIdx < len(parts) {
			eval.Move = strings.Join(parts[moveStartIdx:], " ")
		}
	} else {
		// Old format without mwp
		// Format: "13-11 24-23                0.473 / -0.289"
		// Find the "/" separator
		slashIdx := -1
		for i, part := range parts {
			if part == "/" {
				slashIdx = i
				break
			}
		}

		if slashIdx > 0 {
			// Everything before "/" is the move
			eval.Move = strings.Join(parts[:slashIdx-1], " ")

			// Parse equity (before "/")
			if slashIdx > 0 {
				eval.Equity, _ = strconv.ParseFloat(parts[slashIdx-1], 64)
			}

			// Parse diff (after "/")
			if slashIdx+1 < len(parts) {
				diffStr := strings.Trim(parts[slashIdx+1], "()")
				eval.Diff, _ = strconv.ParseFloat(diffStr, 64)
			}
		}
	}

	return eval
}

// parseProbabilityLine parses the probability detail line that follows an evaluation
// Format: "   0.443  0.113  0.002  -  0.557  0.179  0.003"
// Which represents: Win WinG WinBG - (Lose implied) LoseG LoseBG
func parseProbabilityLine(line string, eval *Evaluation) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}

	// Check if this looks like a probability line
	// Should start with a decimal number and contain a dash separator
	if !regexp.MustCompile(`^\d+\.\d+\s`).MatchString(line) {
		return false
	}

	if !strings.Contains(line, "-") {
		return false
	}

	parts := strings.Fields(line)
	if len(parts) < 7 {
		return false
	}

	// Find the dash separator
	dashIdx := -1
	for i, part := range parts {
		if part == "-" {
			dashIdx = i
			break
		}
	}

	if dashIdx < 3 || dashIdx+3 >= len(parts) {
		return false
	}

	// Parse win probabilities (before dash)
	eval.Win, _ = strconv.ParseFloat(parts[0], 64)
	eval.WinG, _ = strconv.ParseFloat(parts[1], 64)
	eval.WinBG, _ = strconv.ParseFloat(parts[2], 64)

	// Parse lose probabilities (after dash)
	// Note: parts[dashIdx+1] is the lose probability (1 - win), we skip it
	eval.LoseG, _ = strconv.ParseFloat(parts[dashIdx+2], 64)
	eval.LoseBG, _ = strconv.ParseFloat(parts[dashIdx+3], 64)

	return true
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
