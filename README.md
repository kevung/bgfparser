# BGFParser

A Go package for parsing BGBlitz backgammon data files, including SMILE-encoded binary formats.

## Overview

BGFParser is a comprehensive library for reading and analyzing backgammon data from BGBlitz, a popular backgammon analysis software. The package supports two file formats:

- **TXT files**: Human-readable position files containing board states, evaluations, and analysis
- **BGF files**: Binary match files using compressed JSON (with optional SMILE encoding)

## Features

### TXT Parser
- ✅ Parse board positions with checker locations
- ✅ Extract player names and scores
- ✅ Read match information (length, Crawford state)
- ✅ Parse position identifiers (Position-ID, Match-ID, XGID)
- ✅ Extract dice rolls and cube information
- ✅ Parse move evaluations with equity and statistics
- ✅ Read cube decisions (Double/Take, Double/Pass, No Double)
- ✅ Support for multiple languages (English, French)
- ✅ Pip count extraction

### BGF Parser
- ✅ Read BGF file headers (format, version, compression info)
- ✅ Decompress gzip-compressed data
- ✅ Parse JSON data from uncompressed files
- ✅ Basic SMILE decoding (extracts player names, dates, and key metadata)
- ⚠️ Full SMILE decoding is partial (complex nested structures may not parse completely)

## Installation

```bash
go get github.com/unger/bgfparser
```

## Quick Start

### Parsing a Position (TXT) File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/unger/bgfparser"
)

func main() {
    // Parse a position file
    position, err := bgfparser.ParseTXT("game_position.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // Access position data
    fmt.Printf("Players: %s vs %s\n", position.PlayerX, position.PlayerO)
    fmt.Printf("Score: %d-%d\n", position.ScoreX, position.ScoreO)
    fmt.Printf("On Roll: %s with %d-%d\n", 
        position.OnRoll, position.Dice[0], position.Dice[1])
    
    // Check evaluations
    if len(position.Evaluations) > 0 {
        best := position.Evaluations[0]
        fmt.Printf("Best move: %s (Equity: %.3f)\n", 
            best.Move, best.Equity)
    }
}
```

### Parsing a Match (BGF) File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/unger/bgfparser"
)

func main() {
    // Parse a BGF match file
    match, err := bgfparser.ParseBGF("match.bgf")
    if err != nil {
        log.Fatal(err)
    }
    
    // Display match information
    fmt.Printf("Format: %s v%s\n", match.Format, match.Version)
    fmt.Printf("Compressed: %v\n", match.Compress)
    fmt.Printf("Uses SMILE: %v\n", match.UseSmile)
    
    // Check for decoding warnings
    if match.DecodingWarning != "" {
        fmt.Printf("Warning: %s\n", match.DecodingWarning)
    }
    
    // Get match metadata
    info := match.GetMatchInfo()
    for key, value := range info {
        fmt.Printf("%s: %v\n", key, value)
    }
}
```

## Data Types

### Position

Represents a backgammon position with complete state:

```go
type Position struct {
    Board       [26]int           // Checker positions (0=unused, 1-24=points, 25=bar)
    PlayerX     string            // Player X name
    PlayerO     string            // Player O name
    ScoreX      int               // Player X score
    ScoreO      int               // Player O score
    MatchLength int               // Match length in points
    Crawford    bool              // Crawford rule active
    PositionID  string            // BGBlitz Position-ID
    MatchID     string            // BGBlitz Match-ID
    XGID        string            // XG format ID
    OnRoll      string            // "X" or "O"
    Dice        [2]int            // Current dice roll
    CubeValue   int               // Doubling cube value
    CubeOwner   string            // "", "X", or "O"
    OnBar       map[string]int    // Checkers on bar
    PipCount    map[string]int    // Pip counts
    Evaluations []Evaluation      // Move evaluations
    CubeDecision *CubeDecision    // Cube decision analysis
}
```

### Evaluation

Move evaluation with equity and statistics:

```go
type Evaluation struct {
    Rank    int      // Move rank (1 = best)
    Move    string   // Move notation
    Equity  float64  // Equity value
    Diff    float64  // Difference from best
    Win     float64  // Win probability
    WinG    float64  // Gammon win probability
    WinBG   float64  // Backgammon win probability
    LoseG   float64  // Gammon loss probability
    LoseBG  float64  // Backgammon loss probability
    IsBest  bool     // True if marked as best move
}
```

### CubeDecision

Cube action analysis:

```go
type CubeDecision struct {
    Action  string  // "Double/Take", "Double/Pass", or "No Double"
    MWC     float64 // Match Winning Chances
    MWCDiff float64 // Difference in MWC
    EMG     float64 // EMG (normalized equity)
    EMGDiff float64 // Difference in EMG
    IsBest  bool    // True if marked as best action
}
```

### Match

BGF match file data:

```go
type Match struct {
    Format     string                 // File format ("BGF")
    Version    string                 // Format version
    Compressed bool                   // Gzip compression used
    UseSmile   bool                   // SMILE encoding used
    Data       map[string]interface{} // Parsed match data
    DecodingWarning string            // Warning if decoding was incomplete
}
```

## Examples

The package includes three complete example programs in the `examples/` directory:

### 1. parse_txt
Detailed parser for TXT position files showing all available data.

```bash
go run examples/parse_txt/main.go tmp/blunder21_EN.txt
```

### 2. parse_bgf
Parser for BGF match files with format detection.

```bash
go run examples/parse_bgf/main.go tmp/match.bgf
```

### 3. batch_parse
Batch processor for parsing all files in a directory.

```bash
go run examples/batch_parse/main.go tmp/
```

## File Format Details

### TXT Format

Text files contain ASCII art representation of the board plus analysis:

```
 +13-14-15-16-17-18------19-20-21-22-23-24-+   O: TachiAI_V  161
 | X           O    |   | O     O        X |
 | X           O    |   | O     O        X |
 ...
 Position-ID: mGfwATDgc/ABMA    Match-ID: cAllAAAAAAAE
 XGID=-b----E-C---eE---b-d-b--B-:0:0:1:21:0:0:0:3:10
 
 TachiAI_V - 0 player - 0 in a 3 point match.
 player to move 2-1
 
 Evaluation
 ==========
 1) 13-11 24-23                0.473 / -0.289
    0.443  0.113  0.002  -  0.557  0.179  0.003 
```

The parser extracts:
- Board state (visual representation)
- Player information
- Match state and score
- Position identifiers
- Move evaluations with equity values
- Win/loss probabilities
- Cube decisions

### BGF Format

Binary files with two parts:

1. **Header** (JSON, uncompressed):
```json
{"format":"BGF","version":"1.0","compress":true,"useSmile":true}
```

2. **Data** (gzip + SMILE encoded JSON):
   - Complete match data
   - All positions
   - Move history
   - Analysis results

**Note**: SMILE decoding is implemented with basic support. The parser can read SMILE-encoded BGF files and extract player names, dates, match parameters, and other metadata. Complex nested structures may not parse completely, but key information is accessible.

## Project Structure

```
bgfparser/
├── LICENSE              # MIT License
├── README.md           # This file
├── go.mod              # Go module definition
├── types.go            # Core data types
├── txt_parser.go       # TXT format parser
├── bgf_parser.go       # BGF format parser
├── doc/
│   └── BGF_format.md   # BGF format documentation
├── examples/
│   ├── parse_txt/      # TXT parser example
│   ├── parse_bgf/      # BGF parser example
│   └── batch_parse/    # Batch processing example
└── tmp/                # Sample data files
    ├── *.txt           # Position files
    └── *.bgf           # Match files
```

## Testing

Run the batch parser on the sample data:

```bash
# Build examples
go build -o bin/parse_txt ./examples/parse_txt/
go build -o bin/parse_bgf ./examples/parse_bgf/
go build -o bin/batch_parse ./examples/batch_parse/

# Test on sample data
./bin/batch_parse tmp/
```

## Limitations and Future Work

### Current Limitations
- **Board parsing**: Board state extraction from ASCII art is partially implemented
- **SMILE decoding**: Complex nested structures in SMILE format may not fully decode
- **Language support**: Primarily tested with English and French files
- **Incomplete statistics**: Some evaluation statistics may not be fully extracted

### Planned Improvements
- Enhanced SMILE decoder for complex nested structures
- Complete board state parsing from ASCII representation
- Support for additional BGBlitz output formats
- More comprehensive test suite
- Performance optimizations for batch processing

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for:
- Bug fixes
- Additional features
- Documentation improvements
- Test cases
- Language support

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- BGBlitz for creating the analysis software
- The backgammon community for the XGID and position notation standards

## Support

For questions, issues, or suggestions:
- Open an issue on GitHub
- Check the `doc/` directory for format specifications
- Review the `examples/` directory for usage patterns
