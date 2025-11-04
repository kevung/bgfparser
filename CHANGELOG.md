# Changelog

All notable changes to the bgfparser project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-11-04

### Added
- **Web-Ready API**: Complete web interface for HTTP file uploads and in-memory parsing
  - `ParseBGFFromReader(io.Reader)` - Parse BGF from any Reader source
  - `ParseTXTFromReader(io.Reader)` - Parse TXT from any Reader source
  - Works with HTTP uploads, memory buffers, network streams, etc.
  - No file system required - parse directly from uploaded data
  
- **JSON Serialization**: Export to JSON for web APIs
  - `(*Match).ToJSON()` - Export match data as JSON
  - `(*Position).ToJSON()` - Export position data as JSON
  - All structures now have JSON tags
  
- **Web Server Example**: Production-ready web server (`examples/web_server/`)
  - HTML upload interface
  - BGF and TXT file support
  - Summary and full JSON endpoints
  - Health check endpoint
  - Complete HTTP handler examples
  
- **Helper Functions**: Extracted TXT parsing utilities (`txt_parser_helpers.go`)
  - Modular, testable parsing functions
  - Reusable components
  
- **Documentation**:
  - [WEB_INTERFACE.md](doc/WEB_INTERFACE.md) - Complete web integration guide
  - [ARCHITECTURE.md](doc/ARCHITECTURE.md) - System architecture diagrams
  - [QUICK_REFERENCE.md](doc/QUICK_REFERENCE.md) - Quick API reference
  - [WEB_IMPLEMENTATION_SUMMARY.md](doc/WEB_IMPLEMENTATION_SUMMARY.md) - Implementation details
  
- **Tests**: Web API test coverage
  - `TestParseBGFFromReader` - Reader-based parsing
  - `TestParseTXTFromReader` - Reader-based parsing
  - `TestMatchToJSON` - JSON export
  - `TestPositionToJSON` - JSON export

### Changed
- **All Data Structures**: Added JSON tags to all types
  - `Position` - Now JSON-serializable with proper field names
  - `Match` - JSON-serializable with nested data support
  - `Evaluation` - JSON tags for API responses
  - `CubeDecision` - JSON tags for API responses
  
- **File-Based Parsers**: Refactored to use Reader-based functions
  - `ParseBGF()` - Now uses `ParseBGFFromReader()` internally
  - `ParseTXT()` - Now uses `ParseTXTFromReader()` internally
  - Improved error handling with filename context
  
- **README**: Enhanced with web integration examples
  - HTTP upload examples
  - JSON export examples
  - Web server quick start
  - Use cases for web applications

### Design Inspiration
- Architecture inspired by [xgparser](https://github.com/kevung/xgparser)
- Flexible input/output patterns for web backends
- Database-ready structures

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
