package bgfparser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/kevung/bgfparser/internal/smile"
)

// ParseBGFFromReader parses a BGBlitz BGF file from an io.Reader
// This allows parsing BGF files from network streams, memory buffers, HTTP uploads,
// or any io.Reader source.
//
// Example usage with HTTP upload:
//
//	func uploadHandler(w http.ResponseWriter, r *http.Request) {
//	    file, _, _ := r.FormFile("bgffile")
//	    defer file.Close()
//
//	    match, err := bgfparser.ParseBGFFromReader(file)
//	    if err != nil {
//	        http.Error(w, err.Error(), http.StatusBadRequest)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "application/json")
//	    json.NewEncoder(w).Encode(match)
//	}
func ParseBGFFromReader(reader io.Reader) (*Match, error) {
	bufReader := bufio.NewReader(reader)

	// Read first line (JSON header)
	headerLine, err := bufReader.ReadBytes('\n')
	if err != nil {
		return nil, &ParseError{Message: "failed to read header: " + err.Error()}
	}

	// Parse header
	match := &Match{}
	if err := json.Unmarshal(headerLine, match); err != nil {
		return nil, &ParseError{Message: "failed to parse header: " + err.Error()}
	}

	// Read the rest of the data
	restData, err := io.ReadAll(bufReader)
	if err != nil {
		return nil, &ParseError{Message: "failed to read data: " + err.Error()}
	}

	// Decompress if compressed
	var jsonData []byte
	if match.Compress {
		gzReader, err := gzip.NewReader(bytes.NewReader(restData))
		if err != nil {
			return nil, &ParseError{Message: "failed to create gzip reader: " + err.Error()}
		}
		defer gzReader.Close()

		jsonData, err = io.ReadAll(gzReader)
		if err != nil {
			return nil, &ParseError{Message: "failed to decompress: " + err.Error()}
		}
	} else {
		jsonData = restData
	}

	// Handle SMILE encoding
	if match.UseSmile {
		var data interface{}
		if err := smile.Unmarshal(jsonData, &data); err != nil {
			return nil, &ParseError{Message: "failed to decode SMILE: " + err.Error()}
		}

		if dataMap, ok := data.(map[string]interface{}); ok {
			match.Data = dataMap
		} else {
			match.Data = map[string]interface{}{"_data": data}
		}
	} else {
		if err := json.Unmarshal(jsonData, &match.Data); err != nil {
			return nil, &ParseError{Message: "failed to parse JSON: " + err.Error()}
		}
	}

	return match, nil
}

// ParseTXTFromReader parses a BGBlitz TXT position file from an io.Reader
// This allows parsing TXT files from network streams, memory buffers, HTTP uploads,
// or any io.Reader source.
//
// Example usage with in-memory data:
//
//	data := []byte("... TXT content ...")
//	pos, err := bgfparser.ParseTXTFromReader(bytes.NewReader(data))
func ParseTXTFromReader(reader io.Reader) (*Position, error) {
	pos := &Position{
		OnBar:    make(map[string]int),
		PipCount: make(map[string]int),
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	boardLines := []string{}
	inEvaluation := false
	inCubeDecision := false
	evalRank := 0
	var lastEval *Evaluation

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Parse board lines
		if parseBoardLine(line, &boardLines) {
			continue
		}

		// Parse player names and scores
		parsePlayerInfo(line, pos)

		// Parse Position-ID, Match-ID
		parsePositionID(line, pos)

		// Parse XGID
		parseXGIDLine(line, pos)

		// Parse match score
		parseMatchScore(line, pos)

		// Parse current player to move
		parseCurrentPlayer(line, pos)

		// Parse cube value
		if parseCubeValue(line, scanner, pos) {
			continue
		}

		// Handle evaluation sections
		if handleEvaluationSection(line, &inEvaluation, &inCubeDecision, &evalRank) {
			continue
		}

		// Parse evaluations
		if inEvaluation && len(line) > 0 {
			if eval := parseEvaluation(line, &evalRank); eval != nil {
				pos.Evaluations = append(pos.Evaluations, *eval)
				lastEval = &pos.Evaluations[len(pos.Evaluations)-1]
			} else if lastEval != nil {
				// Try to parse probability line for the last evaluation
				if parseProbabilityLine(line, lastEval) {
					lastEval = nil // Reset after parsing probabilities
				}
			}
		}

		// Parse cube decisions
		if inCubeDecision {
			if decision := parseCubeDecision(line); decision != nil {
				pos.CubeDecision = decision
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, &ParseError{Message: err.Error()}
	}

	// Parse the board from collected lines
	if len(boardLines) > 0 {
		parseBoard(pos, boardLines)
	}

	return pos, nil
}

// ToJSON serializes the Match to JSON
func (m *Match) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// ToJSON serializes the Position to JSON
func (p *Position) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
