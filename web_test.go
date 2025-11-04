package bgfparser

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseBGFFromReader(t *testing.T) {
	// Test with a minimal BGF header
	header := `{"format":"BGF","version":"1.0","compress":false,"useSmile":false}` + "\n"
	data := `{"test":"data","playerX":"Player1","playerO":"Player2"}`

	content := header + data
	reader := bytes.NewReader([]byte(content))

	match, err := ParseBGFFromReader(reader)
	if err != nil {
		t.Fatalf("ParseBGFFromReader failed: %v", err)
	}

	if match.Format != "BGF" {
		t.Errorf("Expected format BGF, got %s", match.Format)
	}

	if match.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", match.Version)
	}

	if match.Compress {
		t.Errorf("Expected compress false, got true")
	}

	if match.Data == nil {
		t.Error("Expected Data to be populated")
	}

	if match.Data["test"] != "data" {
		t.Error("Expected test data to be preserved")
	}
}

func TestParseTXTFromReader(t *testing.T) {
	txtContent := `O: Player1 150  X: Player2 140

Position-ID: testpos123    Match-ID: testmatch456
XGID=-b----E-C---eE---b-d-b--B-:0:0:1:21:0:0:0:3:10

Player1 - 5 Player2 - 3 in a 7 point match.
Player2 to move 3-2
`

	reader := strings.NewReader(txtContent)
	pos, err := ParseTXTFromReader(reader)
	if err != nil {
		t.Fatalf("ParseTXTFromReader failed: %v", err)
	}

	if pos.PlayerO != "Player1" {
		t.Errorf("Expected PlayerO to be Player1, got %s", pos.PlayerO)
	}

	if pos.PlayerX != "Player2" {
		t.Errorf("Expected PlayerX to be Player2, got %s", pos.PlayerX)
	}

	if pos.MatchLength != 7 {
		t.Errorf("Expected MatchLength 7, got %d", pos.MatchLength)
	}

	if pos.ScoreO != 5 {
		t.Errorf("Expected ScoreO 5, got %d", pos.ScoreO)
	}

	if pos.ScoreX != 3 {
		t.Errorf("Expected ScoreX 3, got %d", pos.ScoreX)
	}

	if pos.OnRoll != "X" {
		t.Errorf("Expected OnRoll X, got %s", pos.OnRoll)
	}

	if pos.Dice[0] != 3 || pos.Dice[1] != 2 {
		t.Errorf("Expected dice [3 2], got [%d %d]", pos.Dice[0], pos.Dice[1])
	}
}

func TestMatchToJSON(t *testing.T) {
	match := &Match{
		Format:   "BGF",
		Version:  "1.0",
		Compress: false,
		UseSmile: false,
		Data: map[string]interface{}{
			"test": "data",
		},
	}

	jsonData, err := match.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !bytes.Contains(jsonData, []byte("BGF")) {
		t.Error("Expected JSON to contain BGF")
	}

	if !bytes.Contains(jsonData, []byte("test")) {
		t.Error("Expected JSON to contain test data")
	}
}

func TestPositionToJSON(t *testing.T) {
	pos := &Position{
		PlayerX:     "PlayerA",
		PlayerO:     "PlayerB",
		MatchLength: 5,
		ScoreX:      2,
		ScoreO:      3,
		OnRoll:      "X",
		Dice:        [2]int{4, 5},
		OnBar:       make(map[string]int),
		PipCount:    make(map[string]int),
	}

	jsonData, err := pos.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !bytes.Contains(jsonData, []byte("PlayerA")) {
		t.Error("Expected JSON to contain PlayerA")
	}

	if !bytes.Contains(jsonData, []byte("PlayerB")) {
		t.Error("Expected JSON to contain PlayerB")
	}
}
