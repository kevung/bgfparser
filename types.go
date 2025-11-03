// Package bgfparser provides parsers for BGBlitz backgammon data files.
//
// BGBlitz is a backgammon analysis software that stores data in two formats:
//   - TXT: Human-readable position files showing board state and analysis
//   - BGF: Binary match files using JSON with gzip compression and SMILE encoding
//
// This package can parse both formats and extract position data, evaluations,
// match information, and move analysis.
package bgfparser

import (
	"fmt"
)

// Position represents a backgammon position
type Position struct {
	// Board represents the checkers on each point (1-24), bar, and off
	// Points 1-24 for both players, positive for X, negative for O
	Board [26]int // 0=unused, 1-24=points, 25=bar

	// Player information
	PlayerX string
	PlayerO string
	ScoreX  int
	ScoreO  int

	// Match information
	MatchLength int
	Crawford    bool

	// Position identifiers
	PositionID string // BGBlitz Position-ID
	MatchID    string // BGBlitz Match-ID
	XGID       string // XG format ID

	// Current state
	OnRoll    string // "X" or "O"
	Dice      [2]int
	CubeValue int
	CubeOwner string // "", "X", "O"
	OnBar     map[string]int
	PipCount  map[string]int

	// Evaluation data
	Evaluations  []Evaluation
	CubeDecision *CubeDecision
}

// Evaluation represents a move evaluation
type Evaluation struct {
	Rank   int
	Move   string
	Equity float64
	Diff   float64
	Win    float64
	WinG   float64
	WinBG  float64
	LoseG  float64
	LoseBG float64
	IsBest bool
}

// CubeDecision represents a cube decision analysis
type CubeDecision struct {
	Action  string  // "Double/Take", "Double/Pass", "No Double"
	MWC     float64 // Match Winning Chances
	MWCDiff float64
	EMG     float64 // EMG (Normalized equity)
	EMGDiff float64
	IsBest  bool
}

// Match represents a complete backgammon match from a BGF file
type Match struct {
	Format   string `json:"format"`
	Version  string `json:"version"`
	Compress bool   `json:"compress"`
	UseSmile bool   `json:"useSmile"`

	// Match data will be populated from the JSON structure
	Data map[string]interface{}

	// DecodingWarning contains any warnings from partial decoding
	DecodingWarning string
}

// ParseError represents an error during parsing
type ParseError struct {
	File    string
	Line    int
	Message string
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.File, e.Message)
}
