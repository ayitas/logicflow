package engine

import (
	"fmt"
	"sync"
)

// Algorithm is the interface that all sorting/algorithm implementations must satisfy.
// To add a new algorithm:
//  1. Create a new file in internal/algorithm/
//  2. Define a struct implementing this interface
//  3. Call engine.Register() in an init() function
//
// That's it — the algorithm will automatically appear in the API and frontend.
type Algorithm interface {
	// Name returns the unique identifier for this algorithm (e.g., "bubble_sort").
	Name() string

	// DisplayName returns a human-friendly name (e.g., "Bubble Sort").
	DisplayName() string

	// TimeComplexity returns the Big-O time complexity string (e.g., "O(n²)").
	TimeComplexity() string

	// Description returns a brief explanation of the algorithm.
	Description() string

	// Execute runs the algorithm on a copy of the input array and returns
	// the step-by-step execution trace along with comparison and swap counts.
	Execute(arr []int) (steps []Step, comparisons int, swapsMoves int)
}

// --- Global Registry ---

var (
	mu       sync.RWMutex
	registry = make(map[string]Algorithm)
)

// Register adds an algorithm to the global registry.
// It panics if an algorithm with the same name is already registered,
// preventing silent overwrites during init().
func Register(algo Algorithm) {
	mu.Lock()
	defer mu.Unlock()

	name := algo.Name()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("algorithm already registered: %s", name))
	}
	registry[name] = algo
}

// Get retrieves a registered algorithm by its name.
// Returns the algorithm and true if found, or nil and false otherwise.
func Get(name string) (Algorithm, bool) {
	mu.RLock()
	defer mu.RUnlock()

	algo, ok := registry[name]
	return algo, ok
}

// List returns information about all registered algorithms.
func List() []AlgorithmInfo {
	mu.RLock()
	defer mu.RUnlock()

	infos := make([]AlgorithmInfo, 0, len(registry))
	for _, algo := range registry {
		infos = append(infos, AlgorithmInfo{
			Name:           algo.Name(),
			DisplayName:    algo.DisplayName(),
			TimeComplexity: algo.TimeComplexity(),
			Description:    algo.Description(),
		})
	}
	return infos
}

// CopyArray creates a deep copy of an integer slice.
// Algorithms should use this to avoid mutating the original input.
func CopyArray(arr []int) []int {
	cp := make([]int, len(arr))
	copy(cp, arr)
	return cp
}

// SnapshotArray creates a snapshot of the current array state for a Step.
func SnapshotArray(arr []int) []int {
	return CopyArray(arr)
}
