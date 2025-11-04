// Package main provides a web server example for parsing BGBlitz filespackage webserver

// This example demonstrates how to use the bgfparser in HTTP handlers
// for uploading and analyzing BGF and TXT files via a web interface.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kevung/bgfparser"
)

// MatchSummary provides a simplified match summary for API responses
type MatchSummary struct {
	Format      string                 `json:"format"`
	Version     string                 `json:"version"`
	Compressed  bool                   `json:"compressed"`
	UseSmile    bool                   `json:"use_smile"`
	MatchInfo   map[string]interface{} `json:"match_info,omitempty"`
	DataKeys    []string               `json:"data_keys,omitempty"`
	DataPreview map[string]interface{} `json:"data_preview,omitempty"`
}

// PositionSummary provides a simplified position summary
type PositionSummary struct {
	PlayerX        string  `json:"player_x"`
	PlayerO        string  `json:"player_o"`
	Score          string  `json:"score"`
	MatchLength    int     `json:"match_length"`
	OnRoll         string  `json:"on_roll"`
	Dice           [2]int  `json:"dice"`
	CubeValue      int     `json:"cube_value"`
	NumEvals       int     `json:"num_evaluations"`
	HasCubeDecis   bool    `json:"has_cube_decision"`
	BestMove       string  `json:"best_move,omitempty"`
	BestMoveEquity float64 `json:"best_move_equity,omitempty"`
}

// uploadBGFHandler handles BGF file uploads
func uploadBGFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10 MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("bgffile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file into memory
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Parse the BGF file
	match, err := bgfparser.ParseBGFFromReader(bytes.NewReader(fileData))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse BGF file: %v", err), http.StatusBadRequest)
		return
	}

	// Create summary
	info := match.GetMatchInfo()
	summary := MatchSummary{
		Format:     match.Format,
		Version:    match.Version,
		Compressed: match.Compress,
		UseSmile:   match.UseSmile,
		MatchInfo:  info,
	}

	// Add data keys preview
	if match.Data != nil {
		keys := make([]string, 0, len(match.Data))
		preview := make(map[string]interface{})
		count := 0
		for k, v := range match.Data {
			keys = append(keys, k)
			if count < 5 {
				preview[k] = v
				count++
			}
		}
		summary.DataKeys = keys
		summary.DataPreview = preview
	}

	// Add filename to summary
	summary.MatchInfo["filename"] = header.Filename

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// fullBGFHandler returns the complete match as JSON
func fullBGFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file
	file, _, err := r.FormFile("bgffile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read into memory
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Parse
	match, err := bgfparser.ParseBGFFromReader(bytes.NewReader(fileData))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse BGF file: %v", err), http.StatusBadRequest)
		return
	}

	// Return full match as JSON
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := match.ToJSON()
	w.Write(jsonData)
}

// uploadTXTHandler handles TXT file uploads
func uploadTXTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("txtfile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse the TXT file
	pos, err := bgfparser.ParseTXTFromReader(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse TXT file: %v", err), http.StatusBadRequest)
		return
	}

	// Create summary
	summary := PositionSummary{
		PlayerX:      pos.PlayerX,
		PlayerO:      pos.PlayerO,
		Score:        fmt.Sprintf("%s: %d - %s: %d", pos.PlayerO, pos.ScoreO, pos.PlayerX, pos.ScoreX),
		MatchLength:  pos.MatchLength,
		OnRoll:       pos.OnRoll,
		Dice:         pos.Dice,
		CubeValue:    pos.CubeValue,
		NumEvals:     len(pos.Evaluations),
		HasCubeDecis: pos.CubeDecision != nil,
	}

	// Add best move info
	for _, eval := range pos.Evaluations {
		if eval.IsBest {
			summary.BestMove = eval.Move
			summary.BestMoveEquity = eval.Equity
			break
		}
	}

	// Add filename
	response := map[string]interface{}{
		"filename": header.Filename,
		"summary":  summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fullTXTHandler returns the complete position as JSON
func fullTXTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("txtfile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	pos, err := bgfparser.ParseTXTFromReader(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse TXT file: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := pos.ToJSON()
	w.Write(jsonData)
}

// homeHandler serves a simple HTML form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>BGF Parser Web Interface</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        h2 { color: #666; margin-top: 30px; }
        form { margin: 20px 0; }
        input[type="file"] { margin-right: 10px; }
        button { background: #4CAF50; color: white; padding: 10px 20px; border: none; cursor: pointer; }
        button:hover { background: #45a049; }
        iframe { margin-top: 10px; }
        .section { border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>BGBlitz Parser Web Interface</h1>
    
    <div class="section">
        <h2>BGF Files (Match Files)</h2>
        
        <h3>Quick Summary</h3>
        <form action="/upload/bgf" method="post" enctype="multipart/form-data" target="bgf_summary">
            <input type="file" name="bgffile" accept=".bgf" required>
            <button type="submit">Get BGF Summary</button>
        </form>
        <iframe name="bgf_summary" style="width:100%; height:250px; border:1px solid #ccc;"></iframe>

        <h3>Full Match JSON</h3>
        <form action="/full/bgf" method="post" enctype="multipart/form-data" target="bgf_full">
            <input type="file" name="bgffile" accept=".bgf" required>
            <button type="submit">Get Full BGF Match</button>
        </form>
        <iframe name="bgf_full" style="width:100%; height:400px; border:1px solid #ccc;"></iframe>
    </div>

    <div class="section">
        <h2>TXT Files (Position Files)</h2>
        
        <h3>Quick Summary</h3>
        <form action="/upload/txt" method="post" enctype="multipart/form-data" target="txt_summary">
            <input type="file" name="txtfile" accept=".txt" required>
            <button type="submit">Get TXT Summary</button>
        </form>
        <iframe name="txt_summary" style="width:100%; height:250px; border:1px solid #ccc;"></iframe>

        <h3>Full Position JSON</h3>
        <form action="/full/txt" method="post" enctype="multipart/form-data" target="txt_full">
            <input type="file" name="txtfile" accept=".txt" required>
            <button type="submit">Get Full TXT Position</button>
        </form>
        <iframe name="txt_full" style="width:100%; height:400px; border:1px solid #ccc;"></iframe>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// healthHandler provides a health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "bgfparser-web",
	})
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload/bgf", uploadBGFHandler)
	http.HandleFunc("/full/bgf", fullBGFHandler)
	http.HandleFunc("/upload/txt", uploadTXTHandler)
	http.HandleFunc("/full/txt", fullTXTHandler)
	http.HandleFunc("/health", healthHandler)

	port := ":8080"
	fmt.Printf("BGBlitz Parser Web Server\n")
	fmt.Printf("=========================\n\n")
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Printf("\nEndpoints:\n")
	fmt.Printf("  GET  /              - Web interface\n")
	fmt.Printf("  POST /upload/bgf    - Upload BGF file (summary)\n")
	fmt.Printf("  POST /full/bgf      - Upload BGF file (full JSON)\n")
	fmt.Printf("  POST /upload/txt    - Upload TXT file (summary)\n")
	fmt.Printf("  POST /full/txt      - Upload TXT file (full JSON)\n")
	fmt.Printf("  GET  /health        - Health check\n")
	fmt.Printf("\nUpload BGF or TXT files to analyze positions and matches\n\n")

	log.Fatal(http.ListenAndServe(port, nil))
}
