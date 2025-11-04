# bgfparser - Quick API Reference

## Installation

```bash
go get github.com/unger/bgfparser
```

## File-Based Parsing

```go
import "github.com/unger/bgfparser"

// Parse BGF match file
match, err := bgfparser.ParseBGF("match.bgf")

// Parse TXT position file
pos, err := bgfparser.ParseTXT("position.txt")
```

## Web-Ready Parsing (io.Reader)

```go
// Parse from HTTP upload
match, err := bgfparser.ParseBGFFromReader(file)
pos, err := bgfparser.ParseTXTFromReader(file)

// Parse from memory
reader := bytes.NewReader(data)
match, err := bgfparser.ParseBGFFromReader(reader)
```

## JSON Export

```go
// Convert to JSON
jsonData, err := match.ToJSON()
jsonData, err := pos.ToJSON()

// Use directly in HTTP response
w.Header().Set("Content-Type", "application/json")
w.Write(jsonData)
```

## HTTP Upload Handler

```go
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("bgffile")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    match, err := bgfparser.ParseBGFFromReader(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(match)
}
```

## Data Structures

### Match (BGF files)

```go
type Match struct {
    Format   string                 `json:"format"`
    Version  string                 `json:"version"`
    Compress bool                   `json:"compress"`
    UseSmile bool                   `json:"useSmile"`
    Data     map[string]interface{} `json:"data,omitempty"`
}

// Get match info
info := match.GetMatchInfo()
```

### Position (TXT files)

```go
type Position struct {
    PlayerX      string            `json:"player_x"`
    PlayerO      string            `json:"player_o"`
    ScoreX       int               `json:"score_x"`
    ScoreO       int               `json:"score_o"`
    MatchLength  int               `json:"match_length"`
    OnRoll       string            `json:"on_roll"`
    Dice         [2]int            `json:"dice"`
    CubeValue    int               `json:"cube_value"`
    PositionID   string            `json:"position_id"`
    XGID         string            `json:"xgid"`
    Evaluations  []Evaluation      `json:"evaluations,omitempty"`
    CubeDecision *CubeDecision     `json:"cube_decision,omitempty"`
}
```

### Evaluation

```go
type Evaluation struct {
    Rank   int     `json:"rank"`
    Move   string  `json:"move"`
    Equity float64 `json:"equity"`
    Diff   float64 `json:"diff"`
    Win    float64 `json:"win"`
    IsBest bool    `json:"is_best"`
}

// Find best move
for _, eval := range pos.Evaluations {
    if eval.IsBest {
        fmt.Printf("Best: %s (%.3f)\n", eval.Move, eval.Equity)
    }
}
```

## Common Patterns

### Upload + Summary Response

```go
func summaryHandler(w http.ResponseWriter, r *http.Request) {
    file, _, _ := r.FormFile("txtfile")
    defer file.Close()
    
    pos, err := bgfparser.ParseTXTFromReader(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    summary := map[string]interface{}{
        "players": fmt.Sprintf("%s vs %s", pos.PlayerX, pos.PlayerO),
        "score": fmt.Sprintf("%d-%d", pos.ScoreX, pos.ScoreO),
        "match_length": pos.MatchLength,
        "on_roll": pos.OnRoll,
        "num_evals": len(pos.Evaluations),
    }
    
    json.NewEncoder(w).Encode(summary)
}
```

### Database Storage (PostgreSQL JSONB)

```go
func savePosition(db *sql.DB, pos *bgfparser.Position) error {
    jsonData, _ := pos.ToJSON()
    _, err := db.Exec(
        "INSERT INTO positions (position_data) VALUES ($1)",
        jsonData,
    )
    return err
}
```

### Batch Processing

```go
func processDirectory(dir string) error {
    files, _ := filepath.Glob(dir + "/*.txt")
    
    for _, file := range files {
        pos, err := bgfparser.ParseTXT(file)
        if err != nil {
            log.Printf("Error parsing %s: %v", file, err)
            continue
        }
        
        // Process position...
    }
    
    return nil
}
```

## Error Handling

```go
pos, err := bgfparser.ParseTXT("position.txt")
if err != nil {
    if parseErr, ok := err.(*bgfparser.ParseError); ok {
        log.Printf("Parse error in %s: %s", parseErr.File, parseErr.Message)
    }
    return err
}
```

## Run Web Server Example

```bash
cd examples/web_server
go run main.go
# Visit http://localhost:8080
```

## Documentation

- [Full Web Interface Guide](WEB_INTERFACE.md)
- [API Reference](API_REFERENCE.md)
- [README](../README.md)
