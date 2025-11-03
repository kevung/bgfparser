# SMILE Decoder

This package implements decoding of SMILE (Smile Format - binary JSON) as defined in the [SMILE Format Specification](https://github.com/FasterXML/smile-format-specification/blob/master/smile-specification.md).

## Attribution

This is a copy of the SMILE decoder from [LeLuxNet/X](https://gitlab.com/LeLuxNet/X/-/tree/c09411c26dfb/encoding/smile) (commit c09411c26dfb, October 11, 2024), used under the MIT License.

The original source is available at:
- GitLab: https://gitlab.com/LeLuxNet/X/-/tree/c09411c26dfb/encoding/smile
- Go package: https://pkg.go.dev/lelux.net/x/encoding/smile

We've included it as an internal package to:
1. Avoid external dependencies that could disappear
2. Ensure long-term stability of the bgfparser project
3. Maintain full control over the decoder implementation

## License

MIT License - See LICENSE file in this directory.

## Usage

```go
import "github.com/unger/bgfparser/internal/smile"

var result interface{}
err := smile.Unmarshal(smileData, &result)
```

## What is SMILE?

SMILE is a binary JSON format that is more compact and faster to parse than text JSON. It's used by BGBlitz to store backgammon match data efficiently in BGF files.
