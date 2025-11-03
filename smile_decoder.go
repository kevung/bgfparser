package bgfparser

import (
	"encoding/json"
	"fmt"
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
		// Fall back to string extraction
		return extractBasicInfo(data)
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

	// Handle strings (most common: 0x20-0x3F for short ASCII, 0x40-0x7F for shared keys)
	if b >= 0x00 && b < 0x20 {
		// Shared string reference
		return d.readSharedString()
	} else if b >= 0x20 && b < 0x40 {
		// Tiny ASCII (length in low 5 bits)
		return d.readTinyAscii()
	} else if b >= 0x40 && b < 0x80 {
		// Short ASCII (length follows)
		return d.readShortAscii()
	} else if b >= 0x80 && b < 0xC0 {
		// Short Unicode or key reference
		return d.readShortString()
	} else if b == smileStartObject {
		return d.readObject()
	} else if b == smileStartArray {
		return d.readArray()
	} else if b == smileTrue {
		d.offset++
		return true, nil
	} else if b == smileFalse {
		d.offset++
		return false, nil
	} else if b == smileNull {
		d.offset++
		return nil, nil
	} else if b >= 0xC0 && b < 0xE0 {
		// Small integers
		return d.readSmallInt()
	} else if b >= 0xE0 {
		// Long strings or other types
		return d.readLongString()
	}

	return nil, fmt.Errorf("unknown SMILE token: 0x%02x at offset %d", b, d.offset)
}

func (d *smileDecoder) readObject() (map[string]interface{}, error) {
	d.offset++ // Skip 0xFA
	result := make(map[string]interface{})

	for d.offset < len(d.data) {
		if d.data[d.offset] == smileEndObject {
			d.offset++
			return result, nil
		}

		// Read key
		key, err := d.decode()
		if err != nil {
			return result, err
		}

		keyStr, ok := key.(string)
		if !ok {
			return result, fmt.Errorf("object key is not a string: %T", key)
		}

		// Read value
		value, err := d.decode()
		if err != nil {
			// Store partial result
			result[keyStr] = fmt.Sprintf("<decode error: %v>", err)
			continue
		}

		result[keyStr] = value
	}

	return result, nil
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
	length := int(b - 0x40)
	d.offset++

	if d.offset+length > len(d.data) {
		return "", fmt.Errorf("string extends beyond data")
	}

	str := string(d.data[d.offset : d.offset+length])
	d.offset += length

	// Add to shared keys if this looks like a key
	if length > 0 && length < 40 {
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
	// For now, extract what we can
	d.offset++
	// Try to find next structure marker
	start := d.offset
	for d.offset < len(d.data) && d.data[d.offset] >= 0x20 && d.data[d.offset] < 0x80 {
		d.offset++
	}
	return string(d.data[start:d.offset]), nil
}

func (d *smileDecoder) readSmallInt() (int, error) {
	b := d.data[d.offset]
	d.offset++

	// VInt encoding for small integers
	value := int(b & 0x1F)
	if b&0x20 == 0 {
		// Positive
		return value, nil
	}
	// Negative
	return -value, nil
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
