package algorithm

import (
	"logicflow/internal/engine"
)

// InsertionSort implements the Insertion Sort algorithm.
// Time: O(n²) | Space: O(1) | Stable: Yes
// Efficient for small or nearly sorted datasets.
type InsertionSort struct{}

func init() {
	engine.Register(&InsertionSort{})
}

func (is *InsertionSort) Name() string           { return "insertion_sort" }
func (is *InsertionSort) DisplayName() string    { return "Insertion Sort" }
func (is *InsertionSort) Category() string       { return "sorting" }
func (is *InsertionSort) TimeComplexity() string { return "O(n²)" }
func (is *InsertionSort) Description() string {
	return "Builds the sorted array one element at a time by inserting each element into its correct position. Efficient for small or nearly sorted data."
}

func (is *InsertionSort) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	n := len(data)
	steps := make([]engine.Step, 0, n*n)
	comparisons := 0
	moves := 0

	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{0},
		ActionType:   "sorted",
	})

	for i := 1; i < n; i++ {
		key := data[i]
		j := i - 1

		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i},
			ActionType:   "compare",
		})

		for j >= 0 {
			comparisons++
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{j, j + 1},
				ActionType:   "compare",
			})

			if data[j] > key {
				data[j+1] = data[j]
				moves++
				steps = append(steps, engine.Step{
					CurrentState: engine.SnapshotArray(data),
					Highlights:   []int{j, j + 1},
					ActionType:   "shift",
				})
				j--
			} else {
				break
			}
		}

		data[j+1] = key
		moves++
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{j + 1},
			ActionType:   "insert",
		})
	}

	return steps, comparisons, moves
}
