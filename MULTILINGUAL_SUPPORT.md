# Multilingual Support for BGF Text Position Parsing

## Overview

The BGF parser now fully supports parsing position files in multiple languages. This enhancement allows the parser to correctly interpret BGBlitz position files exported in English, French, German, and Japanese.

## Supported Languages

| Language | Evaluation Header | Player O Examples | Player X Examples |
|----------|------------------|-------------------|-------------------|
| English  | "Evaluation"     | Green, Blue, etc. | Red, White, etc.  |
| French   | "Évaluation"     | Vert, Bleu, etc.  | Rouge, Blanc, etc.|
| German   | "Bewertung"      | Grün, Blau, etc.  | Rot, Weiß, etc.   |
| Japanese | "評価"           | 緑, 青, etc.       | 赤, 白, etc.      |

## Changes Made

### 1. Player Information Parsing (`txt_parser_helpers.go`)

**Before:** Required both "O:" and "X:" on the same line  
**After:** Parses player information from separate lines

The parser now correctly handles the board diagram format where player information appears on different lines:
```
 +13-14-15-16-17-18------19-20-21-22-23-24-+   O: Green  52
 ...
 +12-11-10--9--8--7-------6--5--4--3--2--1-+   X: Red  111
```

### 2. Evaluation Section Detection (`txt_parser_helpers.go`)

**Enhanced:** Added detection for multilingual evaluation headers

```go
if strings.Contains(line, "Evaluation") ||    // English
   strings.Contains(line, "Évaluation") ||    // French
   strings.Contains(line, "Bewertung") ||     // German
   strings.Contains(line, "評価") {            // Japanese
```

### 3. Evaluation Line Parsing (`txt_parser.go`)

**Enhanced:** Support for two different evaluation formats

#### Format 1 (Legacy)
```
1) 13-11 24-23                0.473 / -0.289
   0.443  0.113  0.002  -  0.557  0.179  0.003
```

#### Format 2 (New with MWP)
```
  1.   0.124 mwp /  -0.492            19/18, 14/12
       0.254  0.000  0.000  -  0.746  0.338  0.004
```

**Key Improvements:**
- Detects rank markers: both `1)` and `1.` formats
- Filters out probability detail lines (second line of each evaluation)
- Correctly identifies evaluation lines vs. probability lines
- Preserves original line for proper pattern matching

## Testing

### New Test Files

Added comprehensive test coverage in `test/2025-11-04/`:
- `01_checkerPosition_*.txt` - Position with evaluations
- `02_NDT_*.txt` - No Double/Take decisions
- `03_DT_*.txt` - Double/Take decisions
- `04_DP_*.txt` - Double/Pass decisions
- `05_NRT_*.txt` - No Redouble/Take decisions
- `06_RT_*.txt` - Redouble/Take decisions
- `07_RP_*.txt` - Redouble/Pass decisions

Each test file category exists in 4 languages: `_EN`, `_FR`, `_DE`, `_JP`

### Test Suite

Added `txt_parser_multilang_test.go` with:
- `TestParseTXT_Multilingual`: Validates parsing across all languages
  - Player names (language-specific)
  - Pip counts
  - Scores and match length
  - Position IDs and XGIDs
  - Evaluation counts and content

- `TestParseTXT_EvaluationSection`: Validates evaluation parsing
  - Move text extraction
  - Equity values
  - Ranking consistency
  - Language-independent evaluation headers

### Test Results

```
✓ All 28 test file combinations (7 types × 4 languages)
✓ All existing tests continue to pass
✓ Code coverage: 80.6%
```

## Usage Examples

### English File
```go
pos, err := bgfparser.ParseTXT("test/2025-11-04/01_checkerPosition_EN.txt")
// PlayerO: "Green", PlayerX: "Red"
```

### French File
```go
pos, err := bgfparser.ParseTXT("test/2025-11-04/01_checkerPosition_FR.txt")
// PlayerO: "Vert", PlayerX: "Rouge"
```

### German File
```go
pos, err := bgfparser.ParseTXT("test/2025-11-04/01_checkerPosition_DE.txt")
// PlayerO: "Grün", PlayerX: "Rot"
```

### Japanese File
```go
pos, err := bgfparser.ParseTXT("test/2025-11-04/01_checkerPosition_JP.txt")
// PlayerO: "緑", PlayerX: "赤"
```

## Backward Compatibility

All changes are **fully backward compatible**:
- Existing English files continue to work
- Legacy evaluation format still supported
- All existing tests pass
- No API changes

## Implementation Details

### Key Files Modified

1. **`txt_parser_helpers.go`**
   - `parsePlayerInfo()`: Changed OR logic to support separate lines
   - `handleEvaluationSection()`: Added multilingual header detection

2. **`txt_parser.go`**
   - `parseEvaluation()`: 
     - Added support for both `1)` and `1.` rank formats
     - Improved probability line filtering
     - Better detection of evaluation vs. detail lines

3. **`CHANGELOG.md`**
   - Documented all changes in version 1.2.0

### Regex Patterns Used

- Rank detection: `^\s*\d+[.)]` - Matches both "1)" and "1." formats
- Probability filtering: `^\d+\.\d+\s` - Identifies decimal-starting lines
- Ensures only proper evaluation lines are parsed

## Future Enhancements

Potential areas for expansion:
- Support for additional languages (Spanish, Italian, etc.)
- Custom player name patterns
- More evaluation format variations
- Localized error messages
