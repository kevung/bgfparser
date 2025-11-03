package bgfparser

import (
	"encoding/json"
	"fmt"
	"math"
)

// SMILE type tokens (common ones)
const (
	smileHeaderByte1 = 0x3A // ':'
	smileHeaderByte2 = 0x29 // ')'
	smileHeaderByte3 = 0x0A // '\n'

	// String markers
	smileStringShort = 0x20 // Short ASCII string (0x20-0x3F, then 0x40-0x7F for shared)
	smileStringLong  = 0xE0 // Long ASCII string

	// Number markers
	smileInt32  = 0x24 // 32-bit integer
	smileInt64  = 0x28 // 64-bit integer
	smileFloat  = 0x28 // Float/double
	smileBigDec = 0x2A // Big decimal

	// Structural markers
	smileStartObject = 0xFA // Start object (map)
	smileEndObject   = 0xFB // End object
	smileStartArray  = 0xF8 // Start array
	smileEndArray    = 0xF9 // End array

	// Literal values
	smileNull  = 0x21 // null
	smileFalse = 0x22 // false
	smileTrue  = 0x23 // true
)

// DecodeSMILE attempts to decode SMILE-encoded binary JSON data
// SMILE is a binary JSON format: http://wiki.fasterxml.com/SmileFormat
func DecodeSMILE(data []byte) (map[string]interface{}, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short to be SMILE format")
	}

	// Check for SMILE header: :)\n
	offset := 0
	if data[0] == smileHeaderByte1 && data[1] == smileHeaderByte2 && data[2] == smileHeaderByte3 {
		offset = 4 // Skip header including version byte
	}

	// Try to decode
	decoder := &smileDecoder{
		data:   data,
		offset: offset,
		keys:   make([]string, 0, 64), // Shared key buffer
	}

	// Attempt basic decoding
	result, err := decoder.decode()
	if err != nil {
		// Fall back to string extraction but try to get partial data
		partial, _ := extractBasicInfo(data)
		partial["_decodeError"] = err.Error()
		partial["_decodedOffset"] = decoder.offset

		// If we got a partial object, merge it in
		if m, ok := result.(map[string]interface{}); ok && len(m) > 0 {
			partial["_partiallyDecoded"] = m
			return partial, nil // Return success with partial data
		}

		return partial, fmt.Errorf("SMILE decoding incomplete: %v", err)
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	// If result is not a map, wrap it
	return map[string]interface{}{"_data": result}, nil
}

type smileDecoder struct {
	data   []byte
	offset int
	keys   []string // Shared key names
}

func (d *smileDecoder) decode() (interface{}, error) {
	if d.offset >= len(d.data) {
		return nil, fmt.Errorf("unexpected end of data")
	}

	b := d.data[d.offset]

	// Structural markers
	if b == smileStartObject {
		return d.readObject()
	} else if b == smileStartArray {
		return d.readArray()
	} else if b == smileEndObject || b == smileEndArray {
		// These should be handled by their respective read functions
		return nil, fmt.Errorf("unexpected end marker: 0x%02x", b)
	}

	// Literal values
	if b == smileTrue {
		d.offset++
		return true, nil
	} else if b == smileFalse {
		d.offset++
		return false, nil
	} else if b == smileNull {
		d.offset++
		return nil, nil
	}

	// Strings
	if b < 0x20 {
		// Shared string reference (short)
		return d.readSharedString()
	} else if b >= 0x20 && b < 0x40 {
		// Tiny ASCII (length in low 5 bits)
		return d.readTinyAscii()
	} else if b >= 0x40 && b < 0xC0 {
		// Short ASCII (includes 0x40-0x7F and 0x80-0xBF)
		return d.readShortAscii()
	}

	// Integers
	if b >= 0xC0 && b < 0xE0 {
		// Small integers (-16 to 15)
		return d.readSmallInt()
	}

	// Long strings and other types
	if b >= 0xE0 {
		return d.readLongValue()
	}

	// Special numeric types
	if b == 0x24 {
		// 32-bit integer
		return d.readInt32()
	} else if b == 0x25 {
		// 64-bit integer
		return d.readInt64()
	} else if b == 0x26 {
		// BigInteger
		return d.readBigInteger()
	} else if b == 0x28 || b == 0x2A {
		// Float or Double
		return d.readFloat()
	}

	return nil, fmt.Errorf("unknown SMILE token: 0x%02x at offset %d", b, d.offset)
}

func (d *smileDecoder) readObject() (map[string]interface{}, error) {
	d.offset++ // Skip 0xFA
	result := make(map[string]interface{})

	for d.offset < len(d.data) {
		if d.offset >= len(d.data) {
			return result, fmt.Errorf("unexpected end in object")
		}

		b := d.data[d.offset]

		if b == smileEndObject {
			d.offset++
			return result, nil
		}

		// Sanity check: if we hit end array, something is wrong
		if b == smileEndArray {
			return result, fmt.Errorf("unexpected end array in object at offset %d", d.offset)
		}

		// Read key - in SMILE, keys are strings (can be shared refs or new strings)
		var keyStr string
		var err error

		// Keys can be:
		// - Tiny ASCII (0x20-0x3F)
		// - Short ASCII (0x40-0x7F) - NOT added to shared keys
		// - Short ASCII shared (0x80-0xBF) - ADDED to shared keys
		// - Short shared reference (0x00-0x1F) - reference to existing key
		// - Long ASCII (0xE0+)

		if b < 0x20 {
			// Shared key reference (0x00-0x1F)
			keyStr, err = d.readSharedString()
			if err != nil {
				return result, fmt.Errorf("error reading shared key: %v", err)
			}
		} else if b >= 0x20 && b < 0x40 {
			// Tiny ASCII key
			keyStr, err = d.readTinyAscii()
			if err != nil {
				return result, fmt.Errorf("error reading tiny key: %v", err)
			}
		} else if b >= 0x40 && b < 0xC0 {
			// Short ASCII key (includes both 0x40-0x7F and 0x80-0xBF ranges)
			// The 0x80-0xBF range will be added to shared keys automatically
			keyStr, err = d.readShortAscii()
			if err != nil {
				return result, fmt.Errorf("error reading object key: %v", err)
			}
		} else if b >= 0xE0 {
			// Long string
			val, err := d.readLongValue()
			if err != nil {
				return result, fmt.Errorf("error reading long key: %v", err)
			}
			if s, ok := val.(string); ok {
				keyStr = s
			} else {
				keyStr = fmt.Sprintf("%v", val)
			}
		} else {
			return result, fmt.Errorf("unexpected key type marker: 0x%02x at offset %d", b, d.offset)
		}

		// Read value
		value, err := d.decode()
		if err != nil {
			// Store partial result with error indication, then stop
			result[keyStr] = fmt.Sprintf("<decode error: %v>", err)
			// Return what we have so far
			return result, err
		}

		result[keyStr] = value
	}

	return result, fmt.Errorf("object not properly closed")
}

func (d *smileDecoder) readArray() ([]interface{}, error) {
	d.offset++ // Skip 0xF8
	result := make([]interface{}, 0)

	for d.offset < len(d.data) {
		if d.data[d.offset] == smileEndArray {
			d.offset++
			return result, nil
		}

		value, err := d.decode()
		if err != nil {
			return result, err
		}

		result = append(result, value)
	}

	return result, nil
}

func (d *smileDecoder) readTinyAscii() (string, error) {
	b := d.data[d.offset]
	length := int(b - 0x20 + 1)
	d.offset++

	if d.offset+length > len(d.data) {
		return "", fmt.Errorf("string extends beyond data")
	}

	str := string(d.data[d.offset : d.offset+length])
	d.offset += length
	return str, nil
}

func (d *smileDecoder) readShortAscii() (string, error) {
	b := d.data[d.offset]
	var length int

	// Short ASCII strings: 0x40-0xBF
	// 0x40-0x7F: length = byte - 0x40 (0-63 bytes)
	// 0x80-0xBF: length = byte - 0x80 + 1 (1-64 bytes) - these add to shared keys!
	if b >= 0x40 && b < 0x80 {
		length = int(b - 0x40)
	} else if b >= 0x80 && b < 0xC0 {
		length = int(b - 0x80 + 1)
	} else {
		return "", fmt.Errorf("not a short ASCII string: 0x%02x", b)
	}

	d.offset++

	if d.offset+length > len(d.data) {
		return "", fmt.Errorf("string extends beyond data")
	}

	str := string(d.data[d.offset : d.offset+length])
	d.offset += length

	// Strings in 0x80-0xBF range are added to shared keys
	// (these are "long key names")
	if b >= 0x80 && b < 0xC0 {
		d.keys = append(d.keys, str)
	}

	return str, nil
}

func (d *smileDecoder) readShortString() (string, error) {
	b := d.data[d.offset]

	// Check if it's a shared key reference (0x80-0xBF)
	if b >= 0x80 && b < 0xC0 {
		idx := int(b - 0x80)
		d.offset++
		if idx < len(d.keys) {
			return d.keys[idx], nil
		}
		return fmt.Sprintf("<key#%d>", idx), nil
	}

	d.offset++
	return "", fmt.Errorf("unhandled short string type: 0x%02x", b)
}

func (d *smileDecoder) readSharedString() (string, error) {
	idx := int(d.data[d.offset])
	d.offset++
	if idx < len(d.keys) {
		return d.keys[idx], nil
	}
	return fmt.Sprintf("<shared#%d>", idx), nil
}

func (d *smileDecoder) readLongString() (string, error) {
	b := d.data[d.offset]
	d.offset++

	// 0xE0-0xE3: long ASCII strings with length following
	if b >= 0xE0 && b <= 0xE3 {
		// Read variable-length integer for string length
		length, err := d.readVInt()
		if err != nil {
			return "", err
		}

		if d.offset+length > len(d.data) {
			return "", fmt.Errorf("long string extends beyond data")
		}

		str := string(d.data[d.offset : d.offset+length])
		d.offset += length
		return str, nil
	}

	// 0xE4-0xE7: long Unicode strings
	if b >= 0xE4 && b <= 0xE7 {
		length, err := d.readVInt()
		if err != nil {
			return "", err
		}

		if d.offset+length > len(d.data) {
			return "", fmt.Errorf("long unicode string extends beyond data")
		}

		str := string(d.data[d.offset : d.offset+length])
		d.offset += length
		return str, nil
	}

	// 0xE8: Big integer
	// 0xE9-0xEB: floats/decimals
	// For now, try to skip unknown types
	return fmt.Sprintf("<type:0x%02x>", b-1), nil
}

// readVInt reads a variable-length integer (ZigZag encoding)
func (d *smileDecoder) readVInt() (int, error) {
	if d.offset >= len(d.data) {
		return 0, fmt.Errorf("unexpected end reading VInt")
	}

	// SMILE uses a continuation bit encoding
	result := 0
	shift := 0

	for {
		if d.offset >= len(d.data) {
			return 0, fmt.Errorf("unexpected end in VInt")
		}

		b := d.data[d.offset]
		d.offset++

		// Lower 7 bits contribute to value
		result |= int(b&0x7F) << shift
		shift += 7

		// High bit indicates continuation
		if b&0x80 == 0 {
			break
		}

		if shift > 28 {
			return 0, fmt.Errorf("VInt too large")
		}
	}

	return result, nil
}

// readLongValue handles 0xE0+ type markers
func (d *smileDecoder) readLongValue() (interface{}, error) {
	b := d.data[d.offset]

	// 0xE0-0xE3: long ASCII strings
	if b >= 0xE0 && b <= 0xE3 {
		return d.readLongString()
	}

	// 0xE4-0xE7: long Unicode strings
	if b >= 0xE4 && b <= 0xE7 {
		return d.readLongString()
	}

	// 0xE8: Big Integer
	if b == 0xE8 {
		return d.readBigInteger()
	}

	// 0xE9: Float (32-bit)
	if b == 0xE9 {
		return d.readFloat32()
	}

	// 0xEA: Double (64-bit)
	if b == 0xEA {
		return d.readFloat64()
	}

	// 0xEB: Big Decimal
	if b == 0xEB {
		return d.readBigDecimal()
	}

	// Unknown type - skip it
	d.offset++
	return fmt.Sprintf("<unknown:0x%02x>", b), nil
}

// readInt32 reads a 32-bit integer
func (d *smileDecoder) readInt32() (int32, error) {
	d.offset++ // Skip type marker

	if d.offset+4 > len(d.data) {
		return 0, fmt.Errorf("not enough data for int32")
	}

	// Big-endian
	value := int32(d.data[d.offset])<<24 |
		int32(d.data[d.offset+1])<<16 |
		int32(d.data[d.offset+2])<<8 |
		int32(d.data[d.offset+3])
	d.offset += 4

	return value, nil
}

// readInt64 reads a 64-bit integer
func (d *smileDecoder) readInt64() (int64, error) {
	d.offset++ // Skip type marker

	if d.offset+8 > len(d.data) {
		return 0, fmt.Errorf("not enough data for int64")
	}

	// Big-endian
	value := int64(d.data[d.offset])<<56 |
		int64(d.data[d.offset+1])<<48 |
		int64(d.data[d.offset+2])<<40 |
		int64(d.data[d.offset+3])<<32 |
		int64(d.data[d.offset+4])<<24 |
		int64(d.data[d.offset+5])<<16 |
		int64(d.data[d.offset+6])<<8 |
		int64(d.data[d.offset+7])
	d.offset += 8

	return value, nil
}

// readBigInteger reads a variable-length big integer
func (d *smileDecoder) readBigInteger() (string, error) {
	d.offset++ // Skip type marker

	// Read length
	length, err := d.readVInt()
	if err != nil {
		return "", err
	}

	if d.offset+length > len(d.data) {
		return "", fmt.Errorf("big integer extends beyond data")
	}

	// For now, return as hex string
	bytes := d.data[d.offset : d.offset+length]
	d.offset += length

	return fmt.Sprintf("<bigint:%x>", bytes), nil
}

// readFloat32 reads a 32-bit float
func (d *smileDecoder) readFloat32() (float32, error) {
	d.offset++ // Skip type marker

	if d.offset+4 > len(d.data) {
		return 0, fmt.Errorf("not enough data for float32")
	}

	// Read as uint32 first, then convert
	bits := uint32(d.data[d.offset])<<24 |
		uint32(d.data[d.offset+1])<<16 |
		uint32(d.data[d.offset+2])<<8 |
		uint32(d.data[d.offset+3])
	d.offset += 4

	// Convert bits to float
	return float32FromBits(bits), nil
}

// readFloat64 reads a 64-bit double
func (d *smileDecoder) readFloat64() (float64, error) {
	d.offset++ // Skip type marker

	if d.offset+8 > len(d.data) {
		return 0, fmt.Errorf("not enough data for float64")
	}

	// Read as uint64 first
	bits := uint64(d.data[d.offset])<<56 |
		uint64(d.data[d.offset+1])<<48 |
		uint64(d.data[d.offset+2])<<40 |
		uint64(d.data[d.offset+3])<<32 |
		uint64(d.data[d.offset+4])<<24 |
		uint64(d.data[d.offset+5])<<16 |
		uint64(d.data[d.offset+6])<<8 |
		uint64(d.data[d.offset+7])
	d.offset += 8

	// Convert bits to float
	return float64FromBits(bits), nil
}

// readBigDecimal reads a big decimal value
func (d *smileDecoder) readBigDecimal() (string, error) {
	d.offset++ // Skip type marker

	// Big decimal has scale + value
	scale, err := d.readVInt()
	if err != nil {
		return "", err
	}

	// Read the unscaled value length
	length, err := d.readVInt()
	if err != nil {
		return "", err
	}

	if d.offset+length > len(d.data) {
		return "", fmt.Errorf("big decimal extends beyond data")
	}

	bytes := d.data[d.offset : d.offset+length]
	d.offset += length

	return fmt.Sprintf("<decimal:scale=%d,val=%x>", scale, bytes), nil
}

// readFloat is a generic float reader
func (d *smileDecoder) readFloat() (interface{}, error) {
	b := d.data[d.offset]

	if b == 0x28 {
		// Could be float32
		return d.readFloat32()
	} else if b == 0x2A {
		// Could be float64
		return d.readFloat64()
	}

	d.offset++
	return nil, fmt.Errorf("unknown float type: 0x%02x", b)
}

// Float conversion helpers (using math package functions)
func float32FromBits(bits uint32) float32 {
	return math.Float32frombits(bits)
}

func float64FromBits(bits uint64) float64 {
	return math.Float64frombits(bits)
}

func (d *smileDecoder) readSmallInt() (interface{}, error) {
	b := d.data[d.offset]
	d.offset++

	// SMILE small integers: 0xC0-0xDF range
	// 0xC0-0xDF: integers from -16 to 15
	if b >= 0xC0 && b <= 0xDF {
		value := int(b) - 0xD0 // 0xD0 is zero
		return value, nil
	}

	// If it's 0xE0 or higher, it's a different type
	return nil, fmt.Errorf("not a small integer: 0x%02x", b)
}

// extractBasicInfo falls back to string extraction when full decoding fails
func extractBasicInfo(data []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["_smileEncoded"] = true
	result["_dataSize"] = len(data)

	// Extract readable strings
	strings := extractStrings(data, 4)
	if len(strings) > 0 {
		result["_extractedStrings"] = strings[:min(20, len(strings))]
	}

	// Try to extract key-value pairs from the strings
	info := make(map[string]interface{})
	for i := 0; i < len(strings)-1; i++ {
		key := strings[i]
		// Common field names
		if isLikelyFieldName(key) && i+1 < len(strings) {
			value := strings[i+1]
			info[key] = value
		}
	}

	if len(info) > 0 {
		result["_partialData"] = info
	}

	return result, fmt.Errorf("full SMILE decoding not implemented; partial data extracted")
}

// isLikelyFieldName checks if a string looks like a field name
func isLikelyFieldName(s string) bool {
	if len(s) < 3 || len(s) > 30 {
		return false
	}
	// Check for common patterns
	commonFields := []string{
		"matchlen", "flags", "date", "name", "player",
		"event", "location", "round", "comment", "site",
		"rating", "rank", "score", "points", "games",
	}
	for _, field := range commonFields {
		if s == field || contains(s, field) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)
}

// extractStrings finds printable ASCII strings in binary data
func extractStrings(data []byte, minLen int) []string {
	var strings []string
	var current []byte

	for _, b := range data {
		// Check if byte is printable ASCII (space to ~)
		if b >= 32 && b <= 126 {
			current = append(current, b)
		} else {
			if len(current) >= minLen {
				strings = append(strings, string(current))
			}
			current = nil
		}
	}

	// Don't forget the last string
	if len(current) >= minLen {
		strings = append(strings, string(current))
	}

	return strings
}

// attemptSMILEDecode tries various strategies to decode SMILE data
func attemptSMILEDecode(data []byte) (map[string]interface{}, error) {
	// Strategy 1: Check if it's actually JSON (some files might not use SMILE)
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err == nil {
		return result, nil
	}

	// Strategy 2: Try SMILE decoding
	return DecodeSMILE(data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
