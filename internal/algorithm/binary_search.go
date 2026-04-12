package algorithm

import (
	"sort"

	"logicflow/internal/engine"
)

// BinarySearch implements the Binary Search algorithm.
// Time: O(log n) | Space: O(1)
// Requires a sorted array. Divides search space in half each step.
type BinarySearch struct{}

func init() {
	engine.Register(&BinarySearch{})
}

func (b *BinarySearch) Name() string           { return "binary_search" }
func (b *BinarySearch) DisplayName() string    { return "Binary Search" }
func (b *BinarySearch) Category() string       { return "searching" }
func (b *BinarySearch) TimeComplexity() string { return "O(log n)" }
func (b *BinarySearch) Description() string {
	return "Divides the sorted array in half repeatedly, eliminating half the remaining elements each step. Requires a pre-sorted array. Extremely efficient for large datasets."
}

func (b *BinarySearch) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	target := params.Target

	// Binary search requires sorted array — sort it first
	sort.Ints(data)

	steps := make([]engine.Step, 0, 32)
	comparisons := 0
	checks := 0

	// Show the sorted array first
	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   makeRange(0, len(data)-1),
		ActionType:   "sorted",
	})

	low, high := 0, len(data)-1

	for low <= high {
		mid := low + (high-low)/2
		comparisons++
		checks++

		// Show the current search range
		rangeHighlights := make([]int, 0)
		for i := low; i <= high; i++ {
			rangeHighlights = append(rangeHighlights, i)
		}
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   rangeHighlights,
			ActionType:   "partition",
		})

		// Check mid element
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{mid},
			ActionType:   "check",
		})

		if data[mid] == target {
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{mid},
				ActionType:   "found",
			})
			return steps, comparisons, checks
		} else if data[mid] < target {
			// Eliminate left half
			eliminateHighlights := make([]int, 0)
			for i := low; i <= mid; i++ {
				eliminateHighlights = append(eliminateHighlights, i)
			}
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   eliminateHighlights,
				ActionType:   "eliminate",
			})
			low = mid + 1
		} else {
			// Eliminate right half
			eliminateHighlights := make([]int, 0)
			for i := mid; i <= high; i++ {
				eliminateHighlights = append(eliminateHighlights, i)
			}
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   eliminateHighlights,
				ActionType:   "eliminate",
			})
			high = mid - 1
		}
	}

	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{},
		ActionType:   "not_found",
	})

	return steps, comparisons, checks
}
