# bgfparser Web Architecture

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Client Applications                      │
├─────────────┬─────────────┬─────────────┬──────────────────┤
│   Web UI    │  Mobile App │  REST API   │  Command Line    │
│   Browser   │             │   Client    │     Tools        │
└─────┬───────┴─────┬───────┴──────┬──────┴────────┬─────────┘
      │             │              │               │
      │ HTTP POST   │ HTTP POST    │ HTTP GET/POST │ File Path
      │ multipart   │ multipart    │ JSON          │
      │             │              │               │
┌─────▼─────────────▼──────────────▼───────────────▼─────────┐
│                     Web Server Layer                        │
│                  (examples/web_server)                      │
├─────────────────────────────────────────────────────────────┤
│  Handlers:                                                  │
│  • uploadBGFHandler()    - Upload BGF summary              │
│  • fullBGFHandler()      - Full BGF JSON                   │
│  • uploadTXTHandler()    - Upload TXT summary              │
│  • fullTXTHandler()      - Full TXT JSON                   │
│  • healthHandler()       - Health check                    │
└─────────┬─────────────────────┬─────────────────────────────┘
          │                     │
          │ io.Reader           │ Filename string
          │ (HTTP upload)       │ (File path)
          │                     │
┌────────────────────────────────────────────────────────┐
│                   BGFParser Package                    │
│                   (github.com/kevung/bgfparser)              │
└────────────────────────────────────────────────────────┘
│  Web-Ready API:                                             │
│  • ParseBGFFromReader(io.Reader) -> *Match                  │
│  • ParseTXTFromReader(io.Reader) -> *Position              │
│                                                              │
│  File-Based API:                                            │
│  • ParseBGF(filename) -> *Match                             │
│  • ParseTXT(filename) -> *Position                          │
│                                                              │
│  JSON Export:                                               │
│  • (*Match).ToJSON() -> []byte                              │
│  • (*Position).ToJSON() -> []byte                           │
└──────────────┬────────────────────┬──────────────────────────┘
               │                    │
               │                    │
┌──────────────▼────────┐  ┌────────▼──────────────────────────┐
│  Input Parsers        │  │  Data Structures                  │
├───────────────────────┤  ├───────────────────────────────────┤
│  • bgf_parser.go      │  │  • types.go                       │
│  • txt_parser.go      │  │    - Position (JSON tags)         │
│  • web.go             │  │    - Match (JSON tags)            │
│  • txt_parser_helpers │  │    - Evaluation (JSON tags)       │
│                       │  │    - CubeDecision (JSON tags)     │
└───────────────────────┘  └───────────────────────────────────┘
               │
               │
┌──────────────▼────────────────────────────────────────────┐
│              Internal Components                          │
├───────────────────────────────────────────────────────────┤
│  • internal/smile - SMILE binary JSON decoder             │
│  • gzip decompression                                     │
│  • JSON parsing                                           │
└───────────────────────────────────────────────────────────┘
```

## Data Flow

### File Upload Flow

```
User Browser
    │
    │ 1. Upload BGF/TXT file via HTML form
    ▼
HTTP Server (port 8080)
    │
    │ 2. r.FormFile("bgffile")
    ▼
ParseBGFFromReader(file)
    │
    │ 3. Read header
    │ 4. Decompress (if needed)
    │ 5. Decode SMILE (if needed)
    │ 6. Parse JSON
    ▼
Match/Position struct
    │
    │ 7. ToJSON()
    ▼
HTTP Response (JSON)
    │
    │ 8. Send to browser
    ▼
User sees JSON data
```

### In-Memory Processing Flow

```
Application
    │
    │ 1. []byte data from any source
    ▼
bytes.NewReader(data)
    │
    │ 2. Create io.Reader
    ▼
ParseBGFFromReader(reader)
    │
    │ 3. Parse directly from memory
    ▼
Match/Position struct
    │
    │ 4. Use in application
    ▼
Database, API, etc.
```

## Input Sources

The parser accepts data from multiple sources via `io.Reader`:

```
┌──────────────────────────────────────────────────────────┐
│                    Input Sources                         │
├──────────────┬──────────────┬──────────────┬────────────┤
│ File Upload  │ HTTP Stream  │ Memory       │ File       │
├──────────────┼──────────────┼──────────────┼────────────┤
│ multipart.   │ net.Conn     │ bytes.       │ os.File    │
│ File         │              │ Reader       │            │
└──────┬───────┴──────┬───────┴──────┬───────┴─────┬──────┘
       │              │              │             │
       └──────────────┴──────────────┴─────────────┘
                      │
                      ▼
              io.Reader interface
                      │
                      ▼
        ParseBGFFromReader(reader)
        ParseTXTFromReader(reader)
```

## Output Formats

The parser provides data in multiple formats:

```
┌─────────────────────────────────────────────────────────┐
│              Parsed Data (Go structs)                   │
└──────────────────┬──────────────────┬───────────────────┘
                   │                  │
         ┌─────────▼─────────┐  ┌─────▼──────────┐
         │   Direct Access   │  │  JSON Export   │
         ├───────────────────┤  ├────────────────┤
         │ pos.PlayerX       │  │ pos.ToJSON()   │
         │ pos.Evaluations   │  │ match.ToJSON() │
         │ match.Data        │  │                │
         └─────────┬─────────┘  └────────┬───────┘
                   │                     │
        ┌──────────▼──────────┐ ┌────────▼────────┐
        │  Go Application     │ │  Web Response   │
        │  • Logic            │ │  • HTTP API     │
        │  • Processing       │ │  • Database     │
        │  • Analysis         │ │  • Storage      │
        └─────────────────────┘ └─────────────────┘
```

## API Endpoints Example

```
                    HTTP Endpoints
    ┌──────────────────┴──────────────────┐
    │                                     │
GET /                                POST /upload/bgf
    │                                     │
    │ HTML Interface                      │ multipart/form-data
    │                                     │
    ▼                                     ▼
┌──────────────────┐              ┌──────────────────┐
│  Web UI Form     │              │ ParseBGFFromReader│
│  - BGF upload    │              │      ↓           │
│  - TXT upload    │              │  Match struct    │
│  - Results view  │              │      ↓           │
└──────────────────┘              │ Summary JSON     │
                                  └──────────────────┘
                                           │
                                           ▼
                                    HTTP Response
                                    {
                                      "format": "BGF",
                                      "version": "1.0",
                                      "match_info": {...}
                                    }
```

## Database Integration Pattern

```
HTTP Upload
    │
    ▼
ParseBGFFromReader()
    │
    ▼
Match struct
    │
    ├─→ ToJSON() ──→ PostgreSQL JSONB
    │                 INSERT INTO matches
    │                 (match_data) VALUES ($1)
    │
    ├─→ Extract fields ──→ SQL Table
    │                      INSERT INTO matches
    │                      (player1, player2, ...)
    │
    └─→ Direct struct ──→ MongoDB
                          collection.InsertOne(match)
```

## Use Case Scenarios

### 1. Web Analysis Service
```
User uploads BGF → Parse → Analyze → Return statistics
```

### 2. Match Database
```
Upload BGF → Parse → Extract → Store in DB → Query later
```

### 3. Real-time API
```
POST /analyze → Parse from request → Compute → JSON response
```

### 4. Batch Processing
```
Directory of files → Parse each → Aggregate stats → Report
```

## Security Layers

```
┌──────────────────────────────────────────────────┐
│              Security Considerations             │
├──────────────────────────────────────────────────┤
│  1. File size limits (10 MB default)            │
│  2. File type validation (.bgf, .txt)           │
│  3. Rate limiting (requests/minute)              │
│  4. Input sanitization                           │
│  5. Error handling (no data leakage)             │
│  6. Timeout limits                               │
└──────────────────────────────────────────────────┘
```

## Performance Characteristics

```
File Upload → Parse → Response
    ↓          ↓         ↓
  <1ms      <10ms     <5ms   (typical match file)
  
Stream:  No buffering, processes as received
Memory:  Minimal allocation, efficient parsing
CPU:     Single-threaded, fast decompression
```

## Key Benefits

1. **No Temporary Files** - Direct parsing from upload
2. **Flexible Input** - Works with any io.Reader
3. **JSON Ready** - All structs serializable
4. **Type Safe** - Go structs with full type info
5. **Database Ready** - Easy PostgreSQL/MongoDB integration
6. **Web Native** - Designed for HTTP handlers
7. **Backward Compatible** - File-based API unchanged

## Next Steps

- Add authentication middleware
- Implement rate limiting
- Add caching layer
- Create Swagger docs
- Add streaming for large files
- Deploy with Docker
