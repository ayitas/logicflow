package algorithm

import (
	"logicflow/internal/engine"
)

// SelectionSort implements the Selection Sort algorithm.
// Time: O(n²) | Space: O(1) | Stable: No
// Minimizes the number of swaps — at most n-1 swaps.
type SelectionSort struct{}

func init() {
	engine.Register(&SelectionSort{})
}

func (s *SelectionSort) Name() string           { return "selection_sort" }
func (s *SelectionSort) DisplayName() string    { return "Selection Sort" }
func (s *SelectionSort) TimeComplexity() string { return "O(n²)" }
func (s *SelectionSort) Description() string {
	return "Finds the minimum element from the unsorted part and places it at the beginning. Minimizes the number of swaps performed."
}

func (s *SelectionSort) Execute(arr []int) ([]engine.Step, int, int) {
	data := engine.CopyArray(arr)
	n := len(data)
	steps := make([]engine.Step, 0, n*n)
	comparisons := 0
	swaps := 0

	for i := 0; i < n-1; i++ {
		minIdx := i

		for j := i + 1; j < n; j++ {
			// Compare current element with minimum
			comparisons++
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{minIdx, j},
				ActionType:   "compare",
			})

			if data[j] < data[minIdx] {
				minIdx = j
			}
		}

		// Swap the found minimum with the first unsorted element
		if minIdx != i {
			data[i], data[minIdx] = data[minIdx], data[i]
			swaps++
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{i, minIdx},
				ActionType:   "swap",
			})
		}

		// Mark element as sorted
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i},
			ActionType:   "sorted",
		})
	}

	// Mark last element as sorted
	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{n - 1},
		ActionType:   "sorted",
	})

	return steps, comparisons, swaps
}
