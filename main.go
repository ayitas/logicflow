// LogicFlow — Algorithm Visualizer
//
// A high-performance algorithm visualization server built with Go.
// The server exposes a REST API for running algorithms (sorting, searching, etc.)
// and serves the frontend static files.
//
// Architecture:
//   - internal/engine/    → Core types and algorithm registry
//   - internal/algorithm/ → Algorithm implementations (plugin pattern)
//   - internal/handler/   → HTTP route handlers
//   - static/             → Frontend (HTML + CSS + JS + D3.js)
//
// To add a new algorithm, simply create a new file in internal/algorithm/
// that implements the engine.Algorithm interface and registers via init().

package main

import (
	"fmt"
	"log"
	"net/http"

	// Import algorithm package for its init() side effects.
	// Each algorithm file registers itself with the engine registry.
	_ "logicflow/internal/algorithm"

	"logicflow/internal/engine"
	"logicflow/internal/handler"
)

func main() {
	const port = ":8080"

	// Log registered algorithms grouped by category
	algos := engine.List()
	fmt.Printf("LogicFlow — Algorithm Visualizer\n")
	fmt.Printf("   Registered algorithms: %d\n\n", len(algos))

	currentCategory := ""
	for _, a := range algos {
		if a.Category != currentCategory {
			currentCategory = a.Category
			fmt.Printf("   [%s]\n", currentCategory)
		}
		fmt.Printf("   - %-25s %s\n", a.DisplayName, a.TimeComplexity)
	}
	fmt.Println()

	// API routes
	http.HandleFunc("/execute", corsMiddleware(handler.ExecuteHandler))
	http.HandleFunc("/sort", corsMiddleware(handler.ExecuteHandler)) // backward compatibility
	http.HandleFunc("/algorithms", corsMiddleware(handler.AlgorithmsHandler))

	// Serve frontend static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	fmt.Printf("   Server listening on http://localhost%s\n\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// corsMiddleware adds CORS headers to allow frontend development on different ports.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
