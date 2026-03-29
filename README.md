# LogicFlow — Sorting Algorithm Visualizer

A high-performance, full-stack algorithm visualizer designed for backend engineers to analyze and compare sorting efficiency in real-time through animated step-by-step execution traces.

Built with **Go** (Backend) and **D3.js v7** (Frontend). Zero external Go dependencies.

---

## About

LogicFlow visualizes how sorting algorithms work internally — not just the final result, but every single comparison, swap, and merge operation rendered as an animated bar chart.

Reading pseudocode or Big-O notation alone doesn't give a clear picture of how algorithms actually behave under different data conditions. LogicFlow solves this by providing:

- **Step-by-step execution traces** generated server-side by the Go backend, giving an accurate representation of each algorithm's internal operations.
- **Real-time performance metrics** including comparison count, swap/move count, time complexity label, and actual server processing time in microseconds.
- **Visual comparison** across algorithms on the same dataset, making it straightforward to reason about algorithmic trade-offs.

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    Browser (Client)                   │
│  ┌────────────────────────────────────────────────┐   │
│  │  Vanilla JS + D3.js v7                        │   │
│  │  - SVG bar chart with animated transitions    │   │
│  │  - Async step-by-step playback                │   │
│  │  - Real-time metric counters                  │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│         GET /algorithms    POST /sort                 │
│         (auto-discovery)   (execution trace)          │
└────────────────────┼──────────────────────────────────┘
                     │
┌────────────────────┼──────────────────────────────────┐
│                Go HTTP Server (:8080)                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  handler/sort.go                               │   │
│  │  - POST /sort  → run algorithm, return trace   │   │
│  │  - GET /algorithms → list registered algos     │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  engine/registry.go (Plugin Registry)          │   │
│  │  - Algorithm interface                         │   │
│  │  - Register() / Get() / List()                 │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  algorithm/                                    │   │
│  │  - bubble.go    → Bubble Sort     O(n²)        │   │
│  │  - selection.go → Selection Sort  O(n²)        │   │
│  │  - insertion.go → Insertion Sort  O(n²)        │   │
│  │  - merge.go     → Merge Sort     O(n log n)    │   │
│  │  - quick.go     → Quick Sort     O(n log n)    │   │
│  │  - [future algorithms...]                      │   │
│  └────────────────────────────────────────────────┘   │
└───────────────────────────────────────────────────────┘
```

---

## Project Structure

```
logicflow/
├── main.go                        # Entry point — HTTP server, routing, static file serving
├── go.mod                         # Go module definition
│
├── internal/                      # Internal packages (unexported outside module)
│   ├── engine/
│   │   ├── types.go               # Shared types: Step, SortRequest, SortResponse, Metadata
│   │   └── registry.go            # Algorithm interface + global plugin registry
│   │
│   ├── algorithm/                 # Algorithm implementations (one file per algorithm)
│   │   ├── bubble.go              # Bubble Sort — O(n²), stable, early termination
│   │   ├── selection.go           # Selection Sort — O(n²), minimizes swaps
│   │   ├── insertion.go           # Insertion Sort — O(n²), efficient for nearly sorted data
│   │   ├── merge.go               # Merge Sort — O(n log n), stable, divide & conquer
│   │   └── quick.go               # Quick Sort — O(n log n) avg, Lomuto partition
│   │
│   └── handler/
│       └── sort.go                # HTTP handlers for POST /sort and GET /algorithms
│
└── static/                        # Frontend (served as static files)
    ├── index.html                 # HTML5 — SVG container, controls, metrics panel
    ├── style.css                  # Dark theme, glassmorphism, responsive (mobile-first)
    └── script.js                  # D3.js v7 animation engine, Fetch API, async playback
```

---

## Tech Stack

### Backend

| Component | Detail |
|-----------|--------|
| Language | Go 1.24 |
| HTTP Server | `net/http` (standard library, zero external dependencies) |
| Architecture | Plugin/Registry pattern — designed for easy algorithm extension |
| API Format | JSON REST |

### Frontend

| Component | Detail |
|-----------|--------|
| Language | Vanilla JavaScript (no frameworks) |
| Visualization | D3.js v7.9.0 — SVG-based bar chart with animated transitions |
| Styling | Vanilla CSS — dark theme, glassmorphism, CSS Grid/Flexbox |
| Typography | Google Fonts (Inter, JetBrains Mono) |
| Responsive | Mobile-first with breakpoints for smartphone, tablet, and desktop |

---

## Getting Started

### Prerequisites

- Go >= 1.24

### Run

```bash
git clone <repository-url>
cd logicflow
go run main.go
```

Expected output:

```
LogicFlow — Algorithm Visualizer
   Registered algorithms: 5
   - Bubble Sort          O(n²)
   - Selection Sort       O(n²)
   - Insertion Sort       O(n²)
   - Merge Sort           O(n log n)
   - Quick Sort           O(n log n)

   Server listening on http://localhost:8080
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

### Build Binary (Optional)

```bash
go build -o logicflow .
./logicflow
```

---

## Usage

1. **Select an algorithm** from the dropdown (all 5 sorting algorithms are listed with their time complexity).
2. **Adjust array size** using the slider (5–100 elements).
3. **Adjust speed** using the speed slider (1ms–200ms delay per step).
4. **Generate** a new random array by clicking "Generate".
5. **Start** the visualization — the array is sent to the Go backend, which returns a full execution trace. The frontend then animates through each step.
6. **Stop** the animation at any time.

### Bar Colors

| Color | Meaning |
|-------|---------|
| Gray/Purple | Default (unprocessed) |
| Yellow | Currently being compared |
| Red | Being swapped |
| Cyan | Pivot / Partition point |
| Green | Sorted (final position) |

### Displayed Metrics

- **Comparisons** — total comparisons performed
- **Swaps/Moves** — total element exchanges or movements
- **Time Complexity** — Big-O label for the selected algorithm
- **Server Time** — actual Go backend execution time in microseconds

---

## API Reference

### `GET /algorithms`

Returns a list of all algorithms currently registered in the backend registry. The frontend calls this on page load to dynamically populate the dropdown.

**Response:**

```json
[
  {
    "name": "bubble_sort",
    "display_name": "Bubble Sort",
    "time_complexity": "O(n²)",
    "description": "Repeatedly steps through the list, compares adjacent elements, and swaps them if they are in the wrong order."
  }
]
```

### `POST /sort`

Executes the specified sorting algorithm on the given array and returns a step-by-step execution trace.

**Request:**

```json
{
  "algorithm": "bubble_sort",
  "array": [5, 3, 1, 4, 2]
}
```

**Response:**

```json
{
  "steps": [
    {
      "current_state": [5, 3, 1, 4, 2],
      "highlights": [0, 1],
      "action_type": "compare"
    },
    {
      "current_state": [3, 5, 1, 4, 2],
      "highlights": [0, 1],
      "action_type": "swap"
    }
  ],
  "metadata": {
    "execution_time_us": 3,
    "comparisons": 10,
    "swaps_moves": 7,
    "time_complexity": "O(n²)",
    "algorithm_name": "Bubble Sort"
  }
}
```

---

## Adding a New Algorithm

The backend uses a **plugin/registry pattern**. Adding a new algorithm requires creating a single file — no modifications to any existing code.

### Steps

**1.** Create a new file in `internal/algorithm/`, for example `heap.go`:

```go
package algorithm

import "logicflow/internal/engine"

type HeapSort struct{}

func init() {
    engine.Register(&HeapSort{})
}

func (h *HeapSort) Name() string           { return "heap_sort" }
func (h *HeapSort) DisplayName() string    { return "Heap Sort" }
func (h *HeapSort) TimeComplexity() string { return "O(n log n)" }
func (h *HeapSort) Description() string    { return "Uses a binary heap data structure to sort." }

func (h *HeapSort) Execute(arr []int) ([]engine.Step, int, int) {
    data := engine.CopyArray(arr)
    steps := make([]engine.Step, 0)
    comparisons, swaps := 0, 0

    // Implement the algorithm here.
    // For each significant operation, append a Step:
    //   steps = append(steps, engine.Step{
    //       CurrentState: engine.SnapshotArray(data),
    //       Highlights:   []int{i, j},
    //       ActionType:   "compare",  // or "swap", "merge", etc.
    //   })

    return steps, comparisons, swaps
}
```

**2.** Restart the server. The new algorithm automatically appears in the frontend dropdown.

### How It Works

- Go's `init()` functions execute automatically when a package is imported.
- `main.go` imports `_ "logicflow/internal/algorithm"`, which triggers all `init()` registrations.
- The frontend fetches `GET /algorithms` on page load to discover all registered algorithms dynamically — no hardcoded values on the client side.

---

## License

MIT