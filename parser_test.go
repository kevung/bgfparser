package bgfparser_test

import (
	"testing"

	"github.com/kevung/bgfparser"
)

func TestParseTXT_ValidFile(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder21_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if pos == nil {
		t.Fatal("ParseTXT returned nil position")
	}

	// Check basic fields
	if pos.MatchLength != 3 {
		t.Errorf("MatchLength = %d, want 3", pos.MatchLength)
	}

	if pos.OnRoll != "X" {
		t.Errorf("OnRoll = %s, want X", pos.OnRoll)
	}

	if pos.Dice[0] != 2 || pos.Dice[1] != 1 {
		t.Errorf("Dice = %v, want [2 1]", pos.Dice)
	}

	// Check evaluations were parsed
	if len(pos.Evaluations) == 0 {
		t.Error("No evaluations parsed")
	}
}

func TestParseTXT_FrenchFile(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder32_FR.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed on French file: %v", err)
	}

	if pos == nil {
		t.Fatal("ParseTXT returned nil position")
	}

	// French files should parse the same way
	if pos.MatchLength != 3 {
		t.Errorf("MatchLength = %d, want 3", pos.MatchLength)
	}

	if len(pos.Evaluations) == 0 {
		t.Error("No evaluations parsed from French file")
	}
}

func TestParseTXT_WithCubeDecision(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/BlunderCubeOffered_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if len(pos.CubeDecisions) == 0 {
		t.Error("CubeDecisions not parsed")
	}
}

func TestParseTXT_NonExistentFile(t *testing.T) {
	_, err := bgfparser.ParseTXT("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParseBGF_ValidFile(t *testing.T) {
	match, err := bgfparser.ParseBGF("tmp/TachiAI_V_player_Nov_2__2025__16_55.bgf")

	// SMILE decoding is now supported, so we should not get an error
	if err != nil {
		t.Fatalf("ParseBGF failed: %v", err)
	}

	if match == nil {
		t.Fatal("ParseBGF returned nil match")
	}

	if match.Format != "BGF" {
		t.Errorf("Format = %s, want BGF", match.Format)
	}

	if match.Version != "1.0" {
		t.Errorf("Version = %s, want 1.0", match.Version)
	}

	if !match.Compress {
		t.Error("Expected Compress to be true")
	}

	if !match.UseSmile {
		t.Error("Expected UseSmile to be true")
	}
}

func TestParseBGF_NonExistentFile(t *testing.T) {
	_, err := bgfparser.ParseBGF("nonexistent.bgf")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestPosition_XGIDParsing(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder21_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if pos.XGID == "" {
		t.Error("XGID not parsed")
	}

	// Check that XGID parsing extracted some values
	if pos.CubeValue == 0 {
		t.Error("CubeValue not extracted from XGID")
	}
}

func TestEvaluation_Ranking(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder21_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if len(pos.Evaluations) < 2 {
		t.Skip("Need at least 2 evaluations for ranking test")
	}

	// First evaluation should be rank 1
	if pos.Evaluations[0].Rank != 1 {
		t.Errorf("First evaluation rank = %d, want 1", pos.Evaluations[0].Rank)
	}

	// Evaluations should be in order
	for i := 1; i < len(pos.Evaluations); i++ {
		if pos.Evaluations[i].Rank != i+1 {
			t.Errorf("Evaluation %d rank = %d, want %d",
				i, pos.Evaluations[i].Rank, i+1)
		}
	}
}

func TestEvaluation_EquityValues(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder21_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if len(pos.Evaluations) == 0 {
		t.Fatal("No evaluations to test")
	}

	// All evaluations should have equity values
	for i, eval := range pos.Evaluations {
		if eval.Equity == 0 && eval.Move != "" {
			t.Errorf("Evaluation %d has zero equity", i)
		}
	}
}
