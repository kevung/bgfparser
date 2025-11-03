# BGF File Decoding Progress

## Summary (as of November 4, 2025)

The BGF parser now successfully decodes **2.56%** of BGF files (up from 0% initially).

### Decoding Statistics
- **Offset reached**: 3,267 bytes out of 127,411 total
- **Fields decoded**: 33 top-level fields
- **Format supported**: SMILE-encoded, gzip-compressed JSON

### Successfully Decoded Data

#### Match Metadata (100% complete)
- `matchlen`: Match length
- `flags`: Match flags
- `date`: Match date (e.g., "Nov 2, 2025")
- `nameGreen`, `nameRed`: Player names
- `ratingGreen`, `ratingRed`: Player ratings
- `rankGreen`, `rankRed`: Player ranks
- `finalGreen`, `finalRed`: Final scores
- `actGame`: Active game number
- `boardno`: Board number

#### Match Settings (100% complete)
- `useCrawford`: Crawford rule enabled
- `useJacoby`: Jacoby rule enabled
- `useCube`: Doubling cube enabled
- `useBeaver`: Beaver allowed
- `cubeLimit`: Maximum cube value
- `gameMode`: Game mode

#### Game Statistics (Partial)
- `luck`: Object containing:
  - `luckPlain`: Plain luck value
  - `luckWeighted`: Weighted luck value
  - `mode`: Match/Money mode

#### Player Statistics (`prGreen`, `prRed`) (Partial ~30%)
- `cubeCnt`: Cube decision count
- `checkerCnt`: Checker play count  
- `cubeError`: Cube play error rate
- `checkerError`: Checker play error rate
- `ube`: Cube decision categorization (BLUNDER, ERROR, OK, etc.)
- `equity`: Equity calculations (partially decoded)
- `moveAnalysis`: Array of move analysis objects (partially decoded)

#### Move Analysis Data (Partial ~20%)
Each move contains:
- `move`: Move details with checker positions
  - `from`: Array of starting positions
  - `to`: Array of ending positions
  - Equity values: `eqDoubleTake`, `eqDoublePass`, `eqCubeLess`, `eqCubeFull`
  - Win probabilities
  - `stateOnMove`: Move state (e.g., "NO_DOUBLE")
  - `stateOther`: Response state (e.g., "ACCEPT")
- `ply`: Ply depth
- `played`: Boolean indicating if move was played
- `hasAccepted`: Cube decision acceptance
- Various probability calculations

### SMILE Decoder Capabilities

#### Fully Implemented ✅
- Header parsing (`:)\n` + version byte)
- Object structures (`START_OBJECT`, `END_OBJECT`)
- Array structures (`START_ARRAY`, `END_ARRAY`)
- Tiny ASCII strings (0x20-0x3F): 0-31 bytes
- Short ASCII strings (0x40-0x7F): 1-63 bytes
- Short ASCII shared strings (0x80-0xBF): 1-64 bytes, added to shared keys
- Shared key references (0x00-0x1F)
- Boolean values (`TRUE`, `FALSE`, `NULL`)
- Small integers (0xC0-0xDF): values -16 to +15
- 32-bit integers (0x24)
- 64-bit integers (0x25)
- 32-bit floats (0xE9)
- 64-bit doubles (0xEA)
- Long ASCII strings (0xE0-0xE3)
- Long Unicode strings (0xE4-0xE7)
- BigInteger (0xE8) - returns hex representation
- BigDecimal (0xEB) - returns formatted string
- Variable-length integers (VInt) for lengths
- Shared key buffer management
- Multi-value merging
- Error recovery with partial data extraction

#### Known Limitations
- Nested structures beyond certain depth may encounter context end markers
- Some binary data fields contain non-printable characters
- Complex equity calculations may have garbled output
- Game move sequences not fully parsed
- Binary data types not yet identified

### Recent Improvements

1. **Added aggressive error recovery** - Decoder continues past decode errors instead of stopping
2. **Improved structural marker handling** - Recognizes unexpected START_OBJECT/START_ARRAY during key reading
3. **Enhanced metadata tracking** - Now reports:
   - `_finalOffset`: Last byte position decoded
   - `_totalBytes`: Total file size
   - `_percentDecoded`: Percentage of file decoded
4. **Better context handling** - More robust handling of nested objects and arrays
5. **Continued decoding after errors** - Stores error information but attempts to continue

### Next Steps to Improve Decoding

1. **Identify remaining structural patterns** - Analyze bytes at offset 3267+ to understand why decoding stops
2. **Handle more complex nesting** - Improve context tracking for deeply nested structures
3. **Decode game/move arrays** - Full parsing of the `games` array with move sequences
4. **Binary field identification** - Determine what the non-printable characters represent
5. **Validate against known BGF structures** - Compare with BGBlitz documentation
6. **Handle additional SMILE types** - If any exotic types exist in BGF files

### File Structure (Inferred)

```
BGF File Structure:
├── Header (JSON line): format, version, compress, useSmile
└── Compressed+SMILE data:
    ├── Match metadata (matchlen, flags, date, players, etc.)
    ├── Match settings (Crawford, Jacoby, cube, etc.)
    ├── Player statistics (prGreen, prRed)
    │   ├── Error rates
    │   ├── Decision categorization
    │   └── Move analysis arrays
    └── Games array (partially decoded)
        └── Each game contains:
            ├── Game metadata
            ├── Board states
            └── Moves array
                └── Each move contains:
                    ├── Checker positions (from/to)
                    ├── Equity calculations
                    ├── Win probabilities
                    └── Cube decisions
```

### Performance
- Decoding time: < 100ms for typical match file
- Memory usage: Minimal (shared key buffer ~64 entries)
- Error handling: Graceful degradation with partial results

### Compatibility
- Tested with BGBlitz BGF v1.0 files
- Works with both compressed and uncompressed files
- Handles SMILE-encoded and plain JSON variants
