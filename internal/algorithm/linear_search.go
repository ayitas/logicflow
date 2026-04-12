package algorithm

import (
	"logicflow/internal/engine"
)

// LinearSearch implements the Linear Search algorithm.
// Time: O(n) | Space: O(1)
// Searches sequentially through every element. Works on unsorted arrays.
type LinearSearch struct{}

func init() {
	engine.Register(&LinearSearch{})
}

func (l *LinearSearch) Name() string           { return "linear_search" }
func (l *LinearSearch) DisplayName() string    { return "Linear Search" }
func (l *LinearSearch) Category() string       { return "searching" }
func (l *LinearSearch) TimeComplexity() string { return "O(n)" }
func (l *LinearSearch) Description() string {
	return "Sequentially checks each element of the array until the target is found or the end is reached. Works on any array regardless of order."
}

func (l *LinearSearch) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	target := params.Target
	steps := make([]engine.Step, 0, len(data))
	comparisons := 0
	checks := 0

	for i := 0; i < len(data); i++ {
		comparisons++
		checks++

		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i},
			ActionType:   "check",
		})

		if data[i] == target {
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{i},
				ActionType:   "found",
			})
			return steps, comparisons, checks
		}

		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{i},
			ActionType:   "eliminate",
		})
	}

	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   []int{},
		ActionType:   "not_found",
	})

	return steps, comparisons, checks
}
