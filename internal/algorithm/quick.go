package algorithm

import (
	"logicflow/internal/engine"
)

// QuickSort implements the Quick Sort algorithm using Lomuto partition.
// Time: O(n log n) avg, O(n²) worst | Space: O(log n) | Stable: No
type QuickSort struct{}

func init() {
	engine.Register(&QuickSort{})
}

func (q *QuickSort) Name() string           { return "quick_sort" }
func (q *QuickSort) DisplayName() string    { return "Quick Sort" }
func (q *QuickSort) Category() string       { return "sorting" }
func (q *QuickSort) TimeComplexity() string { return "O(n log n)" }
func (q *QuickSort) Description() string {
	return "Divide and Conquer algorithm that selects a pivot and partitions the array around it. Very fast in practice with O(n log n) average time."
}

func (q *QuickSort) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	tracker := &quickTracker{
		steps:       make([]engine.Step, 0, len(data)*len(data)),
		comparisons: 0,
		swaps:       0,
	}

	quickSortRecursive(data, 0, len(data)-1, tracker)

	tracker.steps = append(tracker.steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   makeRange(0, len(data)-1),
		ActionType:   "sorted",
	})

	return tracker.steps, tracker.comparisons, tracker.swaps
}

type quickTracker struct {
	steps       []engine.Step
	comparisons int
	swaps       int
}

func quickSortRecursive(data []int, low, high int, t *quickTracker) {
	if low >= high {
		if low == high {
			t.steps = append(t.steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{low},
				ActionType:   "sorted",
			})
		}
		return
	}

	pivotIdx := lomutoPartition(data, low, high, t)

	t.steps = append(t.steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{pivotIdx},
		ActionType:   "sorted",
	})

	quickSortRecursive(data, low, pivotIdx-1, t)
	quickSortRecursive(data, pivotIdx+1, high, t)
}

func lomutoPartition(data []int, low, high int, t *quickTracker) int {
	pivot := data[high]

	t.steps = append(t.steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{high},
		ActionType:   "partition",
	})

	i := low - 1

	for j := low; j < high; j++ {
		t.comparisons++

		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{j, high},
			ActionType:   "compare",
		})

		if data[j] <= pivot {
			i++
			if i != j {
				data[i], data[j] = data[j], data[i]
				t.swaps++
				t.steps = append(t.steps, engine.Step{
					CurrentState: engine.SnapshotArray(data),
					Highlights:   []int{i, j},
					ActionType:   "swap",
				})
			}
		}
	}

	i++
	if i != high {
		data[i], data[high] = data[high], data[i]
		t.swaps++
		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i, high},
			ActionType:   "swap",
		})
	}

	return i
}
