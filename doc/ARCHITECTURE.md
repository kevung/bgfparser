# bgfparser Architecture

## System Overview

```
Client → Web Server → bgfparser Package → Parsers → Data Structures
```

## Data Flow

### HTTP Upload
```
HTTP Upload → ParseBGFFromReader(io.Reader) → *Match → ToJSON() → JSON Response
HTTP Upload → ParseTXTFromReader(io.Reader) → *Position → ToJSON() → JSON Response
```

### File-Based
```
File Path → ParseBGF(filename) → *Match → Data Structure
File Path → ParseTXT(filename) → *Position → Data Structure
```

## Package Structure

```
bgfparser/
├── types.go              # Data structures (Match, Position, Evaluation, etc.)
├── bgf_parser.go         # BGF binary format parser
├── txt_parser.go         # TXT position parser
├── txt_parser_helpers.go # TXT parsing utilities
├── web.go                # JSON export methods
├── internal/smile/       # SMILE decoder (BGF compression)
└── examples/             # Usage examples
    ├── parse_bgf/        # BGF file parsing
    ├── parse_txt/        # TXT file parsing
    ├── batch_parse/      # Batch processing
    └── web_server/       # HTTP upload server
```

## Core Components

### Parsers
- **BGF Parser**: Binary match file format with SMILE compression
- **TXT Parser**: Text position files with multilingual support (EN/FR/DE/JP)

### Data Structures
- **Match**: BGF file data (format, version, compressed data)
- **Position**: TXT position data (board, score, evaluations, cube decisions)
- **Evaluation**: Move analysis (rank, equity, probabilities)
- **CubeDecision**: Cube action analysis (action, MWC, equity)

### Web Integration
- **Reader-based parsing**: No temp files required
- **JSON serialization**: All structures support `ToJSON()`
- **HTTP handlers**: Ready-to-use upload examples

## Design Principles

1. **Flexible Input**: Support files, HTTP uploads, memory buffers
2. **Zero Dependencies**: Standard library only (except internal SMILE decoder)
3. **JSON-First**: All structures JSON-serializable for web APIs
4. **Multilingual**: Language-agnostic parsing for international support

## Extension Points

### Adding New Parsers
1. Create `format_parser.go`
2. Implement `ParseFormat(filename string) (*Type, error)`
3. Add `ParseFormatFromReader(io.Reader)` for web support
4. Define types in `types.go`

### Adding New Fields
1. Update types in `types.go`
2. Add JSON tags for serialization
3. Implement parsing in appropriate parser
4. Update examples

## Performance Characteristics

- **BGF Parsing**: O(n) where n = compressed data size
- **TXT Parsing**: O(n) where n = file lines
- **Memory**: Loads entire file into memory
- **Concurrency**: Thread-safe for read operations

## Dependencies

- Go standard library
- `internal/smile`: SMILE compression decoder (bundled, Apache 2.0 license)
