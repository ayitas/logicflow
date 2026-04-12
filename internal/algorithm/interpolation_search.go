package algorithm

import (
	"sort"

	"logicflow/internal/engine"
)

// InterpolationSearch implements the Interpolation Search algorithm.
// Time: O(log log n) avg for uniform data, O(n) worst | Space: O(1)
// Requires a sorted, uniformly distributed array.
type InterpolationSearch struct{}

func init() {
	engine.Register(&InterpolationSearch{})
}

func (is *InterpolationSearch) Name() string           { return "interpolation_search" }
func (is *InterpolationSearch) DisplayName() string    { return "Interpolation Search" }
func (is *InterpolationSearch) Category() string       { return "searching" }
func (is *InterpolationSearch) TimeComplexity() string { return "O(log log n)" }
func (is *InterpolationSearch) Description() string {
	return "Estimates the position of the target using interpolation formula based on the value distribution. Faster than binary search for uniformly distributed, sorted data. Falls back to O(n) for skewed distributions."
}

func (is *InterpolationSearch) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	target := params.Target
	n := len(data)

	// Interpolation search requires sorted array
	sort.Ints(data)

	steps := make([]engine.Step, 0, 32)
	comparisons := 0
	checks := 0

	// Show sorted array
	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   makeRange(0, n-1),
		ActionType:   "sorted",
	})

	low, high := 0, n-1

	for low <= high && target >= data[low] && target <= data[high] {
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

		// Calculate interpolated position
		var pos int
		if data[high] == data[low] {
			pos = low
		} else {
			pos = low + int(float64(high-low)*float64(target-data[low])/float64(data[high]-data[low]))
		}

		// Clamp position
		if pos < low {
			pos = low
		}
		if pos > high {
			pos = high
		}

		// Check the interpolated position
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{pos},
			ActionType:   "check",
		})

		if data[pos] == target {
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{pos},
				ActionType:   "found",
			})
			return steps, comparisons, checks
		}

		if data[pos] < target {
			// Eliminate left portion
			eliminateHighlights := make([]int, 0)
			for i := low; i <= pos; i++ {
				eliminateHighlights = append(eliminateHighlights, i)
			}
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   eliminateHighlights,
				ActionType:   "eliminate",
			})
			low = pos + 1
		} else {
			// Eliminate right portion
			eliminateHighlights := make([]int, 0)
			for i := pos; i <= high; i++ {
				eliminateHighlights = append(eliminateHighlights, i)
			}
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   eliminateHighlights,
				ActionType:   "eliminate",
			})
			high = pos - 1
		}
	}

	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{},
		ActionType:   "not_found",
	})

	return steps, comparisons, checks
}
