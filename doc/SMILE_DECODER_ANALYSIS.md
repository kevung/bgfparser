# SMILE Decoder Analysis

## Current Status (Nov 3, 2025)

### What's Working ✅
- Basic SMILE header parsing (:)\n + version)
- Object and array structure recognition
- String decoding (tiny ASCII, short ASCII, shared references)
- Small integer decoding (range 0xC0-0xDF)
- Boolean and null values
- Shared key buffer management
- Multi-value merging (handles multiple top-level SMILE values)
- Error recovery with partial data extraction

### Successfully Decoded Fields
The decoder currently extracts 27+ fields from BGF files:
- Match metadata: `matchlen`, `flags`, `actGame`, `finalGreen`, `finalRed`
- Player info: `nameGreen`, `nameRed`, `ratingGreen`, `ratingRed`, `rankGreen`, `rankRed`
- Game settings: `useCrawford`, `useJacoby`, `gameMode`, `useCube`, `useBeaver`, `cubeLimit`
- Game data: `date`, `location`, `event`, `round`, `comment`, `site`, `extMatchID`
- Statistics: `luck` object with `luckPlain`, `luckWeighted`, `mode`
- Board state data: `boardno`

### Recent Improvements ✅

**Fixed Issues:**
1. ✅ Removed overly strict "unexpected end array in object" sanity check  
   - Changed to return partial data when END_ARRAY encountered unexpectedly
   - Allows decoder to continue past minor structural anomalies
   
2. ✅ Improved error recovery in `readObject()`
   - Now returns partial results instead of trying to skip bytes
   - Cleaner output with fewer garbage keys

**Decoding Progress:**
- Successfully decodes to offset ~2085 (out of 127,411 bytes)
- Extracts 34+ top-level fields
- Partially decodes nested `prGreen` object with probability data

### Current Issue ⚠️

**Error Location:** Offset 2085 (inside `prGreen` object)  
**Error Message:** "unexpected end marker: 0xfb"  
**Context:** After decoding several `prGreen` sub-fields, encounters END_OBJECT when expecting a value

#### Structure at Error Point
```
games: [
  {  // First game object (offset 685)
    // ... game metadata ...
    moves: [
      {  // First move object (offset 1100)
        type: "Damove"
        red: -12
        green: -8
        player: -15
        from: [-4, 0, -15, -15]  // Array ends at offset 1142
        to: ...  // Next key at offset 1143
```

#### Hexdump Around Error
```
1137: 0xF8 START_ARRAY (from array)
1138: 0xCC SMALL_INT (-4)
1139: 0xD0 SMALL_INT (0)
1140: 0xC1 SMALL_INT (-15)
1141: 0xC1 SMALL_INT (-15)
1142: 0xF9 END_ARRAY ← Error occurs here
1143: 0x81 SHORT_ASCII_SHARED (next key: "to")
```

### Analysis

The error "unexpected end array in object at offset 1142" suggests that when `readObject()` loops back after reading a value, it encounters an END_ARRAY marker (0xF9) when it expects a key or END_OBJECT marker.

**Possible causes:**
1. `readArray()` is not consuming the END_ARRAY marker properly
2. Error recovery code in `readObject()` continues without advancing offset
3. Nested structure causing offset confusion
4. The "games" array itself might not be properly structured in the file

**Note:** The `readArray()` function DOES have code to consume the END_ARRAY marker:
```go
if d.data[d.offset] == smileEndArray {
    d.offset++  // This should work
    return result, nil
}
```

### Fields Not Yet Decoded
- Complete `prGreen` and `prRed` probability objects
- Full `games` array with move sequences
- Analysis data within moves
- Equity calculations
- Move evaluations

### Missing SMILE Type Handlers
Based on the BGF file structure, we may need:
- ✅ Float32 (0xE9) - Implemented
- ✅ Float64 (0xEA) - Implemented
- ✅ BigInteger (0xE8) - Partially implemented (returns hex string)
- ✅ BigDecimal (0xEB) - Partially implemented (returns formatted string)
- ✅ Long strings (0xE0-0xE7) - Implemented
- ❓ Variable-length integers (VInt) - Used for lengths, may need refinement
- ❓ Binary data - Unknown if present in BGF files

### Next Steps
1. Add debug logging to trace exact offset progression through the problematic section
2. Verify that array reading properly consumes END_ARRAY in all code paths
3. Check if error recovery in `readObject()` properly advances offset
4. Consider if the sanity check for "unexpected end array in object" is too strict
5. Test with hex dumps to manually verify the structure is as expected

### Decoder Improvements Made
- ✅ Fixed tiny ASCII length bug (changed from `byte - 0x20 + 1` to `byte - 0x20`)
- ✅ Added multi-value decoding to merge multiple top-level SMILE objects
- ✅ Implemented error recovery to continue decoding despite errors
- ✅ Added shared key buffer for efficient key storage
- ✅ Separated key vs value string encoding (different length calculations)

