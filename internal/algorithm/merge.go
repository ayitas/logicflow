package algorithm

import (
	"logicflow/internal/engine"
)

// MergeSort implements the Merge Sort algorithm.
// Time: O(n log n) | Space: O(n) | Stable: Yes
type MergeSort struct{}

func init() {
	engine.Register(&MergeSort{})
}

func (m *MergeSort) Name() string           { return "merge_sort" }
func (m *MergeSort) DisplayName() string    { return "Merge Sort" }
func (m *MergeSort) TimeComplexity() string { return "O(n log n)" }
func (m *MergeSort) Description() string {
	return "Divide and Conquer algorithm that divides the array in half, recursively sorts, and merges. Guaranteed O(n log n) and stable."
}

func (m *MergeSort) Execute(arr []int) ([]engine.Step, int, int) {
	data := engine.CopyArray(arr)
	tracker := &mergeTracker{
		steps:       make([]engine.Step, 0, len(data)*len(data)),
		comparisons: 0,
		moves:       0,
	}

	mergeSortRecursive(data, 0, len(data)-1, tracker)

	// Mark all as sorted
	tracker.steps = append(tracker.steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   makeRange(0, len(data)-1),
		ActionType:   "sorted",
	})

	return tracker.steps, tracker.comparisons, tracker.moves
}

// mergeTracker accumulates steps and counters across recursive calls.
type mergeTracker struct {
	steps       []engine.Step
	comparisons int
	moves       int
}

func mergeSortRecursive(data []int, left, right int, t *mergeTracker) {
	if left >= right {
		return
	}

	mid := left + (right-left)/2

	// Show partition
	t.steps = append(t.steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{left, mid, right},
		ActionType:   "partition",
	})

	mergeSortRecursive(data, left, mid, t)
	mergeSortRecursive(data, mid+1, right, t)
	merge(data, left, mid, right, t)
}

func merge(data []int, left, mid, right int, t *mergeTracker) {
	// Create temporary arrays
	leftArr := make([]int, mid-left+1)
	rightArr := make([]int, right-mid)
	copy(leftArr, data[left:mid+1])
	copy(rightArr, data[mid+1:right+1])

	i, j, k := 0, 0, left

	for i < len(leftArr) && j < len(rightArr) {
		t.comparisons++

		// Show comparison
		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{left + i, mid + 1 + j},
			ActionType:   "compare",
		})

		if leftArr[i] <= rightArr[j] {
			data[k] = leftArr[i]
			i++
		} else {
			data[k] = rightArr[j]
			j++
		}
		t.moves++

		// Show merge placement
		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{k},
			ActionType:   "merge",
		})
		k++
	}

	// Copy remaining elements
	for i < len(leftArr) {
		data[k] = leftArr[i]
		t.moves++
		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{k},
			ActionType:   "merge",
		})
		i++
		k++
	}

	for j < len(rightArr) {
		data[k] = rightArr[j]
		t.moves++
		t.steps = append(t.steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{k},
			ActionType:   "merge",
		})
		j++
		k++
	}
}

// makeRange returns a slice of ints from start to end inclusive.
func makeRange(start, end int) []int {
	r := make([]int, end-start+1)
	for i := range r {
		r[i] = start + i
	}
	return r
}
