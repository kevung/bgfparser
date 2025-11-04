# BGFParser Examples

Example programs demonstrating bgfparser usage.

## Quick Start

```bash
# Build all tools
go build -o bin/parse_txt ./examples/parse_txt/
go build -o bin/parse_bgf ./examples/parse_bgf/
go build -o bin/batch_parse ./examples/batch_parse/
go build -o bin/web_server ./examples/web_server/

# Or run directly
go run examples/parse_txt/main.go tmp/blunder21_EN.txt
go run examples/batch_parse/main.go tmp/
go run examples/web_server/main.go  # http://localhost:8080
```

## Tools

### parse_txt - Parse TXT Position Files

Parse and display a BGBlitz TXT position file.

```bash
./bin/parse_txt <filename.txt>
./bin/parse_txt tmp/blunder21_EN.txt
```

**Shows:** Players, scores, dice, cube, evaluations, cube decisions, position IDs

### parse_bgf - Parse BGF Match Files

Parse and display a BGBlitz BGF binary match file.

```bash
./bin/parse_bgf <filename.bgf>
./bin/parse_bgf tmp/TachiAI_V_player_Nov_2__2025__16_55.bgf
```

**Shows:** Format, version, compression, SMILE encoding, match metadata

### batch_parse - Process Multiple Files

Parse all BGF and TXT files in a directory.

```bash
./bin/batch_parse <directory>
./bin/batch_parse tmp/
```

**Processes:** All .txt and .bgf files, shows summary for each

### web_server - Web Upload Interface

Run a web server for uploading and analyzing files.

```bash
./bin/web_server
# Visit http://localhost:8080
```

**Features:**
- Upload BGF/TXT files via web interface
- JSON API responses
- Endpoints: `/upload/bgf`, `/upload/txt`, `/health`

## Sample Output

### parse_txt

```
=== Position Information ===
Players: player (X) vs TachiAI_V (O)
Score: 0-0 in 3 point match
On Roll: X with 2-1
Cube: 1
Evaluations: 8 moves

Best move: 13-11 24-23 (Equity: 0.473)
```

### parse_bgf

```
Format: BGF v1.0
Compressed: true
SMILE encoding: true
Players: TachiAI_V vs player
Score: 4-0
Decoded fields: 30
```

## Code Templates

### Extract Data

```go
pos, err := bgfparser.ParseTXT("file.txt")
if err != nil {
    log.Fatal(err)
}

equity := pos.Evaluations[0].Equity
bestMove := pos.Evaluations[0].Move
```

### Find Blunders

```go
if len(pos.Evaluations) >= 2 {
    diff := pos.Evaluations[0].Equity - pos.Evaluations[1].Equity
    if diff > 0.100 {
        fmt.Printf("Blunder: %.3f equity loss\n", diff)
    }
}
```

### Batch Processing

```go
files, _ := filepath.Glob("positions/*.txt")
for _, file := range files {
    pos, err := bgfparser.ParseTXT(file)
    if err != nil {
        continue
    }
    // Process position
}
```

## Sample Data

Test files in `tmp/`:

**TXT Files:**
- `blunder21_EN.txt` - English, 2-1 roll
- `blunder32_FR.txt` - French, 3-2 roll
- `blunderBar_FR.txt` - Checker on bar
- `BlunderCubeOffered_EN.txt` - Cube decision

**BGF Files:**
- `TachiAI_V_player_Nov_2__2025__16_55.bgf` - Complete match
- `TachiAI_V_player_Nov_2__2025__17_1.bgf` - Complete match

## Troubleshooting

**File not found:** Use absolute paths or check current directory  
**SMILE encoding:** Expected for BGF files - header info is extracted  
**Empty evaluations:** File may not contain evaluation section

## See Also

- [Main README](../README.md) - Package overview
- [API Reference](../doc/API_REFERENCE.md) - Complete API
- [Quick Reference](../doc/QUICK_REFERENCE.md) - Common patterns
