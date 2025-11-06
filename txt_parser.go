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
	// Note: Board is already parsed from XGID if available
	// This function could parse the ASCII art board representation
	// but for now we rely on XGID parsing which is more accurate

	// Only try to parse checkers on bar from ASCII art if not already set
	for _, line := range lines {
		if strings.Contains(line, "BAR") && (strings.Contains(line, "X") || strings.Contains(line, "O")) {
			// Check for checkers on bar (if not already parsed from XGID)
			parts := strings.Split(line, "BAR")
			if len(parts) > 0 && pos.OnBar["X"] == 0 && pos.OnBar["O"] == 0 {
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
		// Parse board position from first part
		parseXGIDBoard(pos, parts[0])

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

// parseXGIDBoard decodes the board position from XGID format
// XGID board encoding format (26 characters):
//
//	Character 0: X's checkers on bar
//	Characters 1-24: Points 1, 2, 3, ..., 23, 24 (from X's perspective)
//	Character 25: X's checkers borne off
//
// Each character represents:
//
//	'-' = empty point
//	'A'-'O' (uppercase) = 1-15 X checkers
//	'a'-'o' (lowercase) = 1-15 O checkers
func parseXGIDBoard(pos *Position, boardStr string) {
	// Initialize board
	for i := range pos.Board {
		pos.Board[i] = 0
	}
	pos.OnBar["X"] = 0
	pos.OnBar["O"] = 0

	if len(boardStr) < 26 {
		return // Invalid XGID
	}

	// Character 0: X's bar
	if boardStr[0] >= 'A' && boardStr[0] <= 'O' {
		pos.OnBar["X"] = int(boardStr[0] - 'A' + 1)
	} else if boardStr[0] >= 'a' && boardStr[0] <= 'o' {
		pos.OnBar["O"] = int(boardStr[0] - 'a' + 1)
	}

	// Characters 1-24: Points 1, 2, 3, ..., 24
	for i := 1; i <= 24; i++ {
		ch := boardStr[i]
		point := i // Point 1, 2, 3, ..., 24

		if ch == '-' {
			pos.Board[point] = 0
		} else if ch >= 'A' && ch <= 'O' {
			// Uppercase = player X checkers (1-15)
			count := int(ch - 'A' + 1)
			pos.Board[point] = count
		} else if ch >= 'a' && ch <= 'o' {
			// Lowercase = player O checkers (1-15, stored as negative)
			count := int(ch - 'a' + 1)
			pos.Board[point] = -count
		}
	}

	// Character 25: X's borne off (we don't track this in board array)
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

// parseEquityInfo parses equity information lines in cube decision analysis
// Formats:
//
//	"Equity Red (cubeless): 0.139  Std.Dev.: 0.132"
//	"Equity (cubeful)    :  0.226"
func parseEquityInfo(line string, pos *Position) {
	line = strings.TrimSpace(line)

	// Parse cubeless equity and standard deviation
	// English: "Equity ... (cubeless): X.XXX  Std.Dev.: X.XXX"
	// German: "Erwartungswert ... (ohne Doppler): X.XXX  Std.Abw.: X.XXX"
	// Japanese: "期待値 (エクイティ) ... (キューブなし): X.XXX  標準偏差: X.XXX"
	if strings.Contains(strings.ToLower(line), "cubeless") ||
		strings.Contains(line, "ohne Doppler") ||
		strings.Contains(line, "キューブなし") {
		re := regexp.MustCompile(`([+-]?\d+\.\d+)`)
		matches := re.FindAllString(line, -1)
		if len(matches) >= 1 {
			pos.CubelessEquity, _ = strconv.ParseFloat(matches[0], 64)
		}

		// Parse standard deviation
		// English: "Std.Dev.:", German: "Std.Abw.:", Japanese: "標準偏差:"
		if strings.Contains(line, "Std.Dev.") ||
			strings.Contains(line, "Std.Abw.") ||
			strings.Contains(line, "標準偏差") {
			if len(matches) >= 2 {
				pos.EquityStdDev, _ = strconv.ParseFloat(matches[1], 64)
			}
		}
	}

	// Parse cubeful equity
	// English: "Equity (cubeful) : X.XXX"
	// German: "Auszahlungserw. (mit Doppler) : X.XXX"
	// French: "Équité (avec videau) : X.XXX"
	// Japanese: "エクイティ (キューブ有り) : X.XXX"
	if strings.Contains(strings.ToLower(line), "cubeful") ||
		strings.Contains(line, "mit Doppler") ||
		strings.Contains(line, "avec videau") ||
		strings.Contains(line, "キューブ有り") {
		re := regexp.MustCompile(`([+-]?\d+\.\d+)`)
		matches := re.FindAllString(line, -1)
		if len(matches) >= 1 {
			pos.CubefulEquity, _ = strconv.ParseFloat(matches[0], 64)
		}
	}
}

// parseCubeDecision parses a cube decision line
func parseCubeDecision(line string) *CubeDecision {
	line = strings.TrimSpace(line)

	// Must contain a colon and decimal numbers to be a cube decision line
	if !strings.Contains(line, ":") {
		return nil
	}

	// Must have at least one decimal number
	re := regexp.MustCompile(`\d+\.\d+`)
	if !re.MatchString(line) {
		return nil
	}

	decision := &CubeDecision{}

	if strings.Contains(line, "*") {
		decision.IsBest = true
		line = strings.ReplaceAll(line, "*", "")
	}

	// Extract action name (everything before the first colon)
	parts := strings.SplitN(line, ":", 2)
	if len(parts) >= 1 {
		decision.Action = strings.TrimSpace(parts[0])
	}

	// Format: "Action : MWC (MWC_diff) EMG (EMG_diff)"
	// Example: " No Double : 0.226 ( 0.000) 0.287 ( 0.000)"

	// Parse differences in parentheses first (MWC diff, EMG diff)
	reDiff := regexp.MustCompile(`\(([+-]?\s*\d+\.\d+)\)`)
	diffMatches := reDiff.FindAllStringSubmatch(line, -1)

	// Remove parenthesized values to find non-parenthesized decimals
	lineWithoutParens := reDiff.ReplaceAllString(line, "")

	// Now find decimal numbers that are NOT in parentheses
	matches := re.FindAllString(lineWithoutParens, -1)

	// matches[0] = MWC, matches[1] = EMG
	if len(matches) >= 1 {
		decision.MWC, _ = strconv.ParseFloat(matches[0], 64)
	}
	if len(matches) >= 2 {
		decision.EMG, _ = strconv.ParseFloat(matches[1], 64)
	}

	if len(diffMatches) >= 1 {
		// Remove any spaces in the captured group
		diffStr := strings.TrimSpace(diffMatches[0][1])
		decision.MWCDiff, _ = strconv.ParseFloat(diffStr, 64)
	}
	if len(diffMatches) >= 2 {
		diffStr := strings.TrimSpace(diffMatches[1][1])
		decision.EMGDiff, _ = strconv.ParseFloat(diffStr, 64)
	}

	return decision
}
