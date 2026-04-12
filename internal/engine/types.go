// Package engine provides core types and the algorithm registry
// for the LogicFlow algorithm visualizer.
package engine

// Step represents a single step in the algorithm's execution trace.
// Each step captures the array state, which indices are highlighted,
// and what type of action was performed.
type Step struct {
	CurrentState []int  `json:"current_state"`
	Highlights   []int  `json:"highlights"`
	ActionType   string `json:"action_type"` // "compare", "swap", "partition", "merge", "insert", "shift", "sorted", "check", "found", "not_found", "jump", "eliminate"
}

// ExecuteParams contains all parameters needed to run an algorithm.
// Sorting algorithms use only Array; searching algorithms also use Target.
type ExecuteParams struct {
	Array  []int
	Target int
}

// AlgorithmRequest is the incoming JSON payload from the client.
type AlgorithmRequest struct {
	Algorithm string `json:"algorithm"`
	Array     []int  `json:"array"`
	Target    *int   `json:"target,omitempty"` // required for searching algorithms
}

// AlgorithmInfo describes a registered algorithm for the discovery endpoint.
type AlgorithmInfo struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Category       string `json:"category"` // "sorting", "searching"
	TimeComplexity string `json:"time_complexity"`
	Description    string `json:"description"`
}

// AlgorithmResponse is the JSON response returned to the client
// containing the full execution trace and metadata.
type AlgorithmResponse struct {
	Steps    []Step   `json:"steps"`
	Metadata Metadata `json:"metadata"`
}

// Metadata contains performance metrics about the algorithm execution.
type Metadata struct {
	ExecutionTimeMicroseconds int64  `json:"execution_time_us"`
	Comparisons               int    `json:"comparisons"`
	Operations                int    `json:"operations"`     // swaps for sorting, checks for searching
	TimeComplexity            string `json:"time_complexity"`
	AlgorithmName             string `json:"algorithm_name"`
	Category                  string `json:"category"`
	FoundIndex                *int   `json:"found_index,omitempty"` // for searching: index where target was found (nil if N/A, -1 if not found)
}
