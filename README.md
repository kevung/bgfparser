# BGFParser

Go package for parsing BGBlitz backgammon files (BGF binary matches, TXT positions).

## Features

- **TXT Parser**: Board positions, evaluations, cube decisions, multilingual (EN/FR/DE/JP)
- **BGF Parser**: Binary match files with SMILE encoding support
- **Web-Ready**: Parse from `io.Reader` (HTTP uploads, memory, streams)
- **JSON Export**: Database-ready structures with JSON tags

## Installation

```bash
go get github.com/kevung/bgfparser
```

## Quick Start

```go
// Parse position file
pos, _ := bgfparser.ParseTXT("position.txt")
fmt.Printf("%s vs %s, Score: %d-%d\n", 
    pos.PlayerX, pos.PlayerO, pos.ScoreX, pos.ScoreO)

// Parse from HTTP upload
match, _ := bgfparser.ParseBGFFromReader(uploadedFile)

// Export to JSON
jsonData, _ := pos.ToJSON()
```

## CLI Tools

Build and run command-line tools:

```bash
# Build tools
go build -o bin/parse_txt ./examples/parse_txt/
go build -o bin/web_server ./examples/web_server/

# Parse TXT file
./bin/parse_txt position.txt

# Start web server (http://localhost:8080)
./bin/web_server

# Or run directly
go run examples/parse_txt/main.go position.txt
```

## Data Structures

```go
type Position struct {
    Board       [26]int          // Checker positions (1-24=points, 25=bar)
    PlayerX, PlayerO string       // Player names
    ScoreX, ScoreO   int         // Scores
    MatchLength int              // Match length
    XGID        string           // Position ID
    OnRoll      string           // "X" or "O"
    Dice        [2]int           // Dice roll
    CubeValue   int              // Cube value
    Evaluations []Evaluation     // Move analysis
    CubeDecisions []CubeDecision // Cube analysis (all actions)
    CubelessEquity float64       // Cubeless equity
    CubefulEquity  float64       // Cubeful equity
}

type Evaluation struct {
    Rank   int     // 1 = best move
    Move   string  // Move notation
    Equity float64 // Equity value
    Win, WinG, WinBG float64 // Win probabilities
}

type CubeDecision struct {
    Action  string  // "No Double", "Double/Take", "Double/Pass"
    MWC     float64 // Match winning chances
    EMG     float64 // Effective match gammons
}
```

## Examples

See `examples/` directory:
- `parse_txt/` - Parse TXT positions  
- `parse_bgf/` - Parse BGF matches
- `batch_parse/` - Process multiple files
- `web_server/` - HTTP upload interface

## Documentation

- [API Reference](doc/API_REFERENCE.md) - Complete API docs
- [Quick Reference](doc/QUICK_REFERENCE.md) - Common patterns
- [Extracted Information](doc/EXTRACTED_INFORMATION.md) - All parsed fields
- [Multilingual Support](MULTILINGUAL_SUPPORT.md) - Language support details

## License

MIT License - see [LICENSE](LICENSE) file.

## Acknowledgments

- [Frank Berger](https://www.bgblitz.com/) - BGBlitz creator
- [LeLuxNet](https://gitlab.com/LeLuxNet/X) - SMILE decoder

````
