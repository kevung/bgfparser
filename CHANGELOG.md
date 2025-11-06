# Changelog

## [1.2.0] - 2025-11-05

### Added
- **Multilingual Support**: English, French, German, Japanese
- **Probability Extraction**: Win/WinG/WinBG/LoseG/LoseBG from evaluations
- **Dual Format Support**: Legacy `1)` and new `1.` evaluation formats
- **Test Suite**: 28 multilingual test files, 80.6% coverage

### Fixed
- Player info parsing from separate board lines
- Evaluation vs probability line detection

## [1.1.0] - 2025-11-04

### Added
- **Web API**: `ParseBGFFromReader()`, `ParseTXTFromReader()` for HTTP uploads
- **JSON Export**: `ToJSON()` methods for all structures
- **Web Server Example**: Production-ready upload interface
- **Documentation**: WEB_INTERFACE.md, ARCHITECTURE.md, QUICK_REFERENCE.md

### Changed
- All types now JSON-serializable
- File parsers use Reader-based functions internally

## [1.0.0] - 2025-11-04

### Added
- **Complete SMILE Decoder**: 100% SMILE format decoding support
  - Integrated MIT-licensed SMILE decoder as internal package (from LeLuxNet/X)
  - Zero external dependencies
  - Full support for all SMILE data types (objects, arrays, strings, numbers, booleans)
  - Handles complex nested structures in BGF files
  - Extracts complete game data: moves, equity calculations, analysis, cube decisions
  
### Changed
- **BGF Parser**: Now achieves 100% decoding of SMILE-encoded BGF files
  - Replaced partial custom decoder with production-ready implementation
  - All match data now accessible (previously only 2.56% decoded)
  - Proper handling of shared keys and values for memory efficiency
  
### Removed
- Custom SMILE decoder implementation (replaced with complete solution)
- External dependency on `lelux.net/x/encoding/smile`
- Temporary debug files and test artifacts

## [0.1.0] - 2025-11-03

## [0.1.0] - 2025-11-03

### Added
- Initial release of bgfparser package
- **TXT Parser**: Complete BGBlitz position file parser
  - Position data extraction
  - Player names and scores
  - Match information
  - Position identifiers (Position-ID, Match-ID, XGID)
  - Dice and cube state
  - Move evaluations and statistics
  - Cube decision analysis
  - Multi-language support (English, French)
  - Pip count extraction
  
- **BGF Parser**: Binary match file parser
  - Header parsing with format detection
  - Gzip decompression support
  - JSON data extraction
  - SMILE encoding detection (partial decoding)
  
- **Core Types**: Complete data structures
  - Position type with full game state
  - Evaluation type for move analysis
  - CubeDecision type for cube actions
  - Match type for BGF files
  - ParseError type for error handling
  
- **Examples**: Three demonstration programs
  - `parse_txt`: TXT position file parser
  - `parse_bgf`: BGF match file parser
  - `batch_parse`: Batch processing utility
  
- **Documentation**: Comprehensive docs
  - README with quick start guide
  - API reference documentation
  - Package documentation
  - Development guide
  
- **Testing**: Test suite with 9 test cases
- **License**: MIT License

### Known Limitations (resolved in v1.0.0)
- SMILE decoding incomplete (only 2.56% coverage)
- Language support primarily tested with English and French

---

## Project Status

### Current (v1.0.0)
- ✅ **Zero dependencies** - Fully self-contained
- ✅ **100% SMILE decoding** - Complete BGF file support
- ✅ **Production ready** - All core features implemented
- ✅ **Well documented** - Comprehensive API and usage docs
- ✅ **Properly licensed** - MIT with proper attributions
