package main

import (
	"fmt"
	"os"

	"github.com/unger/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_decode <file.bgf>")
		os.Exit(1)
	}

	match, err := bgfparser.ParseBGF(os.Args[1])
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
	}

	fmt.Printf("Decoding Warning: %s\n", match.DecodingWarning)
	fmt.Printf("Data keys: %v\n", getKeys(match.Data))

	for k, v := range match.Data {
		fmt.Printf("%s => %T: %v\n", k, v, v)
	}
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
