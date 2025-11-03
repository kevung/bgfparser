# Changelog

All notable changes to the bgfparser project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-03

### Added
- Initial release of bgfparser package
- TXT parser for BGBlitz position files
  - Complete position parsing
  - Player names and scores extraction
  - Match information parsing
  - Position identifiers (Position-ID, Match-ID, XGID)
  - Dice and cube state extraction
  - Move evaluations parsing
  - Cube decision analysis
  - Multi-language support (English, French)
  - Pip count extraction
- BGF parser for BGBlitz binary match files
  - Header parsing with format detection
  - Gzip decompression support
  - JSON data extraction
  - SMILE encoding detection
- Core data structures
  - Position type with complete state
  - Evaluation type for move analysis
  - CubeDecision type for cube actions
  - Match type for BGF files
  - ParseError type for detailed error handling
- Three example programs
  - parse_txt: Detailed TXT position parser
  - parse_bgf: BGF match file parser
  - batch_parse: Batch processing tool
- Comprehensive test suite
  - 9 test cases covering main functionality
  - Multi-language testing
  - Error handling verification
- Complete documentation
  - README.md with quick start guide
  - API_REFERENCE.md with complete API docs
  - PACKAGE_DOCUMENTATION.md with design patterns
  - DEVELOPMENT.md with development guide
  - PACKAGE_SUMMARY.md with project overview
- MIT License
- Go module support (go.mod)

### Known Limitations
- Board state extraction from ASCII art is partial
- SMILE decoding for BGF files not implemented (requires external library)
- Some evaluation statistics may not be fully extracted
- Language support primarily tested with English and French

## [Unreleased]

### Planned Features
- Full SMILE decoder integration for BGF files
- Complete board state extraction from ASCII representation
- Additional statistics parsing from evaluations
- Export functions to other formats (GNUbg, XG)
- Match analysis tools
- Performance optimizations for batch processing
- Additional language support
- More comprehensive test coverage

---

## Version History

### 1.0.0 - Initial Release
First stable release with full TXT parsing and partial BGF parsing support.
