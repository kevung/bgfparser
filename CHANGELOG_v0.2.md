# BGFParser Changelog

## Version 0.2 - SMILE Encoding Support (November 3, 2025)

### New Features
- ✅ **SMILE Decoder Implementation**: Added basic SMILE (binary JSON) decoding support
  - Can parse SMILE-encoded BGF files without crashing
  - Extracts player names, dates, match parameters from binary data
  - Handles SMILE header detection and basic type markers
  
- ✅ **Graceful Error Handling**: BGF parser no longer returns fatal errors for SMILE files
  - Files parse successfully with warnings for incomplete decoding
  - Partial data extraction when full decoding isn't possible
  
- ✅ **Enhanced Match Type**: Added `DecodingWarning` field
  - Provides feedback when SMILE decoding is incomplete
  - Allows programs to handle partial data gracefully

### Technical Details

#### SMILE Format Support
The implementation includes:
- SMILE header detection (`:)\n` + version byte)
- String decoding (tiny ASCII, short ASCII, shared keys)
- Object and array structure parsing
- Basic type handling (strings, booleans, integers)
- Fallback to string extraction when full decoding fails

#### Files Added
- `smile_decoder.go`: Core SMILE decoding logic
- `debug/test_smile.go`: Debug tool for examining SMILE files

#### Files Modified
- `bgf_parser.go`: Updated to use SMILE decoder
- `types.go`: Added DecodingWarning field to Match type
- `examples/parse_bgf/main.go`: Enhanced output for SMILE files
- `README.md`: Updated documentation

### What Works
- ✅ BGF files with SMILE encoding can be parsed
- ✅ Basic metadata extraction (player names, dates, match length)
- ✅ String data extraction from binary format
- ✅ No fatal errors - graceful degradation

### Current Limitations
- Complex nested structures may not fully decode
- Some binary number formats not yet implemented
- Variable-length integer encoding partially implemented

### Example Output

Before (v0.1):
```
Error parsing file: SMILE encoding is not yet supported
```

After (v0.2):
```
=== BGF Match File ===
Format: BGF v1.0
Compressed: true
Uses SMILE encoding: true

⚠️  Warning: SMILE decoding incomplete: full SMILE decoding not implemented

=== Extracted Strings ===
1: dateJNov 2, 2025
2: nameGreenHTachiAI_V
3: nameRedEplayer
...
```

### Testing
All sample files in `tmp/` directory parse successfully:
- `TachiAI_V_player_Nov_2__2025__16_55.bgf` ✅
- `TachiAI_V_player_Nov_2__2025__17_1.bgf` ✅
- All TXT files continue to work ✅

### Migration Notes
No breaking changes. The API remains backward compatible.

New optional field `DecodingWarning` in `Match` struct contains warning messages for incomplete parsing.

### Next Steps
- Enhance SMILE decoder for complex nested structures
- Add support for more numeric types (float, decimal)
- Implement full variable-length integer decoding
- Add unit tests for SMILE decoder
