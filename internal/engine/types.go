// Package engine provides core types and the algorithm registry
// for the LogicFlow sorting algorithm visualizer.
package engine

// Step represents a single step in the algorithm's execution trace.
// Each step captures the array state, which indices are highlighted,
// and what type of action was performed.
type Step struct {
	CurrentState []int  `json:"current_state"`
	Highlights   []int  `json:"highlights"`
	ActionType   string `json:"action_type"` // "compare", "swap", "partition", "merge", "insert", "shift", "sorted"
}

// SortRequest is the incoming JSON payload from the client.
type SortRequest struct {
	Algorithm string `json:"algorithm"`
	Array     []int  `json:"array"`
}

// AlgorithmInfo describes a registered algorithm for the discovery endpoint.
type AlgorithmInfo struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	TimeComplexity string `json:"time_complexity"`
	Description    string `json:"description"`
}

// SortResponse is the JSON response returned to the client
// containing the full execution trace and metadata.
type SortResponse struct {
	Steps    []Step   `json:"steps"`
	Metadata Metadata `json:"metadata"`
}

// Metadata contains performance metrics about the sort execution.
type Metadata struct {
	ExecutionTimeMicroseconds int64  `json:"execution_time_us"`
	Comparisons               int    `json:"comparisons"`
	SwapsMoves                int    `json:"swaps_moves"`
	TimeComplexity            string `json:"time_complexity"`
	AlgorithmName             string `json:"algorithm_name"`
}
