# BGF Parser API Reference


Import path: `github.com/kevung/bgfparser`

The bgfparser package provides parsers for BGBlitz backgammon data files in both TXT (position) and BGF (match) formats.

---

## Functions

### ParseTXT

```go
func ParseTXT(filename string) (*Position, error)

ParseTXT parses a BGBlitz position text file and returns a Position struct containing all extracted data.

**Parameters:**
- `filename` (string): Path to the TXT file to parse

**Returns:**
- `*Position`: Parsed position data
- `error`: Error if file cannot be read or parsed
**Example:**
```go
    log.Fatal(err)
}
fmt.Printf("Player on roll: %s\n", pos.OnRoll)
```

---

### ParseBGF

```go
func ParseBGF(filename string) (*Match, error)
```

ParseBGF parses a BGBlitz BGF binary match file. BGF files consist of a JSON header followed by gzipped data (optionally SMILE-encoded).

**Parameters:**
- `filename` (string): Path to the BGF file to parse

**Returns:**
- `*Match`: Parsed match data
- `error`: Error if file cannot be read or parsed

**Example:**
```go
match, err := bgfparser.ParseBGF("match.bgf")
if err != nil {
    if parseErr, ok := err.(*bgfparser.ParseError); ok {
        log.Printf("Warning: %s", parseErr.Message)
    }
}
fmt.Printf("Format: %s v%s\n", match.Format, match.Version)
```

---

## Types

### Position

```go
type Position struct {
    Board        [26]int
    PlayerX      string
    PlayerO      string
    ScoreX       int
    ScoreO       int
    MatchLength  int
    Crawford     bool
    PositionID   string
    MatchID      string
    XGID         string
    OnRoll       string
    Dice         [2]int
    CubeValue    int
    CubeOwner    string
    OnBar        map[string]int
    PipCount     map[string]int
    Evaluations  []Evaluation
    CubeDecision *CubeDecision
}
```

Represents a complete backgammon position with all associated metadata.

**Fields:**

- **Board** `[26]int`: Checker positions
  - Index 0: Unused
  - Index 1-24: Board points (positive for X, negative for O)
  - Index 25: Bar
  
- **PlayerX** `string`: Name of player X

- **PlayerO** `string`: Name of player O

- **ScoreX** `int`: Current match score for player X

- **ScoreO** `int`: Current match score for player O

- **MatchLength** `int`: Match length in points

- **Crawford** `bool`: Whether Crawford rule is in effect

- **PositionID** `string`: BGBlitz position identifier

- **MatchID** `string`: BGBlitz match identifier

- **XGID** `string`: eXtreme Gammon position identifier

- **OnRoll** `string`: Player to move ("X" or "O")

- **Dice** `[2]int`: Current dice roll

- **CubeValue** `int`: Current doubling cube value

- **CubeOwner** `string`: Owner of the cube ("", "X", or "O")

- **OnBar** `map[string]int`: Number of checkers on bar per player

- **PipCount** `map[string]int`: Pip count per player

- **Evaluations** `[]Evaluation`: Move evaluations

- **CubeDecision** `*CubeDecision`: Cube decision analysis (nil if not a cube decision)

---

### Evaluation

```go
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
```

Represents a single move evaluation with equity and statistics.

**Fields:**

- **Rank** `int`: Move ranking (1 = best move)

- **Move** `string`: Move in standard notation (e.g., "13-11 24-23")

- **Equity** `float64`: Equity value for this move

- **Diff** `float64`: Difference from best move (negative = worse)

- **Win** `float64`: Probability of winning (0.0 - 1.0)

- **WinG** `float64`: Probability of winning a gammon (0.0 - 1.0)

- **WinBG** `float64`: Probability of winning a backgammon (0.0 - 1.0)

- **LoseG** `float64`: Probability of losing a gammon (0.0 - 1.0)

- **LoseBG** `float64`: Probability of losing a backgammon (0.0 - 1.0)

- **IsBest** `bool`: True if marked with * in the position file

---

### CubeDecision

```go
type CubeDecision struct {
    Action  string
    MWC     float64
    MWCDiff float64
    EMG     float64
    EMGDiff float64
    IsBest  bool
}
```

Represents a cube decision analysis.

**Fields:**

- **Action** `string`: Recommended action
  - `"Double/Take"`: Should double, opponent should take
  - `"Double/Pass"`: Should double, opponent should pass
  - `"No Double"`: Should not double

- **MWC** `float64`: Match Winning Chances (0.0 - 1.0)

- **MWCDiff** `float64`: Difference in MWC from best action

- **EMG** `float64`: EMG (normalized equity)

- **EMGDiff** `float64`: Difference in EMG from best action

- **IsBest** `bool`: True if marked as best action

---

### Match

```go
type Match struct {
    Format   string                 `json:"format"`
    Version  string                 `json:"version"`
    Compress bool                   `json:"compress"`
    UseSmile bool                   `json:"useSmile"`
    Data     map[string]interface{}
}
```

Represents a BGF match file.

**Fields:**

- **Format** `string`: File format identifier (should be "BGF")

- **Version** `string`: Format version (e.g., "1.0")

- **Compress** `bool`: Whether data is gzip compressed

- **UseSmile** `bool`: Whether data uses SMILE binary JSON encoding

- **Data** `map[string]interface{}`: Parsed match data (nil if SMILE encoding used)

**Methods:**

#### GetMatchInfo

```go
func (m *Match) GetMatchInfo() map[string]interface{}
```

Returns a map containing match metadata extracted from header and data.

**Returns:**
- `map[string]interface{}`: Match information including format, version, player names, etc.

**Example:**
```go
info := match.GetMatchInfo()
fmt.Printf("Format: %s\n", info["format"])
```

#### String

```go
func (m *Match) String() string
```

Returns a human-readable string representation of the match.

**Returns:**
- `string`: Formatted match description

**Example:**
```go
fmt.Println(match.String())
// Output: BGF Match (Format: BGF, Version: 1.0, Compressed: true, SMILE: true)
```

---

### ParseError

```go
type ParseError struct {
    File    string
    Line    int
    Message string
}
```

Custom error type for parsing errors.

**Fields:**

- **File** `string`: Filename where error occurred

- **Line** `int`: Line number where error occurred (0 if not line-specific)

- **Message** `string`: Error description

**Methods:**

#### Error

```go
func (e *ParseError) Error() string
```

Returns the error message as a string.

**Example:**
```go
_, err := bgfparser.ParseTXT("invalid.txt")
if parseErr, ok := err.(*bgfparser.ParseError); ok {
    fmt.Printf("Error in %s", parseErr.File)
    if parseErr.Line > 0 {
        fmt.Printf(" at line %d", parseErr.Line)
    }
    fmt.Printf(": %s\n", parseErr.Message)
}
```

---

## Usage Examples

### Example 1: Parse and Display Position

```go
package main

import (
    "fmt"
    "log"
    "github.com/kevung/bgfparser"
)

func main() {
    pos, err := bgfparser.ParseTXT("position.txt")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s vs %s\n", pos.PlayerX, pos.PlayerO)
    fmt.Printf("Score: %d-%d in %d point match\n", 
        pos.ScoreX, pos.ScoreO, pos.MatchLength)
    fmt.Printf("%s to move %d-%d\n", 
        pos.OnRoll, pos.Dice[0], pos.Dice[1])
}
```

### Example 2: Analyze Evaluations

```go
package main

import (
    "fmt"
    "log"
    "github.com/kevung/bgfparser"
)

func main() {
    pos, err := bgfparser.ParseTXT("position.txt")
    if err != nil {
        log.Fatal(err)
    }

    if len(pos.Evaluations) == 0 {
        fmt.Println("No evaluations found")
        return
    }

    best := pos.Evaluations[0]
    fmt.Printf("Best move: %s\n", best.Move)
    fmt.Printf("Equity: %.3f\n", best.Equity)
    
    if len(pos.Evaluations) > 1 {
        second := pos.Evaluations[1]
        blunder := best.Equity - second.Equity
        fmt.Printf("Playing second-best loses %.3f equity\n", blunder)
    }
}
```

### Example 3: Handle Cube Decisions

```go
package main

import (
    "fmt"
    "log"
    "github.com/kevung/bgfparser"
)

func main() {
    pos, err := bgfparser.ParseTXT("position.txt")
    if err != nil {
        log.Fatal(err)
    }

    if pos.CubeDecision == nil {
        fmt.Println("No cube decision in this position")
        return
    }

    cd := pos.CubeDecision
    fmt.Printf("Cube Action: %s\n", cd.Action)
    fmt.Printf("MWC: %.1f%%\n", cd.MWC * 100)
    
    if cd.IsBest {
        fmt.Println("This is the best cube action")
    }
}
```

### Example 4: Batch Processing

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/kevung/bgfparser"
)

func main() {
    files, _ := filepath.Glob("positions/*.txt")
    
    blunderCount := 0
    for _, file := range files {
        pos, err := bgfparser.ParseTXT(file)
        if err != nil {
            log.Printf("Error parsing %s: %v\n", file, err)
            continue
        }
        
        // Check for significant blunders
        if len(pos.Evaluations) >= 2 {
            diff := pos.Evaluations[0].Equity - pos.Evaluations[1].Equity
            if diff > 0.100 { // Blunder threshold
                blunderCount++
                fmt.Printf("Blunder in %s: %.3f\n", file, diff)
            }
        }
    }
    
    fmt.Printf("Total blunders: %d\n", blunderCount)
}
```

---

## Error Handling

The package uses standard Go error handling with a custom `ParseError` type for detailed error information.

### Checking for ParseError

```go
pos, err := bgfparser.ParseTXT("file.txt")
if err != nil {
    if parseErr, ok := err.(*bgfparser.ParseError); ok {
        // Handle parse error with detailed info
        fmt.Printf("Parse error: %s\n", parseErr.Message)
    } else {
        // Handle other errors (file not found, etc.)
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

### Common Errors

- **File not found**: Standard `os.PathError`
- **Invalid format**: `ParseError` with description
- **SMILE encoding**: `ParseError` indicating SMILE not supported
- **Malformed data**: `ParseError` with line number

---

## Performance Notes

- **TXT parsing**: Typically 1-2ms per file
- **BGF parsing**: Typically 10-50ms per file (depends on compression)
- **Memory usage**: ~500 bytes base + ~100-200 bytes per evaluation
- **Thread safety**: Parsers are safe for concurrent use

---

## Limitations

1. **Board state extraction**: Partial implementation from ASCII art
2. **Language support**: Tested with English and French
3. **Statistics**: Some evaluation statistics may not be fully extracted


## See Also

- [README.md](../README.md) - Package overview and quick start
- [PACKAGE_DOCUMENTATION.md](PACKAGE_DOCUMENTATION.md) - Detailed design and usage patterns
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide
