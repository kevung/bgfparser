package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unger/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_bgf <filename.bgf>")
		fmt.Println("Example: parse_bgf matches/game1.bgf")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse the BGF file
	match, err := bgfparser.ParseBGF(filename)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	// Display match information
	fmt.Println("=== BGF Match File ===")
	fmt.Println(match.String())

	fmt.Println("\n=== Header Information ===")
	fmt.Printf("Format: %s\n", match.Format)
	fmt.Printf("Version: %s\n", match.Version)
	fmt.Printf("Compressed: %v\n", match.Compress)
	fmt.Printf("Uses SMILE encoding: %v\n", match.UseSmile)

	// Display any decoding warnings
	if match.DecodingWarning != "" {
		fmt.Printf("\n⚠️  Warning: %s\n", match.DecodingWarning)
		fmt.Println("\nNote: Full SMILE decoding requires additional libraries.")
		fmt.Println("However, we've extracted some information from the binary data:")
	}

	// Display extracted information
	if len(match.Data) > 0 {
		fmt.Println("\n=== Match Data ===")

		// Check if we have a partially decoded structure or full decoded data
		var decodedData map[string]interface{}

		if partialData, ok := match.Data["_partiallyDecoded"].(map[string]interface{}); ok {
			decodedData = partialData
		} else if _, hasError := match.Data["_decodeError"]; !hasError {
			// If no decode error, the Data itself contains the decoded fields
			decodedData = match.Data
		}

		if decodedData != nil && len(decodedData) > 0 {
			fmt.Println("\n--- All Decoded Fields ---")

			// Show all fields sorted by name for readability
			keys := make([]string, 0, len(decodedData))
			for key := range decodedData {
				if !strings.HasPrefix(key, "_") {
					keys = append(keys, key)
				}
			}

			// Display each field
			for _, key := range keys {
				val := decodedData[key]

				// Clean the key name for display
				cleanKey := ""
				for _, r := range key {
					if r >= 32 && r <= 126 {
						cleanKey += string(r)
					} else {
						cleanKey += fmt.Sprintf("\\x%02x", r)
					}
				}

				// Format the value nicely
				switch v := val.(type) {
				case string:
					// Clean up non-printable characters
					cleaned := ""
					hasNonPrintable := false
					for _, r := range v {
						if r >= 32 && r <= 126 {
							cleaned += string(r)
						} else {
							hasNonPrintable = true
						}
					}
					if cleaned != "" {
						if hasNonPrintable {
							fmt.Printf("  %s: %q (contains non-printable chars)\n", cleanKey, cleaned)
						} else {
							fmt.Printf("  %s: %s\n", cleanKey, cleaned)
						}
					} else {
						fmt.Printf("  %s: (empty or all non-printable)\n", cleanKey)
					}
				case map[string]interface{}:
					fmt.Printf("  %s: {\n", cleanKey)
					for subKey, subVal := range v {
						cleanSubKey := ""
						for _, r := range subKey {
							if r >= 32 && r <= 126 {
								cleanSubKey += string(r)
							} else {
								cleanSubKey += fmt.Sprintf("\\x%02x", r)
							}
						}
						fmt.Printf("    %s: %v\n", cleanSubKey, subVal)
					}
					fmt.Printf("  }\n")
				case []interface{}:
					fmt.Printf("  %s: [ %d elements ]\n", cleanKey, len(v))
					if len(v) > 0 && len(v) <= 10 {
						for i, elem := range v {
							fmt.Printf("    [%d]: %v\n", i, elem)
						}
					}
				default:
					fmt.Printf("  %s: %v\n", cleanKey, val)
				}
			}

			fmt.Printf("\nTotal decoded fields: %d\n", len(keys))

			if offset, ok := match.Data["_decodedOffset"].(int); ok {
				fmt.Printf("Decoded up to offset: %d bytes\n", offset)
			}
		}

		info := match.GetMatchInfo()
		showOtherInfo := false
		for key, value := range info {
			if key != "format" && key != "version" && key != "compress" && key != "useSmile" {
				if !showOtherInfo {
					fmt.Println("\n--- Other Information ---")
					showOtherInfo = true
				}
				fmt.Printf("%s: %v\n", key, value)
			}
		}

		// If SMILE encoded, show extracted strings
		if match.UseSmile {
			if strings, ok := match.Data["_extractedStrings"].([]string); ok && len(strings) > 0 {
				fmt.Println("\n=== Extracted Strings (from binary data) ===")
				for i, s := range strings {
					if len(s) > 50 {
						fmt.Printf("%d: %s...\n", i+1, s[:50])
					} else {
						fmt.Printf("%d: %s\n", i+1, s)
					}
				}
			}
			if size, ok := match.Data["_dataSize"].(int); ok {
				fmt.Printf("\nBinary data size: %d bytes\n", size)
			}
		}
	}
}
