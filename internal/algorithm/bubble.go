package algorithm

import (
	"logicflow/internal/engine"
)

// BubbleSort implements the classic Bubble Sort algorithm.
// Time: O(n²) | Space: O(1) | Stable: Yes
type BubbleSort struct{}

func init() {
	engine.Register(&BubbleSort{})
}

func (b *BubbleSort) Name() string           { return "bubble_sort" }
func (b *BubbleSort) DisplayName() string    { return "Bubble Sort" }
func (b *BubbleSort) TimeComplexity() string { return "O(n²)" }
func (b *BubbleSort) Description() string {
	return "Repeatedly steps through the list, compares adjacent elements, and swaps them if they are in the wrong order. Educational baseline algorithm."
}

func (b *BubbleSort) Execute(arr []int) ([]engine.Step, int, int) {
	data := engine.CopyArray(arr)
	n := len(data)
	steps := make([]engine.Step, 0, n*n)
	comparisons := 0
	swaps := 0

	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			// Compare adjacent elements
			comparisons++
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{j, j + 1},
				ActionType:   "compare",
			})

			if data[j] > data[j+1] {
				// Swap
				data[j], data[j+1] = data[j+1], data[j]
				swaps++
				swapped = true
				steps = append(steps, engine.Step{
					CurrentState: engine.SnapshotArray(data),
					Highlights:   []int{j, j + 1},
					ActionType:   "swap",
				})
			}
		}

		// Mark the last element of this pass as sorted
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{n - 1 - i},
			ActionType:   "sorted",
		})

		// Early termination if no swaps occurred
		if !swapped {
			// Mark remaining elements as sorted
			for k := n - 2 - i; k >= 0; k-- {
				steps = append(steps, engine.Step{
					CurrentState: engine.SnapshotArray(data),
					Highlights:   []int{k},
					ActionType:   "sorted",
				})
			}
			break
		}
	}

	return steps, comparisons, swaps
}
