package bgfparser

import (
	"fmt"
	"os"
)

// ParseBGF parses a BGBlitz BGF (binary match) file from disk
// This is a convenience wrapper around ParseBGFFromReader that handles file reading.
//
// BGF files consist of:
// 1. A JSON header line with format info (uncompressed)
// 2. The rest of the file is gzipped JSON data (optionally SMILE encoded)
func ParseBGF(filename string) (*Match, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, &ParseError{File: filename, Message: err.Error()}
	}
	defer file.Close()

	match, err := ParseBGFFromReader(file)
	if err != nil {
		// Add filename to error if not already present
		if parseErr, ok := err.(*ParseError); ok && parseErr.File == "" {
			parseErr.File = filename
			return nil, parseErr
		}
		return nil, err
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
