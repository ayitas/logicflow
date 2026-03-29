// Package handler provides HTTP handlers for the LogicFlow API.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"logicflow/internal/engine"
)

// SortHandler handles POST /sort requests.
// It validates the request, runs the selected algorithm, and returns
// the step-by-step execution trace with performance metadata.
func SortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req engine.SortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if len(req.Array) == 0 {
		http.Error(w, `{"error":"array must not be empty"}`, http.StatusBadRequest)
		return
	}
	if len(req.Array) > 200 {
		http.Error(w, `{"error":"array size must be <= 200"}`, http.StatusBadRequest)
		return
	}

	// Look up the algorithm from the registry
	algo, ok := engine.Get(req.Algorithm)
	if !ok {
		http.Error(w, `{"error":"unknown algorithm: `+req.Algorithm+`"}`, http.StatusBadRequest)
		return
	}

	// Execute and measure time
	start := time.Now()
	steps, comparisons, swapsMoves := algo.Execute(req.Array)
	elapsed := time.Since(start)

	resp := engine.SortResponse{
		Steps: steps,
		Metadata: engine.Metadata{
			ExecutionTimeMicroseconds: elapsed.Microseconds(),
			Comparisons:               comparisons,
			SwapsMoves:                swapsMoves,
			TimeComplexity:            algo.TimeComplexity(),
			AlgorithmName:             algo.DisplayName(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// AlgorithmsHandler handles GET /algorithms requests.
// Returns a list of all registered algorithms so the frontend
// can dynamically populate the algorithm selector.
func AlgorithmsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	algos := engine.List()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(algos); err != nil {
		log.Printf("Error encoding algorithms list: %v", err)
	}
}
