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
	Board [26]int `json:"board"`

	// Player information
	PlayerX string `json:"player_x"`
	PlayerO string `json:"player_o"`
	ScoreX  int    `json:"score_x"`
	ScoreO  int    `json:"score_o"`

	// Match information
	MatchLength int  `json:"match_length"`
	Crawford    bool `json:"crawford"`

	// Position identifiers
	PositionID string `json:"position_id"` // BGBlitz Position-ID
	MatchID    string `json:"match_id"`    // BGBlitz Match-ID
	XGID       string `json:"xgid"`        // XG format ID

	// Current state
	OnRoll    string         `json:"on_roll"` // "X" or "O"
	Dice      [2]int         `json:"dice"`
	CubeValue int            `json:"cube_value"`
	CubeOwner string         `json:"cube_owner"` // "", "X", "O"
	OnBar     map[string]int `json:"on_bar"`
	PipCount  map[string]int `json:"pip_count"`

	// Evaluation data
	Evaluations  []Evaluation  `json:"evaluations,omitempty"`
	CubeDecision *CubeDecision `json:"cube_decision,omitempty"`
}

// Evaluation represents a move evaluation
type Evaluation struct {
	Rank   int     `json:"rank"`
	Move   string  `json:"move"`
	Equity float64 `json:"equity"`
	Diff   float64 `json:"diff"`
	Win    float64 `json:"win"`
	WinG   float64 `json:"win_g"`
	WinBG  float64 `json:"win_bg"`
	LoseG  float64 `json:"lose_g"`
	LoseBG float64 `json:"lose_bg"`
	IsBest bool    `json:"is_best"`
}

// CubeDecision represents a cube decision analysis
type CubeDecision struct {
	Action  string  `json:"action"` // "Double/Take", "Double/Pass", "No Double"
	MWC     float64 `json:"mwc"`    // Match Winning Chances
	MWCDiff float64 `json:"mwc_diff"`
	EMG     float64 `json:"emg"` // EMG (Normalized equity)
	EMGDiff float64 `json:"emg_diff"`
	IsBest  bool    `json:"is_best"`
}

// Match represents a complete backgammon match from a BGF file
type Match struct {
	Format   string `json:"format"`
	Version  string `json:"version"`
	Compress bool   `json:"compress"`
	UseSmile bool   `json:"useSmile"`

	// Match data will be populated from the JSON structure
	Data map[string]interface{} `json:"data,omitempty"`
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
