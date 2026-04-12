package algorithm

import (
	"logicflow/internal/engine"
)

// HeapSort implements the Heap Sort algorithm using a max-heap.
// Time: O(n log n) | Space: O(1) | Stable: No
type HeapSort struct{}

func init() {
	engine.Register(&HeapSort{})
}

func (h *HeapSort) Name() string           { return "heap_sort" }
func (h *HeapSort) DisplayName() string    { return "Heap Sort" }
func (h *HeapSort) Category() string       { return "sorting" }
func (h *HeapSort) TimeComplexity() string { return "O(n log n)" }
func (h *HeapSort) Description() string {
	return "Builds a max-heap from the array, then repeatedly extracts the maximum element and places it at the end. In-place with guaranteed O(n log n) time complexity."
}

func (h *HeapSort) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	n := len(data)
	steps := make([]engine.Step, 0, n*n)
	comparisons := 0
	swaps := 0

	for i := n/2 - 1; i >= 0; i-- {
		siftDown(data, n, i, &steps, &comparisons, &swaps)
	}

	for i := n - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0]
		swaps++
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{0, i},
			ActionType:   "swap",
		})

		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i},
			ActionType:   "sorted",
		})

		siftDown(data, i, 0, &steps, &comparisons, &swaps)
	}

	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{0},
		ActionType:   "sorted",
	})

	return steps, comparisons, swaps
}

func siftDown(data []int, heapSize, i int, steps *[]engine.Step, comparisons, swaps *int) {
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < heapSize {
			*comparisons++
			*steps = append(*steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{largest, left},
				ActionType:   "compare",
			})
			if data[left] > data[largest] {
				largest = left
			}
		}

		if right < heapSize {
			*comparisons++
			*steps = append(*steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{largest, right},
				ActionType:   "compare",
			})
			if data[right] > data[largest] {
				largest = right
			}
		}

		if largest == i {
			break
		}

		data[i], data[largest] = data[largest], data[i]
		*swaps++
		*steps = append(*steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i, largest},
			ActionType:   "swap",
		})

		i = largest
	}
}
