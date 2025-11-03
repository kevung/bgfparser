package bgfparser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ParseBGF parses a BGBlitz BGF (binary match) file
// BGF files consist of:
// 1. A JSON header line with format info (uncompressed)
// 2. The rest of the file is gzipped JSON data (optionally SMILE encoded)
func ParseBGF(filename string) (*Match, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read first line (JSON header)
	headerLine, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, &ParseError{File: filename, Message: "failed to read header: " + err.Error()}
	}

	// Parse header
	match := &Match{}
	if err := json.Unmarshal(headerLine, match); err != nil {
		return nil, &ParseError{File: filename, Message: "failed to parse header: " + err.Error()}
	}

	// Read the rest of the file
	restData, err := io.ReadAll(reader)
	if err != nil {
		return nil, &ParseError{File: filename, Message: "failed to read data: " + err.Error()}
	}

	// Decompress if compressed
	var jsonData []byte
	if match.Compress {
		gzReader, err := gzip.NewReader(bytes.NewReader(restData))
		if err != nil {
			return nil, &ParseError{File: filename, Message: "failed to create gzip reader: " + err.Error()}
		}
		defer gzReader.Close()

		jsonData, err = io.ReadAll(gzReader)
		if err != nil {
			return nil, &ParseError{File: filename, Message: "failed to decompress: " + err.Error()}
		}
	} else {
		jsonData = restData
	}

	// Handle SMILE encoding
	if match.UseSmile {
		// SMILE is a binary JSON format
		// For now, we'll note that SMILE decoding would require a SMILE library
		// The go-smile library could be used, but it's not in the standard library
		// Return match with header info but indicate SMILE is not supported
		return match, &ParseError{
			File:    filename,
			Message: "SMILE encoding is not yet supported. The data is compressed JSON in SMILE binary format. A SMILE decoder library is needed.",
		}
	} else {
		// Parse regular JSON
		if err := json.Unmarshal(jsonData, &match.Data); err != nil {
			return nil, &ParseError{File: filename, Message: "failed to parse JSON: " + err.Error()}
		}
	}

	return match, nil
}

// GetMatchInfo extracts basic match information from a parsed BGF file
func (m *Match) GetMatchInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["format"] = m.Format
	info["version"] = m.Version
	info["compress"] = m.Compress
	info["useSmile"] = m.UseSmile

	// Try to extract common fields from the data
	if m.Data != nil {
		for key, value := range m.Data {
			// Only extract top-level metadata
			switch key {
			case "playerX", "playerO", "matchLength", "score", "date", "event":
				info[key] = value
			}
		}
	}

	return info
}

// String returns a human-readable representation of the match
func (m *Match) String() string {
	info := m.GetMatchInfo()
	return fmt.Sprintf("BGF Match (Format: %s, Version: %s, Compressed: %v, SMILE: %v)",
		info["format"], info["version"], info["compress"], info["useSmile"])
}
