# Complete Information Extraction from BGBlitz TXT Files

This document shows ALL information that the `ParseTXT` parser extracts from BGBlitz position files.

## Data Structure: `Position`

The parser returns a `Position` struct with the following fields:

### 1. Board Position (26-element array)
```go
Board [26]int
```
- **Description**: Checkers on each point (1-24), bar (25), and off (0)
- **Format**: Positive values = X checkers, Negative values = O checkers
- **Example**: `[0 2 0 0 0 -5 0 3 0 0 0 -5 5 -3 0 0 3 0 5 0 0 0 0 -2 0 0]`

### 2. Player Information
```go
PlayerX  string  // Player X name (e.g., "player", "Red", "Rouge", "Rot", "赤")
PlayerO  string  // Player O name (e.g., "TachiAI_V", "Green", "Vert", "Grün", "緑")
ScoreX   int     // Current score of Player X
ScoreO   int     // Current score of Player O
```
- **Multilingual**: Supports English, French, German, Japanese player names
- **Example**: `PlayerX: "Rouge"`, `PlayerO: "Vert"`, `ScoreX: 3`, `ScoreO: 6`

### 3. Match Information
```go
MatchLength  int   // Match length in points (e.g., 3, 7, 9)
Crawford     bool  // Whether Crawford rule is in effect
```
- **Example**: `MatchLength: 7`, `Crawford: false`

### 4. Position Identifiers
```go
PositionID  string  // BGBlitz Position-ID
MatchID     string  // BGBlitz Match-ID
XGID        string  // XG format position ID
```
- **Example**: 
  - `PositionID: "b9sBCIC5bYDQAA"`
  - `MatchID: "QYnoAGAAGAAE"`
  - `XGID: "-B-CBBB---a---A---ABcbbbd-:1:-1:1:21:3:6:0:7:10"`

### 5. Current Position State
```go
OnRoll     string       // "X" or "O" - who is on roll
Dice       [2]int       // Dice rolled [die1, die2]
CubeValue  int          // Current doubling cube value (1, 2, 4, 8, etc.)
CubeOwner  string       // "", "X", or "O" (empty = centered)
OnBar      map[string]int  // Checkers on bar: {"X": 0, "O": 0}
PipCount   map[string]int  // Pip counts: {"X": 167, "O": 161}
```
- **Example**: 
  - `OnRoll: "X"`
  - `Dice: [1, 2]`
  - `CubeValue: 2`
  - `CubeOwner: "O"`
  - `PipCount: {"X": 111, "O": 52}`

### 6. Move Evaluations (Array)
```go
Evaluations []Evaluation
```

Each `Evaluation` contains:
```go
type Evaluation struct {
    Rank   int      // Rank of this move (1, 2, 3, ...)
    Move   string   // Move notation (e.g., "19/18, 14/12")
    Equity float64  // Equity value (e.g., 0.473)
    Diff   float64  // Equity difference from best move (e.g., -0.289)
    
    // Win probabilities
    Win    float64  // Win probability (e.g., 0.443 = 44.3%)
    WinG   float64  // Win gammon probability (e.g., 0.113 = 11.3%)
    WinBG  float64  // Win backgammon probability (e.g., 0.002 = 0.2%)
    
    // Lose probabilities
    LoseG  float64  // Lose gammon probability (e.g., 0.179 = 17.9%)
    LoseBG float64  // Lose backgammon probability (e.g., 0.003 = 0.3%)
    
    IsBest bool     // Whether this is marked as best move
}
```

**Example Evaluation:**
```
Rank: 1
Move: "13-11 24-23"
Equity: 0.473
Diff: -0.289
Win: 0.443 (44.3%)
WinG: 0.113 (11.3%)
WinBG: 0.002 (0.2%)
Lose: 0.557 (55.7%) [calculated as 1 - Win]
LoseG: 0.179 (17.9%)
LoseBG: 0.003 (0.3%)
IsBest: false
```

### 7. Cube Decision (Optional)
```go
CubeDecision *CubeDecision
```

The `CubeDecision` contains:
```go
type CubeDecision struct {
    Action  string   // "Double/Take", "Double/Pass", "No Double"
    MWC     float64  // Match Winning Chances
    MWCDiff float64  // MWC difference from best
    EMG     float64  // EMG (Normalized equity)
    EMGDiff float64  // EMG difference from best
    IsBest  bool     // Whether this is the best cube action
}
```

**Example Cube Decision:**
```
Action: "No Double"
MWC: 0.226
MWCDiff: 0.000
EMG: 0.287
EMGDiff: 0.000
IsBest: true
```

## Complete Example Output

Here's what gets extracted from a typical position file:

```
========================================
COMPLETE POSITION INFORMATION
========================================

PLAYER INFORMATION
├─ Player X: player
├─ Player O: TachiAI_V
├─ Score X:  0
└─ Score O:  0

MATCH INFORMATION
├─ Match Length: 3 points
└─ Crawford:     false

POSITION IDENTIFIERS
├─ Position-ID: mGfwATDgc/ABMA
├─ Match-ID:    cAllAAAAAAAE
└─ XGID:        -b----E-C---eE---b-d-b--B-:0:0:1:21:0:0:0:3:10

CURRENT POSITION STATE
├─ On Roll:      X
├─ Dice:         2-1
├─ Cube Value:   1
└─ Cube Owner:   (centered)

PIP COUNTS
├─ Player X: 167 pips
└─ Player O: 161 pips

BOARD ARRAY
└─ Board: [26 integers representing checker positions]

MOVE EVALUATIONS (8 total)
├─ #1: 13-11 24-23
│   ├─ Equity:      0.4730
│   ├─ Difference:  -0.2890
│   ├─ Win:         0.443 (44.3%)
│   ├─ Win Gammon:  0.113 (11.3%)
│   ├─ Win BG:      0.002 (0.2%)
│   ├─ Lose:        0.557 (55.7%)
│   ├─ Lose Gammon: 0.179 (17.9%)
│   └─ Lose BG:     0.003 (0.3%)
├─ #2: 24-22 24-23
│   └─ [similar structure]
...
└─ #8: 8-6 6-5
    └─ [similar structure]

CUBE DECISION (if present)
├─ Action:   No Double ★ BEST
├─ MWC:      0.2260 (+0.0000)
└─ EMG:      0.2870 (+0.0000)
```

## Multilingual Support

The parser correctly extracts information from files in:
- **English**: "Evaluation", "Green", "Red"
- **French**: "Évaluation", "Vert", "Rouge"
- **German**: "Bewertung", "Grün", "Rot"
- **Japanese**: "評価", "緑", "赤"

All fields are extracted identically regardless of language.

## Usage

### Basic Usage
```go
pos, err := bgfparser.ParseTXT("position.txt")
if err != nil {
    log.Fatal(err)
}

// Access any field
fmt.Printf("Player: %s vs %s\n", pos.PlayerX, pos.PlayerO)
fmt.Printf("Score: %d - %d\n", pos.ScoreX, pos.ScoreO)
fmt.Printf("Pip Count: %d vs %d\n", pos.PipCount["X"], pos.PipCount["O"])

// Iterate through evaluations
for _, eval := range pos.Evaluations {
    fmt.Printf("%d) %s - Equity: %.3f\n", eval.Rank, eval.Move, eval.Equity)
    fmt.Printf("   Win: %.1f%%, WinG: %.1f%%\n", eval.Win*100, eval.WinG*100)
}
```

### JSON Export
```go
jsonData, err := pos.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonData))
```

## Summary of Extracted Fields

| Category | Field | Type | Description |
|----------|-------|------|-------------|
| **Board** | Board | [26]int | Checker positions |
| **Players** | PlayerX | string | Player X name |
| | PlayerO | string | Player O name |
| | ScoreX | int | Player X score |
| | ScoreO | int | Player O score |
| **Match** | MatchLength | int | Match length in points |
| | Crawford | bool | Crawford rule status |
| **IDs** | PositionID | string | BGBlitz Position-ID |
| | MatchID | string | BGBlitz Match-ID |
| | XGID | string | XG format ID |
| **State** | OnRoll | string | "X" or "O" |
| | Dice | [2]int | Dice rolled |
| | CubeValue | int | Cube value |
| | CubeOwner | string | "", "X", or "O" |
| | OnBar | map[string]int | Checkers on bar |
| | PipCount | map[string]int | Pip counts |
| **Evaluations** | Rank | int | Move rank |
| | Move | string | Move notation |
| | Equity | float64 | Equity value |
| | Diff | float64 | Equity difference |
| | Win | float64 | Win probability |
| | WinG | float64 | Win gammon probability |
| | WinBG | float64 | Win backgammon probability |
| | LoseG | float64 | Lose gammon probability |
| | LoseBG | float64 | Lose backgammon probability |
| | IsBest | bool | Best move marker |
| **Cube** | Action | string | Cube action |
| | MWC | float64 | Match winning chances |
| | MWCDiff | float64 | MWC difference |
| | EMG | float64 | EMG value |
| | EMGDiff | float64 | EMG difference |
| | IsBest | bool | Best action marker |

**Total: 30+ distinct fields extracted from each position file!**
