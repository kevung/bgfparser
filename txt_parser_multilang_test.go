package bgfparser_test

import (
	"testing"

	"github.com/kevung/bgfparser"
)

// TestParseTXT_Multilingual tests parsing of position files in different languages
func TestParseTXT_Multilingual(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		lang    string
		playerO string
		playerX string
		scoreO  int
		scoreX  int
		match   int
		onRoll  string
		dice    [2]int
		posID   string
		matchID string
		xgid    string
	}{
		{
			name:    "English",
			file:    "test/2025-11-04/01_checkerPosition_EN.txt",
			lang:    "EN",
			playerO: "Green",
			playerX: "Red",
			scoreO:  6,
			scoreX:  3,
			match:   7,
			onRoll:  "X",
			dice:    [2]int{1, 2},
			posID:   "b9sBCIC5bYDQAA",
			matchID: "QYnoAGAAGAAE",
			xgid:    "-B-CBBB---a---A---ABcbbbd-:1:-1:1:21:3:6:0:7:10",
		},
		{
			name:    "French",
			file:    "test/2025-11-04/01_checkerPosition_FR.txt",
			lang:    "FR",
			playerO: "Vert",
			playerX: "Rouge",
			scoreO:  6,
			scoreX:  3,
			match:   7,
			onRoll:  "X",
			dice:    [2]int{1, 2},
			posID:   "b9sBCIC5bYDQAA",
			matchID: "QYnoAGAAGAAE",
			xgid:    "-B-CBBB---a---A---ABcbbbd-:1:-1:1:21:3:6:0:7:10",
		},
		{
			name:    "German",
			file:    "test/2025-11-04/01_checkerPosition_DE.txt",
			lang:    "DE",
			playerO: "Grün",
			playerX: "Rot",
			scoreO:  6,
			scoreX:  3,
			match:   7,
			onRoll:  "X",
			dice:    [2]int{1, 2},
			posID:   "b9sBCIC5bYDQAA",
			matchID: "QYnoAGAAGAAE",
			xgid:    "-B-CBBB---a---A---ABcbbbd-:1:-1:1:21:3:6:0:7:10",
		},
		{
			name:    "Japanese",
			file:    "test/2025-11-04/01_checkerPosition_JP.txt",
			lang:    "JP",
			playerO: "緑",
			playerX: "赤",
			scoreO:  6,
			scoreX:  3,
			match:   7,
			onRoll:  "X",
			dice:    [2]int{1, 2},
			posID:   "b9sBCIC5bYDQAA",
			matchID: "QYnoAGAAGAAE",
			xgid:    "-B-CBBB---a---A---ABcbbbd-:1:-1:1:21:3:6:0:7:10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := bgfparser.ParseTXT(tt.file)
			if err != nil {
				t.Fatalf("ParseTXT failed for %s: %v", tt.lang, err)
			}

			if pos == nil {
				t.Fatal("ParseTXT returned nil position")
			}

			// Check player names
			if pos.PlayerO != tt.playerO {
				t.Errorf("PlayerO = %q, want %q", pos.PlayerO, tt.playerO)
			}
			if pos.PlayerX != tt.playerX {
				t.Errorf("PlayerX = %q, want %q", pos.PlayerX, tt.playerX)
			}

			// Check scores
			if pos.ScoreO != tt.scoreO {
				t.Errorf("ScoreO = %d, want %d", pos.ScoreO, tt.scoreO)
			}
			if pos.ScoreX != tt.scoreX {
				t.Errorf("ScoreX = %d, want %d", pos.ScoreX, tt.scoreX)
			}

			// Check match length
			if pos.MatchLength != tt.match {
				t.Errorf("MatchLength = %d, want %d", pos.MatchLength, tt.match)
			}

			// Check on roll
			if pos.OnRoll != tt.onRoll {
				t.Errorf("OnRoll = %q, want %q", pos.OnRoll, tt.onRoll)
			}

			// Check dice
			if pos.Dice[0] != tt.dice[0] || pos.Dice[1] != tt.dice[1] {
				t.Errorf("Dice = %v, want %v", pos.Dice, tt.dice)
			}

			// Check IDs
			if pos.PositionID != tt.posID {
				t.Errorf("PositionID = %q, want %q", pos.PositionID, tt.posID)
			}
			if pos.MatchID != tt.matchID {
				t.Errorf("MatchID = %q, want %q", pos.MatchID, tt.matchID)
			}
			if pos.XGID != tt.xgid {
				t.Errorf("XGID = %q, want %q", pos.XGID, tt.xgid)
			}

			// Check evaluations
			if len(pos.Evaluations) != 5 {
				t.Errorf("Expected 5 evaluations, got %d", len(pos.Evaluations))
			}

			// Check pip counts (should be 52 for O, 111 for X based on the files)
			if pos.PipCount["O"] != 52 {
				t.Errorf("PipCount[O] = %d, want 52", pos.PipCount["O"])
			}
			if pos.PipCount["X"] != 111 {
				t.Errorf("PipCount[X] = %d, want 111", pos.PipCount["X"])
			}
		})
	}
}

// TestParseTXT_EvaluationSection tests that evaluation sections are parsed correctly
// regardless of the language used for the "Evaluation" header
func TestParseTXT_EvaluationSection(t *testing.T) {
	tests := []struct {
		name string
		file string
		lang string
	}{
		{"English", "test/2025-11-04/01_checkerPosition_EN.txt", "EN"},
		{"French", "test/2025-11-04/01_checkerPosition_FR.txt", "FR"},
		{"German", "test/2025-11-04/01_checkerPosition_DE.txt", "DE"},
		{"Japanese", "test/2025-11-04/01_checkerPosition_JP.txt", "JP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := bgfparser.ParseTXT(tt.file)
			if err != nil {
				t.Fatalf("ParseTXT failed: %v", err)
			}

			if len(pos.Evaluations) == 0 {
				t.Fatal("No evaluations parsed")
			}

			// Check first evaluation has expected move
			firstMove := "19/18, 14/12"
			if pos.Evaluations[0].Move != firstMove {
				t.Errorf("First move = %q, want %q", pos.Evaluations[0].Move, firstMove)
			}

			// Check equity values are parsed (EMG equity, not MWP)
			expectedEquity := -0.492
			if pos.Evaluations[0].Equity != expectedEquity {
				t.Errorf("First equity = %f, want %f", pos.Evaluations[0].Equity, expectedEquity)
			}

			// Check ranking
			for i, eval := range pos.Evaluations {
				if eval.Rank != i+1 {
					t.Errorf("Evaluation %d: rank = %d, want %d", i, eval.Rank, i+1)
				}
			}
		})
	}
}

// TestParseTXT_ProbabilityExtraction tests that win/lose probabilities are correctly extracted
func TestParseTXT_ProbabilityExtraction(t *testing.T) {
	pos, err := bgfparser.ParseTXT("tmp/blunder21_EN.txt")
	if err != nil {
		t.Fatalf("ParseTXT failed: %v", err)
	}

	if len(pos.Evaluations) == 0 {
		t.Fatal("No evaluations parsed")
	}

	// Check first evaluation has probabilities
	eval := pos.Evaluations[0]

	// Expected values from the file:
	// 1) 13-11 24-23                0.473 / -0.289
	//    0.443  0.113  0.002  -  0.557  0.179  0.003

	if eval.Win == 0 {
		t.Error("Win probability not extracted")
	}

	expectedWin := 0.443
	if eval.Win != expectedWin {
		t.Errorf("Win = %.3f, want %.3f", eval.Win, expectedWin)
	}

	expectedWinG := 0.113
	if eval.WinG != expectedWinG {
		t.Errorf("WinG = %.3f, want %.3f", eval.WinG, expectedWinG)
	}

	expectedWinBG := 0.002
	if eval.WinBG != expectedWinBG {
		t.Errorf("WinBG = %.3f, want %.3f", eval.WinBG, expectedWinBG)
	}

	expectedLoseG := 0.179
	if eval.LoseG != expectedLoseG {
		t.Errorf("LoseG = %.3f, want %.3f", eval.LoseG, expectedLoseG)
	}

	expectedLoseBG := 0.003
	if eval.LoseBG != expectedLoseBG {
		t.Errorf("LoseBG = %.3f, want %.3f", eval.LoseBG, expectedLoseBG)
	}
}
