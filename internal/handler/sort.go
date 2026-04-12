// Package handler provides HTTP handlers for the LogicFlow API.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"logicflow/internal/engine"
)

// ExecuteHandler handles POST /execute requests.
// It validates the request, runs the selected algorithm, and returns
// the step-by-step execution trace with performance metadata.
func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req engine.AlgorithmRequest
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

	// Searching algorithms require a target value
	if algo.Category() == "searching" && req.Target == nil {
		http.Error(w, `{"error":"target value is required for searching algorithms"}`, http.StatusBadRequest)
		return
	}

	// Build execute params
	params := engine.ExecuteParams{
		Array: req.Array,
	}
	if req.Target != nil {
		params.Target = *req.Target
	}

	// Execute and measure time
	start := time.Now()
	steps, comparisons, operations := algo.Execute(params)
	elapsed := time.Since(start)

	// Determine found_index for searching algorithms
	var foundIndex *int
	if algo.Category() == "searching" {
		idx := -1
		if len(steps) > 0 {
			lastStep := steps[len(steps)-1]
			if lastStep.ActionType == "found" && len(lastStep.Highlights) > 0 {
				idx = lastStep.Highlights[0]
			}
		}
		foundIndex = &idx
	}

	resp := engine.AlgorithmResponse{
		Steps: steps,
		Metadata: engine.Metadata{
			ExecutionTimeMicroseconds: elapsed.Microseconds(),
			Comparisons:               comparisons,
			Operations:                operations,
			TimeComplexity:            algo.TimeComplexity(),
			AlgorithmName:             algo.DisplayName(),
			Category:                  algo.Category(),
			FoundIndex:                foundIndex,
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
