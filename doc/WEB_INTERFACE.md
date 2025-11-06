# Web Integration Guide

## Overview

Parse HTTP uploads, memory buffers, and network streams without temporary files using `io.Reader`-based functions.

## API Functions

```go
// Parse from any io.Reader
func ParseBGFFromReader(reader io.Reader) (*Match, error)
func ParseTXTFromReader(reader io.Reader) (*Position, error)

// Export to JSON
func (m *Match) ToJSON() ([]byte, error)
func (p *Position) ToJSON() ([]byte, error)
```

## HTTP Upload Handler

```go
import (
    "encoding/json"
    "net/http"
    "github.com/kevung/bgfparser"
)

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

## Custom Response

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
        "num_evals": len(pos.Evaluations),
    }
    
    json.NewEncoder(w).Encode(summary)
}
```

## Database Integration

```go
// PostgreSQL JSONB
func savePosition(db *sql.DB, pos *bgfparser.Position) error {
    jsonData, _ := pos.ToJSON()
    _, err := db.Exec(
        "INSERT INTO positions (position_data) VALUES ($1)",
        jsonData,
    )
    return err
}

// MongoDB
func saveToMongo(coll *mongo.Collection, pos *bgfparser.Position) error {
    var doc map[string]interface{}
    jsonData, _ := pos.ToJSON()
    json.Unmarshal(jsonData, &doc)
    _, err := coll.InsertOne(context.Background(), doc)
    return err
}
```

## Web Server Example

```bash
cd examples/web_server
go run main.go
# Visit http://localhost:8080
```

## Full Example

See `examples/web_server/main.go` for a complete web server with:
- BGF/TXT upload handlers
- Summary and full JSON endpoints
- Health check
- HTML upload interface
