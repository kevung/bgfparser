# Development Guide

## Setup

```bash
git clone https://github.com/kevung/bgfparser.git
cd bgfparser
go build ./...
go test ./...
```

## Code Style

- Follow Go conventions: `gofmt`, `go vet`
- PascalCase for exported, camelCase for unexported
- Document all exported functions

## Adding Features

### New Parser Function

1. Define types in `types.go`
2. Implement in parser file (e.g., `txt_parser.go`)
3. Add example in `examples/`
4. Write tests

### New File Format

1. Create `format_parser.go`
2. Implement `ParseFormat(filename string) (*Type, error)`
3. Add to README

## Testing

```bash
go test ./...          # Run all tests
go test -cover ./...   # With coverage
go test -run TestName  # Specific test
```

## Building Examples

```bash
go build -o bin/parse_txt ./examples/parse_txt/
go build -o bin/parse_bgf ./examples/parse_bgf/
go build -o bin/batch_parse ./examples/batch_parse/
go build -o bin/web_server ./examples/web_server/
```

## Contributing

1. Fork repository
2. Create feature branch
3. Make changes with tests
4. Run `gofmt`, `go vet`, `go test`
5. Submit pull request
