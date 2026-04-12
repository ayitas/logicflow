package algorithm

import (
	"math"
	"sort"

	"logicflow/internal/engine"
)

// JumpSearch implements the Jump Search algorithm.
// Time: O(√n) | Space: O(1)
// Requires a sorted array. Jumps ahead by fixed steps, then does linear search.
type JumpSearch struct{}

func init() {
	engine.Register(&JumpSearch{})
}

func (j *JumpSearch) Name() string           { return "jump_search" }
func (j *JumpSearch) DisplayName() string    { return "Jump Search" }
func (j *JumpSearch) Category() string       { return "searching" }
func (j *JumpSearch) TimeComplexity() string { return "O(√n)" }
func (j *JumpSearch) Description() string {
	return "Jumps ahead by a fixed block size (√n), then performs a linear search within the identified block. Requires a pre-sorted array. Balances between linear and binary search."
}

func (j *JumpSearch) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
	data := engine.CopyArray(params.Array)
	target := params.Target
	n := len(data)

	// Jump search requires sorted array
	sort.Ints(data)

	steps := make([]engine.Step, 0, n)
	comparisons := 0
	checks := 0

	// Show sorted array
	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   makeRange(0, n-1),
		ActionType:   "sorted",
	})

	blockSize := int(math.Sqrt(float64(n)))
	prev := 0

	// Phase 1: Jump forward in blocks
	for {
		curr := prev + blockSize - 1
		if curr >= n {
			curr = n - 1
		}

		comparisons++
		checks++

		// Show the jump
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{curr},
			ActionType:   "jump",
		})

		// Check if we've jumped past the target
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   []int{curr},
			ActionType:   "check",
		})

		if data[curr] >= target {
			// Target is in [prev, curr] block — do linear search
			break
		}

		// Eliminate the block we just jumped over
		eliminateHighlights := make([]int, 0)
		for i := prev; i <= curr; i++ {
			eliminateHighlights = append(eliminateHighlights, i)
		}
		steps = append(steps, engine.Step{
			CurrentState: engine.SnapshotArray(data),
			Highlights:   eliminateHighlights,
			ActionType:   "eliminate",
		})

		prev = curr + 1

		if prev >= n {
			steps = append(steps, engine.Step{
				CurrentState: engine.SnapshotArray(data),
				Highlights:   []int{},
				ActionType:   "not_found",
			})
			return steps, comparisons, checks
		}
	}

	// Phase 2: Linear search within the block [prev, prev+blockSize)
	end := prev + blockSize
	if end > n {
		end = n
	}

	// Highlight the block being searched
	blockHighlights := make([]int, 0)
	for i := prev; i < end; i++ {
		blockHighlights = append(blockHighlights, i)
	}
	steps = append(steps, engine.Step{
		CurrentState: engine.SnapshotArray(data),
		Highlights:   blockHighlights,
		ActionType:   "partition",
	})

	for i := prev; i < end; i++ {
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

		if data[i] > target {
			break
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
