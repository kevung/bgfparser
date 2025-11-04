# Package Documentation# Package Documentation



## Package Structure## Package Structure



### Core FilesThe `bgfparser` package is organized into several files, each with a specific purpose:



- **types.go** - Data structures (Position, Evaluation, CubeDecision, Match)### Core Files

- **txt_parser.go** - TXT position file parser  

- **txt_parser_helpers.go** - Helper functions for TXT parsing#### types.go

- **bgf_parser.go** - BGF binary match file parserDefines all data structures used throughout the package:

- **web.go** - io.Reader-based parsers for web integration- `Position`: Complete backgammon position with all metadata

- `Evaluation`: Move evaluation with equity and probabilities

## Design- `CubeDecision`: Cube action analysis

- `Match`: BGF match file container

### TXT Parser- `ParseError`: Custom error type for parsing failures

Line-by-line state machine using regex patterns to extract:

- Board state from ASCII art#### txt_parser.go

- Player metadata and scoresImplements parsing for BGBlitz TXT position files:

- Move evaluations with equity- `ParseTXT(filename string) (*Position, error)`: Main entry point

- Cube decisions- `parseBoard(pos *Position, lines []string)`: Board state extraction

- Multi-language support (EN, FR)- `parseXGID(pos *Position, xgid string)`: XGID parser

- `parseEvaluation(line string, rank *int) *Evaluation`: Move evaluation parser

### BGF Parser- `parseCubeDecision(line string) *CubeDecision`: Cube decision parser

Two-phase JSON reader:

1. Parse header (format, version, compression)#### bgf_parser.go

2. Decompress gzip and extract match dataImplements parsing for BGBlitz BGF binary match files:

3. Detect SMILE encoding (partial support)- `ParseBGF(filename string) (*Match, error)`: Main entry point

- `(*Match) GetMatchInfo() map[string]interface{}`: Extract metadata

### Error Handling- `(*Match) String() string`: Human-readable representation

Custom `ParseError` type includes filename, line number, and descriptive message for debugging.

## Design Decisions

## Common Usage

### Parsing Approach

### Basic Parsing

```go**TXT Files**: Line-by-line state machine parser

// Parse files- Scans through file sequentially

pos, err := bgfparser.ParseTXT("position.txt")- Maintains state (in evaluation section, in cube decision, etc.)

match, err := bgfparser.ParseBGF("match.bgf")- Uses regular expressions for pattern matching

- Extracts structured data from ASCII art and formatted text

// Parse from io.Reader (web uploads)

pos, err := bgfparser.ParseTXTFromReader(file)**BGF Files**: Two-phase reader

match, err := bgfparser.ParseBGFFromReader(file)1. Parse uncompressed JSON header

2. Decompress and parse main data

// Export to JSON3. Detect SMILE encoding (full decoding requires external library)

jsonData, err := pos.ToJSON()

```### Error Handling



### Access EvaluationsThe package uses a custom `ParseError` type that includes:

```go- Filename

// Best move (first in array)- Line number (when applicable)

if len(pos.Evaluations) > 0 {- Descriptive error message

    best := pos.Evaluations[0]

    fmt.Printf("%s: %.3f\n", best.Move, best.Equity)This provides clear context for debugging parsing issues.

}

### Data Representation

// Find marked best

for _, eval := range pos.Evaluations {**Board State**: Array of 26 integers

    if eval.IsBest {- Index 0: Unused

        fmt.Printf("Best: %s\n", eval.Move)- Index 1-24: Points (positive for X, negative for O)

    }- Index 25: Bar

}

```**Maps for Dynamic Data**: Used for:

- `OnBar`: Checkers on bar per player

### Cube Decisions- `PipCount`: Pip count per player

```go

if pos.CubeDecision != nil {This allows flexible representation without fixed player ordering.

    fmt.Printf("Action: %s (MWC: %.1f%%)\n", 

        pos.CubeDecision.Action, ## API Usage Patterns

        pos.CubeDecision.MWC * 100)

}### Basic Parsing

```

```go

## Extension// Parse a file

position, err := bgfparser.ParseTXT("file.txt")

### Custom Analysisif err != nil {

```go    // Handle error

import "github.com/kevung/bgfparser"    if parseErr, ok := err.(*bgfparser.ParseError); ok {

        log.Printf("Parse error at line %d: %s", parseErr.Line, parseErr.Message)

func FindBlunders(pos *bgfparser.Position, threshold float64) []bgfparser.Evaluation {    }

    var blunders []bgfparser.Evaluation    return

    for _, eval := range pos.Evaluations[1:] {}

        if -eval.Diff >= threshold {

            blunders = append(blunders, eval)// Use the position

        }fmt.Printf("Player on roll: %s\n", position.OnRoll)

    }```

    return blunders

}### Accessing Evaluations

```

```go

### Adding Parsers// Get best move

1. Define types in `types.go`if len(position.Evaluations) > 0 {

2. Create parser file (e.g., `newformat_parser.go`)    best := position.Evaluations[0]

3. Implement `ParseNewFormat(filename string) (*YourType, error)`    fmt.Printf("Best: %s (%.3f)\n", best.Move, best.Equity)

4. Follow existing error patterns    

5. Add tests    // Compare with second best

    if len(position.Evaluations) > 1 {

## Performance        diff := position.Evaluations[0].Equity - position.Evaluations[1].Equity

        fmt.Printf("Advantage: %.3f\n", diff)

- **TXT parsing**: ~1-2ms per file    }

- **BGF parsing**: ~10-50ms (depends on size/compression)}

- **Memory**: ~500 bytes base + ~150 bytes per evaluation

- **Tip**: Use goroutines for batch processing// Find marked best move

for _, eval := range position.Evaluations {

## Testing    if eval.IsBest {

        fmt.Printf("Marked best: %s\n", eval.Move)

Test coverage includes:    }

- Valid and invalid files}

- Multi-language support```

- Edge cases (empty, malformed)

- Various positions (cube, bar, Crawford)### Working with Cube Decisions



See `parser_test.go` for examples.```go

if position.CubeDecision != nil {

## See Also    cd := position.CubeDecision

    fmt.Printf("Recommended: %s\n", cd.Action)

- [API Reference](API_REFERENCE.md) - Complete API documentation    

- [Quick Reference](QUICK_REFERENCE.md) - Common patterns    switch cd.Action {

- [BGF Format](BGF_FORMAT_SPECIFICATION.md) - BGF file format details    case "Double/Take":

- [Web Interface](WEB_INTERFACE.md) - HTTP upload integration        fmt.Println("Should double and opponent should take")

    case "Double/Pass":
        fmt.Println("Should double and opponent should pass")
    case "No Double":
        fmt.Println("Should not double")
    }
    
    fmt.Printf("MWC: %.1f%%\n", cd.MWC * 100)
}
```

### Batch Processing

```go
files := []string{"game1.txt", "game2.txt", "game3.txt"}

for _, file := range files {
    pos, err := bgfparser.ParseTXT(file)
    if err != nil {
        log.Printf("Error parsing %s: %v\n", file, err)
        continue
    }
    
    // Process position
    processPosition(pos)
}
```

## Extension Points

### Adding New Parsers

To add support for a new format:

1. Define new types in `types.go` if needed
2. Create a new file (e.g., `newformat_parser.go`)
3. Implement `ParseNewFormat(filename string) (*YourType, error)`
4. Follow the existing error handling pattern
5. Add tests and examples

### Custom Analysis

The package can be extended with analysis functions:

```go
package myanalysis

import "github.com/kevung/bgfparser"

// AnalyzeBlunders finds significant equity errors
func AnalyzeBlunders(pos *bgfparser.Position, threshold float64) []bgfparser.Evaluation {
    var blunders []bgfparser.Evaluation
    
    if len(pos.Evaluations) < 2 {
        return blunders
    }
    
    for _, eval := range pos.Evaluations[1:] {
        if -eval.Diff >= threshold {
            blunders = append(blunders, eval)
        }
    }
    
    return blunders
}
```

## Performance Considerations

### Memory Usage

- Position objects are relatively lightweight (~500 bytes base + evaluations)
- Each evaluation adds ~100-200 bytes
- BGF match data can be large (depends on match length)

### Parsing Speed

- TXT files: ~1-2ms per file (typical position with evaluations)
- BGF files: ~10-50ms per file (depending on compression and size)
- Batch processing benefits from parallel processing

### Optimization Tips

1. **Reuse Position objects** when parsing many files
2. **Use batch parsing** with goroutines for large datasets
3. **Stream processing** for very large match collections
4. **Index frequently accessed fields** in your own data structures

## Testing Strategy

### Unit Tests

Each parser should have tests for:
- Valid input files
- Malformed files
- Edge cases (empty files, missing sections)
- Different languages
- Various position types (cube offered, on bar, etc.)

### Integration Tests

Test with:
- Real BGBlitz output files
- Various match lengths and scores
- Different cube values
- Crawford situations

### Example Test Structure

```go
func TestParseTXT(t *testing.T) {
    tests := []struct {
        name    string
        file    string
        wantErr bool
    }{
        {"Valid position", "testdata/valid.txt", false},
        {"Missing file", "testdata/missing.txt", true},
        {"French format", "testdata/french.txt", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pos, err := ParseTXT(tt.file)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseTXT() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && pos == nil {
                t.Error("ParseTXT() returned nil position")
            }
        })
    }
}
```

## Common Issues and Solutions

### Issue: Player names not parsed

**Solution**: The parser looks for "O:" and "X:" markers. Ensure the file follows standard BGBlitz format.

### Issue: Evaluations missing

**Solution**: Check for "Evaluation" or "Évaluation" header. The parser supports both English and French.

### Issue: Cube value incorrect

**Solution**: Cube values are extracted from XGID and board display. Verify the file has proper formatting.

### Issue: SMILE decoding fails

**Solution**: SMILE is not yet fully supported. You can:
1. Export matches without SMILE encoding
2. Integrate a SMILE decoder library
3. Use the header information available

## Future Development

### Planned Features

1. **Full board parsing**: Complete extraction of checker positions from ASCII art
2. **SMILE support**: Integration with SMILE decoder
3. **Statistics aggregation**: Functions to aggregate data across multiple positions
4. **Match analysis**: Tools for analyzing entire matches
5. **Export functions**: Convert positions to other formats (XGID, GNUbg, etc.)

### Community Contributions Welcome

Areas where contributions would be valuable:
- Additional file format support
- Performance improvements
- More comprehensive parsing
- Additional analysis functions
- Better error messages
- Localization support
