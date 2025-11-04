# Development Guide# Development Guide



## Setup## Setting Up Development Environment



```bash### Prerequisites

git clone https://github.com/kevung/bgfparser.git

cd bgfparser- Go 1.21 or later

go build ./...- Git

go test ./...- Text editor or IDE (VS Code, GoLand, etc.)

```

### Clone and Build

## Code Style

```bash

- Follow standard Go conventions# Clone the repository

- Use `gofmt` and `go vet`git clone https://github.com/kevung/bgfparser.git

- Document all exported functionscd bgfparser



```go# Build the package

// ParseTXT parses a BGBlitz position text file.go build ./...

func ParseTXT(filename string) (*Position, error) {

    // ...# Run tests (when available)

}go test ./...

```

# Build examples

## Adding Featuresgo build -o bin/parse_txt ./examples/parse_txt/

go build -o bin/parse_bgf ./examples/parse_bgf/

### New Parser Functiongo build -o bin/batch_parse ./examples/batch_parse/

```

1. Define types in `types.go`

2. Implement in appropriate file (e.g., `txt_parser.go`)## Code Style

3. Add example in `examples/`

4. Write tests### Go Conventions



### New File FormatFollow standard Go conventions:

- Use `gofmt` for formatting

1. Create `format_parser.go`- Run `go vet` for static analysis

2. Implement `ParseFormat(filename string) (*Type, error)`- Follow effective Go guidelines

3. Add to README

### Naming Conventions

## Testing

- **Exported functions**: PascalCase (e.g., `ParseTXT`)

```bash- **Unexported functions**: camelCase (e.g., `parseBoard`)

# Run tests- **Constants**: PascalCase or ALL_CAPS for enum-like values

go test ./...- **Package name**: lowercase, single word



# With coverage### Documentation

go test -cover ./...

All exported functions and types must have documentation comments:

# Specific test

go test -run TestParseTXT```go

```// ParseTXT parses a BGBlitz position text file and returns

// a Position struct containing all extracted data.

### Test Structure// Returns an error if the file cannot be read or parsed.

func ParseTXT(filename string) (*Position, error) {

```go    // ...

func TestParseTXT(t *testing.T) {}

    pos, err := ParseTXT("testdata/sample.txt")```

    if err != nil {

        t.Fatalf("ParseTXT failed: %v", err)## Adding New Features

    }

    ### Adding a Parser Function

    if pos.PlayerX != "expected" {

        t.Errorf("got %s, want %s", pos.PlayerX, "expected")1. **Define types** in `types.go`:

    }```go

}type NewFeature struct {

```    Field1 string

    Field2 int

## Building Examples}

```

```bash

go build -o bin/parse_txt ./examples/parse_txt/2. **Implement parser** in appropriate file:

go build -o bin/parse_bgf ./examples/parse_bgf/```go

go build -o bin/batch_parse ./examples/batch_parse/func parseNewFeature(line string) (*NewFeature, error) {

go build -o bin/web_server ./examples/web_server/    // Implementation

```}

```

## Contributing

3. **Integrate** with main parser:

1. Fork the repository```go

2. Create feature branch// In ParseTXT or ParseBGF

3. Make changes with testsif strings.Contains(line, "feature marker") {

4. Run `gofmt`, `go vet`, `go test`    feature, err := parseNewFeature(line)

5. Submit pull request    if err != nil {

        return nil, err

Areas needing contribution:    }

- Full board parsing from ASCII art    pos.NewFeature = feature

- SMILE decoder integration}

- Additional file formats```

- Performance improvements

- Better error messages4. **Add example** in `examples/`:

```go
if pos.NewFeature != nil {
    fmt.Printf("Feature: %v\n", pos.NewFeature)
}
```

### Adding File Format Support

1. Create new parser file (e.g., `xml_parser.go`)
2. Implement `ParseXML(filename string) (*YourType, error)`
3. Add appropriate types to `types.go`
4. Create example in `examples/parse_xml/`
5. Update README.md with new format

## Testing Guidelines

### Writing Tests

Create test files alongside source files:

```go
// txt_parser_test.go
package bgfparser

import "testing"

func TestParseTXT(t *testing.T) {
    pos, err := ParseTXT("testdata/sample.txt")
    if err != nil {
        t.Fatalf("ParseTXT failed: %v", err)
    }
    
    if pos.PlayerX != "expected" {
        t.Errorf("PlayerX = %s, want %s", pos.PlayerX, "expected")
    }
}
```

### Test Data

Store test files in `testdata/` directory:
```
testdata/
├── valid_position.txt
├── french_position.txt
├── cube_decision.txt
└── match.bgf
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestParseTXT
```

## Debugging

### Common Debugging Approaches

1. **Add logging**:
```go
import "log"

log.Printf("Parsing line %d: %s", lineNum, line)
```

2. **Use debugger**: Set breakpoints in VS Code or GoLand

3. **Print intermediate state**:
```go
fmt.Printf("DEBUG: pos=%+v\n", pos)
```

4. **Test with minimal input**: Create simplified test cases

### Debugging Parsers

For regex issues:
```go
import "regexp"

re := regexp.MustCompile(`pattern`)
matches := re.FindStringSubmatch(line)
fmt.Printf("Matches: %v\n", matches)
```

For file reading issues:
```go
data, err := os.ReadFile(filename)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File content:\n%s\n", data)
```

## Performance Profiling

### CPU Profiling

```go
import (
    "runtime/pprof"
    "os"
)

f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()

// Your code here
```

### Memory Profiling

```go
import (
    "runtime/pprof"
    "os"
)

f, _ := os.Create("mem.prof")
pprof.WriteHeapProfile(f)
f.Close()
```

### Benchmark Tests

```go
func BenchmarkParseTXT(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ParseTXT("testdata/sample.txt")
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem
```

## Release Process

### Version Numbering

Follow semantic versioning (semver):
- MAJOR: Breaking API changes
- MINOR: New features, backward compatible
- PATCH: Bug fixes

### Creating a Release

1. Update version in documentation
2. Update CHANGELOG.md
3. Run all tests
4. Tag the release:
```bash
git tag v1.0.0
git push origin v1.0.0
```

### Changelog Format

```markdown
## [1.0.0] - 2025-11-03

### Added
- TXT parser for position files
- BGF parser for match files
- Example programs

### Changed
- Improved error messages

### Fixed
- Parsing of French format files
```

## Contributing

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `go fmt`, `go vet`
6. Submit pull request

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation updated
- [ ] No unnecessary dependencies
- [ ] Error handling is appropriate
- [ ] Examples work correctly

## Useful Tools

### Development Tools

- **gofmt**: Format code
- **go vet**: Static analysis
- **golint**: Linting (deprecated, use golangci-lint)
- **golangci-lint**: Comprehensive linter
- **godoc**: Generate documentation

### IDE Recommendations

- **VS Code** with Go extension
- **GoLand** by JetBrains
- **Vim/Neovim** with vim-go

### Commands

```bash
# Format all code
go fmt ./...

# Check for issues
go vet ./...

# Install linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Generate documentation
godoc -http=:6060
```

## Resources

### Go Resources
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Blog](https://blog.golang.org/)

### Backgammon Resources
- [XGID Format](http://www.gnu.org/software/gnubg/manual/html_node/A-Technical-Description-of-the-XGID.html)
- [GNUbg](https://www.gnu.org/software/gnubg/)
- BGBlitz documentation

## Support

For development questions:
- Open an issue on GitHub
- Check existing issues and pull requests
- Review the documentation in `doc/`
