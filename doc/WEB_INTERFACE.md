# Web Interface for BGF Parser

## Overview

The bgfparser package provides flexible parsing functions that work seamlessly with web applications, HTTP uploads, in-memory buffers, and any `io.Reader` source.

## Key Features

- **Flexible Input**: Parse from files, HTTP uploads, memory buffers, or any `io.Reader`
- **JSON Serializable**: All structures have JSON tags for easy API responses
- **Web-Ready**: Designed for integration with web backends and REST APIs
- **No File System Required**: Parse directly from uploaded data without temp files

## API Functions

### BGF Parsing

#### ParseBGF (File-based)
```go
func ParseBGF(filename string) (*Match, error)
```
Parse a BGF file from disk. Use this for command-line tools or file-based workflows.

#### ParseBGFFromReader (Web-ready)
```go
func ParseBGFFromReader(reader io.Reader) (*Match, error)
```
Parse a BGF file from any `io.Reader` source - HTTP uploads, memory buffers, network streams, etc.

### TXT Parsing

#### ParseTXT (File-based)
```go
func ParseTXT(filename string) (*Position, error)
```
Parse a TXT position file from disk.

#### ParseTXTFromReader (Web-ready)
```go
func ParseTXTFromReader(reader io.Reader) (*Position, error)
```
Parse a TXT position file from any `io.Reader` source.

### JSON Export

Both `Match` and `Position` types have `ToJSON()` methods:

```go
func (m *Match) ToJSON() ([]byte, error)
func (p *Position) ToJSON() ([]byte, error)
```

## Usage Examples

### 1. HTTP File Upload Handler

```go
package main

import (
    "net/http"
    "encoding/json"
    "io"
    
    "github.com/unger/bgfparser"
)

func uploadBGFHandler(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    r.ParseMultipartForm(10 << 20) // 10 MB max
    
    // Get the uploaded file
    file, _, err := r.FormFile("bgffile")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Parse directly from upload (no temp file needed)
    match, err := bgfparser.ParseBGFFromReader(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(match)
}
```

### 2. Parse from Memory Buffer

```go
package main

import (
    "bytes"
    "github.com/unger/bgfparser"
)

func parseFromMemory(data []byte) (*Match, error) {
    reader := bytes.NewReader(data)
    return bgfparser.ParseBGFFromReader(reader)
}
```

### 3. Parse TXT from HTTP Upload

```go
func uploadTXTHandler(w http.ResponseWriter, r *http.Request) {
    file, _, _ := r.FormFile("txtfile")
    defer file.Close()
    
    pos, err := bgfparser.ParseTXTFromReader(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Get JSON output
    jsonData, _ := pos.ToJSON()
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}
```

### 4. REST API with Summary Endpoint

```go
type PositionSummary struct {
    PlayerX     string  `json:"player_x"`
    PlayerO     string  `json:"player_o"`
    MatchLength int     `json:"match_length"`
    BestMove    string  `json:"best_move"`
    BestEquity  float64 `json:"best_equity"`
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
    file, _, _ := r.FormFile("txtfile")
    defer file.Close()
    
    pos, err := bgfparser.ParseTXTFromReader(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Create summary
    summary := PositionSummary{
        PlayerX:     pos.PlayerX,
        PlayerO:     pos.PlayerO,
        MatchLength: pos.MatchLength,
    }
    
    // Find best move
    for _, eval := range pos.Evaluations {
        if eval.IsBest {
            summary.BestMove = eval.Move
            summary.BestEquity = eval.Equity
            break
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(summary)
}
```

### 5. Stream Processing

```go
import (
    "io"
    "net/http"
)

func downloadAndParse(url string) (*Match, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Parse directly from HTTP response body
    return bgfparser.ParseBGFFromReader(resp.Body)
}
```

## Complete Web Server Example

See `examples/web_server/main.go` for a complete web server implementation that includes:

- HTML form interface for file uploads
- Summary endpoints for quick analysis
- Full JSON endpoints for complete data
- Both BGF and TXT file support
- Health check endpoint

To run the example:

```bash
cd examples/web_server
go run main.go
```

Then visit http://localhost:8080 in your browser.

## API Endpoints Pattern

A typical web backend might expose these endpoints:

```
POST /api/bgf/parse      - Upload BGF file, return full match JSON
POST /api/bgf/summary    - Upload BGF file, return summary
POST /api/txt/parse      - Upload TXT file, return full position JSON
POST /api/txt/summary    - Upload TXT file, return summary
POST /api/txt/evaluate   - Upload TXT file, return evaluations
GET  /health             - Health check
```

## Database Integration

The JSON-serializable structures can be easily stored in databases:

### PostgreSQL JSONB Example

```sql
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255),
    uploaded_at TIMESTAMP DEFAULT NOW(),
    match_data JSONB
);

INSERT INTO matches (filename, match_data)
VALUES ($1, $2);
```

```go
func saveMatchToDB(db *sql.DB, filename string, match *bgfparser.Match) error {
    jsonData, _ := match.ToJSON()
    _, err := db.Exec(
        "INSERT INTO matches (filename, match_data) VALUES ($1, $2)",
        filename, jsonData,
    )
    return err
}
```

### MongoDB Example

```go
import "go.mongodb.org/mongo-driver/mongo"

func saveMatchToMongo(collection *mongo.Collection, match *bgfparser.Match) error {
    _, err := collection.InsertOne(context.TODO(), match)
    return err
}
```

## Data Structures

All structures are JSON-serializable with appropriate tags:

### Match Structure
```go
type Match struct {
    Format   string                 `json:"format"`
    Version  string                 `json:"version"`
    Compress bool                   `json:"compress"`
    UseSmile bool                   `json:"useSmile"`
    Data     map[string]interface{} `json:"data,omitempty"`
    DecodingWarning string          `json:"decoding_warning,omitempty"`
}
```

### Position Structure
```go
type Position struct {
    Board        [26]int           `json:"board"`
    PlayerX      string            `json:"player_x"`
    PlayerO      string            `json:"player_o"`
    ScoreX       int               `json:"score_x"`
    ScoreO       int               `json:"score_o"`
    MatchLength  int               `json:"match_length"`
    PositionID   string            `json:"position_id"`
    MatchID      string            `json:"match_id"`
    XGID         string            `json:"xgid"`
    OnRoll       string            `json:"on_roll"`
    Dice         [2]int            `json:"dice"`
    CubeValue    int               `json:"cube_value"`
    Evaluations  []Evaluation      `json:"evaluations,omitempty"`
    CubeDecision *CubeDecision     `json:"cube_decision,omitempty"`
    // ... other fields
}
```

## Performance Considerations

- **Memory Efficiency**: Parse directly from streams without buffering entire file
- **No Temp Files**: Avoid file system overhead
- **Streaming**: For large files, data is processed as it's read
- **JSON Encoding**: Fast marshaling with standard library

## Error Handling

All parsing functions return detailed errors:

```go
match, err := bgfparser.ParseBGFFromReader(reader)
if err != nil {
    if parseErr, ok := err.(*bgfparser.ParseError); ok {
        // Handle specific parse error
        log.Printf("Parse error: %s", parseErr.Message)
    }
    return err
}
```

## Testing

Example test for HTTP handler:

```go
func TestUploadHandler(t *testing.T) {
    // Create test file
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    
    // Read test BGF file
    file, _ := os.Open("testdata/match.bgf")
    part, _ := writer.CreateFormFile("bgffile", "match.bgf")
    io.Copy(part, file)
    writer.Close()
    
    // Create request
    req := httptest.NewRequest("POST", "/upload", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    // Test handler
    w := httptest.NewRecorder()
    uploadHandler(w, req)
    
    // Verify response
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Security Considerations

1. **File Size Limits**: Always set max upload size
   ```go
   r.ParseMultipartForm(10 << 20) // 10 MB
   ```

2. **File Type Validation**: Check file extensions/content
   ```go
   if !strings.HasSuffix(header.Filename, ".bgf") {
       http.Error(w, "Invalid file type", http.StatusBadRequest)
       return
   }
   ```

3. **Input Validation**: Validate parsed data before use
4. **Rate Limiting**: Implement rate limiting for upload endpoints
5. **Timeouts**: Set reasonable timeouts for parsing operations

## Next Steps

- See `examples/web_server/main.go` for a complete working example
- Check `doc/API_REFERENCE.md` for detailed API documentation
- Review test files for more usage patterns
