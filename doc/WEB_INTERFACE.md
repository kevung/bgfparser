# Web Integration Guide# Web Interface for BGF Parser



## Overview## Overview



BGFParser provides io.Reader-based functions for seamless web integration - parse HTTP uploads, memory buffers, and network streams without temporary files.The bgfparser package provides flexible parsing functions that work seamlessly with web applications, HTTP uploads, in-memory buffers, and any `io.Reader` source.



## Core Functions## Key Features



```go- **Flexible Input**: Parse from files, HTTP uploads, memory buffers, or any `io.Reader`

// Parse from any io.Reader source- **JSON Serializable**: All structures have JSON tags for easy API responses

func ParseBGFFromReader(reader io.Reader) (*Match, error)- **Web-Ready**: Designed for integration with web backends and REST APIs

func ParseTXTFromReader(reader io.Reader) (*Position, error)- **No File System Required**: Parse directly from uploaded data without temp files



// Export to JSON## API Functions

func (m *Match) ToJSON() ([]byte, error)

func (p *Position) ToJSON() ([]byte, error)### BGF Parsing

```

#### ParseBGF (File-based)

## HTTP Upload Handler```go

func ParseBGF(filename string) (*Match, error)

### Basic Example```

Parse a BGF file from disk. Use this for command-line tools or file-based workflows.

```go

import (#### ParseBGFFromReader (Web-ready)

    "encoding/json"```go

    "net/http"func ParseBGFFromReader(reader io.Reader) (*Match, error)

    "github.com/kevung/bgfparser"```

)Parse a BGF file from any `io.Reader` source - HTTP uploads, memory buffers, network streams, etc.



func uploadHandler(w http.ResponseWriter, r *http.Request) {### TXT Parsing

    file, _, err := r.FormFile("bgffile")

    if err != nil {#### ParseTXT (File-based)

        http.Error(w, err.Error(), http.StatusBadRequest)```go

        returnfunc ParseTXT(filename string) (*Position, error)

    }```

    defer file.Close()Parse a TXT position file from disk.

    

    match, err := bgfparser.ParseBGFFromReader(file)#### ParseTXTFromReader (Web-ready)

    if err != nil {```go

        http.Error(w, err.Error(), http.StatusBadRequest)func ParseTXTFromReader(reader io.Reader) (*Position, error)

        return```

    }Parse a TXT position file from any `io.Reader` source.

    

    w.Header().Set("Content-Type", "application/json")### JSON Export

    json.NewEncoder(w).Encode(match)

}Both `Match` and `Position` types have `ToJSON()` methods:

```

```go

### Custom Responsefunc (m *Match) ToJSON() ([]byte, error)

func (p *Position) ToJSON() ([]byte, error)

```go```

func summaryHandler(w http.ResponseWriter, r *http.Request) {

    file, _, _ := r.FormFile("txtfile")## Usage Examples

    defer file.Close()

    ### 1. HTTP File Upload Handler

    pos, err := bgfparser.ParseTXTFromReader(file)

    if err != nil {```go

        http.Error(w, err.Error(), http.StatusBadRequest)package main

        return

    }import (

        "net/http"

    summary := map[string]interface{}{    "encoding/json"

        "players": pos.PlayerX + " vs " + pos.PlayerO,    "io"

        "score": fmt.Sprintf("%d-%d/%d", pos.ScoreX, pos.ScoreO, pos.MatchLength),    

        "evaluations": len(pos.Evaluations),    "github.com/kevung/bgfparser"

        "best_move": pos.Evaluations[0].Move,)

        "equity": pos.Evaluations[0].Equity,

    }func uploadBGFHandler(w http.ResponseWriter, r *http.Request) {

        // Parse multipart form

    json.NewEncoder(w).Encode(summary)    r.ParseMultipartForm(10 << 20) // 10 MB max

}    

```    // Get the uploaded file

    file, _, err := r.FormFile("bgffile")

## Memory Buffer Parsing    if err != nil {

        http.Error(w, err.Error(), http.StatusBadRequest)

```go        return

import "bytes"    }

    defer file.Close()

// Parse from byte slice    

data := []byte{...} // BGF file contents    // Parse directly from upload (no temp file needed)

reader := bytes.NewReader(data)    match, err := bgfparser.ParseBGFFromReader(file)

match, err := bgfparser.ParseBGFFromReader(reader)    if err != nil {

```        http.Error(w, err.Error(), http.StatusBadRequest)

        return

## Full Web Server Example    }

    

```go    // Return JSON response

package main    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(match)

import (}

    "encoding/json"```

    "log"

    "net/http"### 2. Parse from Memory Buffer

    "github.com/kevung/bgfparser"

)```go

package main

func main() {

    http.HandleFunc("/upload/bgf", handleBGF)import (

    http.HandleFunc("/upload/txt", handleTXT)    "bytes"

    log.Fatal(http.ListenAndServe(":8080", nil))    "github.com/kevung/bgfparser"

})



func handleBGF(w http.ResponseWriter, r *http.Request) {func parseFromMemory(data []byte) (*Match, error) {

    if r.Method != "POST" {    reader := bytes.NewReader(data)

        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)    return bgfparser.ParseBGFFromReader(reader)

        return}

    }```

    

    file, header, err := r.FormFile("file")### 3. Parse TXT from HTTP Upload

    if err != nil {

        http.Error(w, err.Error(), http.StatusBadRequest)```go

        returnfunc uploadTXTHandler(w http.ResponseWriter, r *http.Request) {

    }    file, _, _ := r.FormFile("txtfile")

    defer file.Close()    defer file.Close()

        

    match, err := bgfparser.ParseBGFFromReader(file)    pos, err := bgfparser.ParseTXTFromReader(file)

    if err != nil {    if err != nil {

        http.Error(w, "Parse error: "+err.Error(), http.StatusBadRequest)        http.Error(w, err.Error(), http.StatusBadRequest)

        return        return

    }    }

        

    response := map[string]interface{}{    // Get JSON output

        "filename": header.Filename,    jsonData, _ := pos.ToJSON()

        "format":   match.Format,    w.Header().Set("Content-Type", "application/json")

        "version":  match.Version,    w.Write(jsonData)

        "data":     match.Data,}

    }```

    

    w.Header().Set("Content-Type", "application/json")### 4. REST API with Summary Endpoint

    json.NewEncoder(w).Encode(response)

}```go

type PositionSummary struct {

func handleTXT(w http.ResponseWriter, r *http.Request) {    PlayerX     string  `json:"player_x"`

    if r.Method != "POST" {    PlayerO     string  `json:"player_o"`

        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)    MatchLength int     `json:"match_length"`

        return    BestMove    string  `json:"best_move"`

    }    BestEquity  float64 `json:"best_equity"`

    }

    file, _, err := r.FormFile("file")

    if err != nil {func summaryHandler(w http.ResponseWriter, r *http.Request) {

        http.Error(w, err.Error(), http.StatusBadRequest)    file, _, _ := r.FormFile("txtfile")

        return    defer file.Close()

    }    

    defer file.Close()    pos, err := bgfparser.ParseTXTFromReader(file)

        if err != nil {

    pos, err := bgfparser.ParseTXTFromReader(file)        http.Error(w, err.Error(), http.StatusBadRequest)

    if err != nil {        return

        http.Error(w, "Parse error: "+err.Error(), http.StatusBadRequest)    }

        return    

    }    // Create summary

        summary := PositionSummary{

    w.Header().Set("Content-Type", "application/json")        PlayerX:     pos.PlayerX,

    json.NewEncoder(w).Encode(pos)        PlayerO:     pos.PlayerO,

}        MatchLength: pos.MatchLength,

```    }

    

## Database Storage    // Find best move

    for _, eval := range pos.Evaluations {

### PostgreSQL JSONB        if eval.IsBest {

            summary.BestMove = eval.Move

```go            summary.BestEquity = eval.Equity

import "database/sql"            break

        }

func savePosition(db *sql.DB, pos *bgfparser.Position) error {    }

    jsonData, err := pos.ToJSON()    

    if err != nil {    w.Header().Set("Content-Type", "application/json")

        return err    json.NewEncoder(w).Encode(summary)

    }}

    ```

    _, err = db.Exec(

        "INSERT INTO positions (position_data, player_x, player_o) VALUES ($1, $2, $3)",### 5. Stream Processing

        jsonData, pos.PlayerX, pos.PlayerO,

    )```go

    return errimport (

}    "io"

```    "net/http"

)

### Query Example

func downloadAndParse(url string) (*Match, error) {

```sql    resp, err := http.Get(url)

CREATE TABLE positions (    if err != nil {

    id SERIAL PRIMARY KEY,        return nil, err

    position_data JSONB,    }

    player_x TEXT,    defer resp.Body.Close()

    player_o TEXT,    

    created_at TIMESTAMP DEFAULT NOW()    // Parse directly from HTTP response body

);    return bgfparser.ParseBGFFromReader(resp.Body)

}

-- Find positions with high equity```

SELECT player_x, player_o, position_data->>'xgid' as xgid

FROM positions## Complete Web Server Example

WHERE (position_data->'evaluations'->0->>'equity')::float > 0.5;

```See `examples/web_server/main.go` for a complete web server implementation that includes:



## Error Handling- HTML form interface for file uploads

- Summary endpoints for quick analysis

```go- Full JSON endpoints for complete data

func handleUpload(w http.ResponseWriter, r *http.Request) {- Both BGF and TXT file support

    file, _, err := r.FormFile("file")- Health check endpoint

    if err != nil {

        sendError(w, "Upload failed", http.StatusBadRequest)To run the example:

        return

    }```bash

    defer file.Close()cd examples/web_server

    go run main.go

    pos, err := bgfparser.ParseTXTFromReader(file)```

    if err != nil {

        if parseErr, ok := err.(*bgfparser.ParseError); ok {Then visit http://localhost:8080 in your browser.

            sendError(w, fmt.Sprintf("Parse error: %s", parseErr.Message), 

                     http.StatusBadRequest)## API Endpoints Pattern

        } else {

            sendError(w, "Internal error", http.StatusInternalServerError)A typical web backend might expose these endpoints:

        }

        return```

    }POST /api/bgf/parse      - Upload BGF file, return full match JSON

    POST /api/bgf/summary    - Upload BGF file, return summary

    w.Header().Set("Content-Type", "application/json")POST /api/txt/parse      - Upload TXT file, return full position JSON

    json.NewEncoder(w).Encode(pos)POST /api/txt/summary    - Upload TXT file, return summary

}POST /api/txt/evaluate   - Upload TXT file, return evaluations

GET  /health             - Health check

func sendError(w http.ResponseWriter, message string, code int) {```

    w.Header().Set("Content-Type", "application/json")

    w.WriteHeader(code)## Database Integration

    json.NewEncoder(w).Encode(map[string]string{"error": message})

}The JSON-serializable structures can be easily stored in databases:

```

### PostgreSQL JSONB Example

## Running the Example

```sql

See `examples/web_server/main.go` for a complete implementation:CREATE TABLE matches (

    id SERIAL PRIMARY KEY,

```bash    filename VARCHAR(255),

cd examples/web_server    uploaded_at TIMESTAMP DEFAULT NOW(),

go run main.go    match_data JSONB

# Server starts on http://localhost:8080);

```

INSERT INTO matches (filename, match_data)

## See AlsoVALUES ($1, $2);

```

- [Quick Reference](QUICK_REFERENCE.md) - Common patterns

- [API Reference](API_REFERENCE.md) - Complete API```go

- [Package Documentation](PACKAGE_DOCUMENTATION.md) - Design and architecturefunc saveMatchToDB(db *sql.DB, filename string, match *bgfparser.Match) error {

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
