# BGF File Format Specification

## Overview

BGF (BGBlitz Format) stores complete backgammon match data in JSON format with gzip compression and SMILE binary encoding for efficiency.

**Format:** JSON → SMILE encoding → GZIP compression  
**File Extension:** `.bgf`

---

## File Structure

### 1. Header (Single-line JSON)

```json
{
  "version": "0.00011825",
  "compression": "gzip",
  "encoding": "smile"
}
```

### 2. Body (Match Data)

Complete match information in JSON format after decompression and SMILE decoding.

---

## Data Structure

### Match Object

```json
{
  "matchlen": 3,
  "nameGreen": "TachiAI_V",
  "nameRed": "player",
  "date": "Nov 2, 2025",
  "finalGreen": 4,
  "finalRed": 0,
  "useCube": true,
  "useCrawford": true,
  "games": [...]
}
```

**Key Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `matchlen` | int | Match length (points to win) |
| `nameGreen`, `nameRed` | string | Player names |
| `date` | string | Match date |
| `finalGreen`, `finalRed` | int | Final scores |
| `useCube`, `useCrawford`, `useJacoby` | bool | Match rules |
| `prGreen`, `prRed` | object | Performance ratings |
| `luck` | object | Overall luck statistics |
| `games` | array | Array of game objects |

### Game Object

```json
{
  "scoreGreen": 0,
  "scoreRed": 0,
  "wonPoints": 4,
  "wasForfeit": false,
  "wasResignation": false,
  "moves": [...]
}
```

### Move Object

Each move contains complete analysis:

```json
{
  "type": "amove",
  "player": 1,
  "red": 6,
  "green": 4,
  "from": [18, 12, -1, -1],
  "to": [12, 2, -1, -1],
  "equity": {...},
  "luck": {...},
  "pr": {...},
  "moveAnalysis": [...]
}
```

**Move Fields:**

- `player`: 1 (green) or -1 (red)
- `red`, `green`: Dice values
- `from`, `to`: Checker positions (1-24, 0=off, 25=bar, -1=unused)
- `ply`: Analysis depth

### Equity Object

Position evaluation:

```json
{
  "cube": 1,
  "equity": -1.047,
  "matchEquity": 0.247,
  "myWins": 0.181,
  "oppWins": 0.819,
  "myGammon": 0.018,
  "oppGammon": 0.388,
  "cubeDecision": {...}
}
```

**Fields:**
- `equity`: Money game equity
- `matchEquity`: Match winning probability (0-1)
- `myWins`, `oppWins`: Win probabilities
- `myGammon`, `oppGammon`: Gammon probabilities
- `myBackGammon`, `oppBackGammon`: Backgammon probabilities

### Cube Decision

```json
{
  "eqNoDouble": 0.310,
  "eqDoubleTake": 0.195,
  "eqDoublePass": 0.755,
  "stateOnMove": "NO_DOUBLE",
  "stateOther": "ACCEPT"
}
```

**Cube States:** `NO_DOUBLE`, `DOUBLE`, `ACCEPT`, `PASS`

### Move Analysis

Array of alternatives with evaluations:

```json
[
  {
    "move": {"from": [18,12,-1,-1], "to": [12,2,-1,-1]},
    "eq": {"equity": -1.047, "matchEquity": 0.247},
    "played": true,
    "ply": 2
  },
  {
    "move": {"from": [18,14,-1,-1], "to": [12,8,-1,-1]},
    "eq": {"equity": -1.096, "matchEquity": 0.236},
    "played": false,
    "ply": 2
  }
]
```

First entry (played=true) is the actual move played.

### Performance Rating (PR)

```json
{
  "checkerCnt": 26,
  "cubeCnt": 1,
  "checkerError": -0.023,
  "cubeError": 0.000,
  "checkerErrMWP": -0.006,
  "cubeErrMWP": 0.000,
  "mapCntMove": {
    "OK": 25,
    "QUESTIONABLE": 1,
    "ERROR": 0,
    "BLUNDER": 0
  }
}
```

**Quality Ratings:** `OK`, `QUESTIONABLE`, `ERROR`, `BLUNDER`, `LARGE_BLUNDER`

### Luck

Per-move and overall luck calculations:

```json
{
  "luckPlain": -0.005,
  "luckWeighted": -0.021,
  "mode": "Match"
}
```

---

## Data Completeness

BGF files contain **complete match information**:

✅ **Metadata:** Players, date, event, ratings  
✅ **Match Settings:** Rules, match length, cube settings  
✅ **Game Results:** Scores, points won, resignation/forfeit  
✅ **Move-by-Move Data:**
  - Dice rolls
  - Checker movements
  - Position evaluations
  - Win/gammon/backgammon probabilities
  - Equity values (money & match)
  
✅ **Analysis:**
  - Multi-ply evaluation
  - Alternative moves with equity
  - Cube decisions with equity values
  - Error analysis (checker & cube)
  - Luck measurements
  
✅ **Performance:** Error tracking, move quality ratings, overall statistics

---

## Position Encoding

**Board Positions:** 1-24 (standard backgammon numbering)  
**Special Positions:**
- `0` = Borne off
- `25` = Bar
- `-1` = Unused (for moves with < 4 checkers)

---

## Technical Notes

### SMILE Encoding

SMILE (Smile Markup Internet Language Encoding) is a binary JSON format:
- ~50% smaller than text JSON
- Faster parsing
- Maintains full JSON compatibility
- All JSON data types supported

### Decompression

Standard GZIP compression applied after SMILE encoding.

### Parsing Requirements

1. GZIP decompression
2. SMILE decoder
3. JSON parser

The bgfparser handles all steps automatically:

```go
import "github.com/kevung/bgfparser"

match, err := bgfparser.ParseBGF("file.bgf")
// Access data
fmt.Println(match.NameGreen, "vs", match.NameRed)
for _, game := range match.Games {
    for _, move := range game.Moves {
        fmt.Printf("Move: %d-%d, Equity: %.4f\n", 
            move.Red, move.Green, move.Equity.MatchEquity)
    }
}
```

---

## File Size

Sample: `TachiAI_V_player_Nov_2__2025__16_55.bgf`
- **Compressed:** ~15 KB
- **Uncompressed:** ~150 KB (estimated)
- **Content:** 1 game, 27+ moves, 2-ply analysis

---

## See Also

- [API Reference](API_REFERENCE.md) - Parser API
- [Quick Reference](QUICK_REFERENCE.md) - Usage examples
- [Package Documentation](PACKAGE_DOCUMENTATION.md) - Design details
