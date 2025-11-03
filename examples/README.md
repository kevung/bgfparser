# BGFParser Examples

This directory contains example programs demonstrating the use of the bgfparser package.

## Examples

### 1. parse_txt - Position File Parser

Parses a single BGBlitz TXT position file and displays all extracted information.

**Usage:**
```bash
go run examples/parse_txt/main.go <filename.txt>
```

**Example:**
```bash
go run examples/parse_txt/main.go tmp/blunder21_EN.txt
```

**Output:**
- Position information (players, scores, match length)
- Dice roll and cube state
- Position identifiers (Position-ID, Match-ID, XGID)
- Pip counts
- Move evaluations with equity values
- Cube decisions (if present)

### 2. parse_bgf - Match File Parser

Parses a BGBlitz BGF binary match file and displays header information.

**Usage:**
```bash
go run examples/parse_bgf/main.go <filename.bgf>
```

**Example:**
```bash
go run examples/parse_bgf/main.go tmp/TachiAI_V_player_Nov_2__2025__16_55.bgf
```

**Output:**
- BGF format information
- Version and compression details
- SMILE encoding detection
- Match metadata (when available)

**Note:** Files using SMILE encoding will show header information but note that full decoding requires an external library.

### 3. batch_parse - Batch Processor

Processes all TXT and BGF files in a directory, providing a summary of each.

**Usage:**
```bash
go run examples/batch_parse/main.go <directory>
```

**Example:**
```bash
go run examples/batch_parse/main.go tmp/
```

**Output:**
- Summary for each file found
- Player information and scores
- Number of evaluations
- Best moves
- Error messages for parsing failures

## Building Examples

You can build standalone executables:

```bash
# Build all examples
go build -o bin/parse_txt ./examples/parse_txt/
go build -o bin/parse_bgf ./examples/parse_bgf/
go build -o bin/batch_parse ./examples/batch_parse/

# Run built executables
./bin/parse_txt tmp/blunder21_EN.txt
./bin/parse_bgf tmp/TachiAI_V_player_Nov_2__2025__16_55.bgf
./bin/batch_parse tmp/
```

## Example Output

### parse_txt Output

```
=== Position Information ===
Players: player (X) vs TachiAI_V (O)
Score: TachiAI_V 0 - 0 player
Match Length: 3 points
On Roll: X
Dice: 2-1
Cube: 1
Pip Count: X=167, O=161

=== Identifiers ===
Position-ID: mGfwATDgc/ABMA
Match-ID: cAllAAAAAAAE
XGID: -b----E-C---eE---b-d-b--B-:0:0:1:21:0:0:0:3:10

=== Move Evaluations ===
1) 13-11 24-23
   Equity: 0.473 (-0.289)
2) 24-22 24-23
   Equity: 0.469 (-0.331)
...
```

### parse_bgf Output

```
=== BGF Match File ===
BGF Match (Format: BGF, Version: 1.0, Compressed: true, SMILE: true)

=== Header Information ===
Format: BGF
Version: 1.0
Compressed: true
Uses SMILE encoding: true

Note: SMILE encoding is detected but not yet fully supported.
```

### batch_parse Output

```
=== Parsing TXT: blunder21_EN.txt ===
Players:  vs 
Score: 0-0 in a 3 point match
On Roll: X with 2-1
Evaluations: 8 moves analyzed

=== Parsing BGF: match.bgf ===
Format: BGF v1.0
Compressed: true, SMILE: true
Note: SMILE encoding detected (binary JSON format)
```

## Using Examples as Templates

These examples can be used as templates for your own programs:

### Extract Specific Data

```go
// Based on parse_txt example
pos, err := bgfparser.ParseTXT("file.txt")
if err != nil {
    log.Fatal(err)
}

// Get only what you need
equity := pos.Evaluations[0].Equity
bestMove := pos.Evaluations[0].Move
```

### Process Multiple Files

```go
// Based on batch_parse example
files, _ := filepath.Glob("positions/*.txt")
for _, file := range files {
    pos, err := bgfparser.ParseTXT(file)
    if err != nil {
        continue
    }
    // Process position
}
```

### Analyze Blunders

```go
// Find significant errors
if len(pos.Evaluations) >= 2 {
    diff := pos.Evaluations[0].Equity - pos.Evaluations[1].Equity
    if diff > 0.100 {
        fmt.Printf("Blunder found: %.3f equity loss\n", diff)
    }
}
```

## Sample Data

The `tmp/` directory contains sample files for testing:

**TXT Files:**
- `blunder21_EN.txt` - Simple position with 2-1 roll
- `blunder32_FR.txt` - French format with 3-2 roll
- `blunderBar_FR.txt` - Position with checker on bar
- `BlunderCubeOffered_EN.txt` - Cube decision position
- And more...

**BGF Files:**
- `TachiAI_V_player_Nov_2__2025__16_55.bgf` - Complete match
- `TachiAI_V_player_Nov_2__2025__17_1.bgf` - Complete match

## Troubleshooting

### File Not Found
```
Error: open file.txt: no such file or directory
```
**Solution:** Check the file path. Use absolute paths or ensure you're in the correct directory.

### SMILE Encoding
```
Error: SMILE encoding is not yet supported
```
**Solution:** This is expected for BGF files. The parser extracts header information but cannot decode SMILE-encoded data without an external library.

### Empty Evaluations
```
No evaluations parsed
```
**Solution:** The file may not contain an evaluation section. Check that the file is a complete BGBlitz position export.

## Further Reading

- [API Reference](../doc/API_REFERENCE.md) - Complete API documentation
- [Package Documentation](../doc/PACKAGE_DOCUMENTATION.md) - Design and patterns
- [Development Guide](../doc/DEVELOPMENT.md) - Contributing guidelines
