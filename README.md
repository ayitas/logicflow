# LogicFlow — Algorithm Visualizer

A high-performance, full-stack algorithm visualizer designed for backend engineers to analyze and compare algorithm efficiency in real-time through animated step-by-step execution traces.

Built with **Go** (Backend) and **D3.js v7** (Frontend). Zero external Go dependencies.

---

## About

LogicFlow visualizes how algorithms work internally — not just the final result, but every single comparison, swap, search check, and elimination rendered as an animated bar chart.

Reading pseudocode or Big-O notation alone doesn't give a clear picture of how algorithms actually behave under different data conditions. LogicFlow solves this by providing:

- **Step-by-step execution traces** generated server-side by the Go backend, giving an accurate representation of each algorithm's internal operations.
- **Real-time performance metrics** including comparison count, operation count, time complexity label, and actual server processing time in microseconds.
- **Visual comparison** across algorithms on the same dataset, making it straightforward to reason about algorithmic trade-offs.
- **Multi-category support** — currently supports sorting and searching algorithms, with architecture designed for future categories (graph, pathfinding, etc.).

---

## Supported Algorithms

### Sorting (6 algorithms)

| Algorithm | Time Complexity | Notes |
|-----------|----------------|-------|
| Bubble Sort | O(n²) | Stable, early termination |
| Selection Sort | O(n²) | Minimizes swaps |
| Insertion Sort | O(n²) | Efficient for nearly sorted data |
| Merge Sort | O(n log n) | Stable, divide & conquer |
| Quick Sort | O(n log n) avg | Lomuto partition |
| Heap Sort | O(n log n) | In-place, guaranteed performance |

### Searching (4 algorithms)

| Algorithm | Time Complexity | Notes |
|-----------|----------------|-------|
| Linear Search | O(n) | Works on unsorted arrays |
| Binary Search | O(log n) | Requires sorted array |
| Jump Search | O(√n) | Block-based, requires sorted array |
| Interpolation Search | O(log log n) avg | Best for uniform distributions |

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    Browser (Client)                   │
│  ┌────────────────────────────────────────────────┐   │
│  │  Vanilla JS + D3.js v7                        │   │
│  │  - SVG bar chart with animated transitions    │   │
│  │  - Category-aware UI (sorting vs searching)   │   │
│  │  - Real-time metric counters                  │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│         GET /algorithms    POST /execute              │
│         (auto-discovery)   (execution trace)          │
└────────────────────┼──────────────────────────────────┘
                     │
┌────────────────────┼──────────────────────────────────┐
│                Go HTTP Server (:8080)                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  handler/sort.go                               │   │
│  │  - POST /execute  → run algorithm, return trace│   │
│  │  - GET /algorithms → list registered algos     │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  engine/registry.go (Plugin Registry)          │   │
│  │  - Algorithm interface (with Category)         │   │
│  │  - Register() / Get() / List()                 │   │
│  └─────────────────┬──────────────────────────────┘   │
│                    │                                  │
│  ┌─────────────────┴──────────────────────────────┐   │
│  │  algorithm/ (one file per algorithm)           │   │
│  │  - Sorting: bubble, selection, insertion,      │   │
│  │             merge, quick, heap                 │   │
│  │  - Searching: linear, binary, jump,            │   │
│  │               interpolation                    │   │
│  │  - [future categories...]                      │   │
│  └────────────────────────────────────────────────┘   │
└───────────────────────────────────────────────────────┘
```

---

## Project Structure

```
logicflow/
├── main.go                            # Entry point — HTTP server, routing, static file serving
├── go.mod                             # Go module definition
│
├── internal/                          # Internal packages (unexported outside module)
│   ├── engine/
│   │   ├── types.go                   # Shared types: Step, AlgorithmRequest, AlgorithmResponse, ExecuteParams
│   │   └── registry.go               # Algorithm interface + global plugin registry
│   │
│   ├── algorithm/                     # Algorithm implementations (one file per algorithm)
│   │   ├── bubble.go                  # Bubble Sort — O(n²)
│   │   ├── selection.go               # Selection Sort — O(n²)
│   │   ├── insertion.go               # Insertion Sort — O(n²)
│   │   ├── merge.go                   # Merge Sort — O(n log n)
│   │   ├── quick.go                   # Quick Sort — O(n log n)
│   │   ├── heap.go                    # Heap Sort — O(n log n)
│   │   ├── linear_search.go           # Linear Search — O(n)
│   │   ├── binary_search.go           # Binary Search — O(log n)
│   │   ├── jump_search.go             # Jump Search — O(√n)
│   │   └── interpolation_search.go    # Interpolation Search — O(log log n)
│   │
│   └── handler/
│       └── sort.go                    # HTTP handlers for POST /execute and GET /algorithms
│
└── static/                            # Frontend (served as static files)
    ├── index.html                     # HTML5 — SVG container, controls, metrics panel
    ├── style.css                      # Dark theme, glassmorphism, responsive (mobile-first)
    └── script.js                      # D3.js v7 animation engine, category-aware UI
```

---

## Tech Stack

### Backend

| Component | Detail |
|-----------|--------|
| Language | Go 1.24 |
| HTTP Server | `net/http` (standard library, zero external dependencies) |
| Architecture | Plugin/Registry pattern with category support |
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
   Registered algorithms: 10

   [searching]
   - Binary Search             O(log n)
   - Interpolation Search      O(log log n)
   - Jump Search               O(√n)
   - Linear Search             O(n)
   [sorting]
   - Bubble Sort               O(n²)
   - Heap Sort                 O(n log n)
   - Insertion Sort            O(n²)
   - Merge Sort                O(n log n)
   - Quick Sort                O(n log n)
   - Selection Sort            O(n²)

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

### Sorting Algorithms

1. Select a sorting algorithm from the dropdown (grouped under "Sorting").
2. Adjust array size and animation speed using the sliders.
3. Click "Generate" for a new random array.
4. Click "Start" to execute — the array is sent to the Go backend, which returns a step-by-step trace. The frontend animates through each operation.

### Searching Algorithms

1. Select a searching algorithm from the dropdown (grouped under "Searching").
2. A "Search Target" input field appears — enter the value to find.
3. Click "Start" to execute — the backend runs the search and returns a trace showing every comparison and elimination.
4. When the animation finishes, a result banner shows whether the target was found (and at which index) or not found.

Note: Binary Search, Jump Search, and Interpolation Search require sorted arrays. The backend automatically sorts the array before searching.

### Bar Colors

**Sorting mode:**

| Color | Meaning |
|-------|---------|
| Gray/Purple | Default (unprocessed) |
| Yellow | Currently being compared |
| Red | Being swapped |
| Cyan | Pivot / Partition point |
| Green | Sorted (final position) |

**Searching mode:**

| Color | Meaning |
|-------|---------|
| Gray/Purple | Default |
| Yellow | Currently being checked |
| Cyan | Jump / Search range |
| Dimmed gray | Eliminated from search |
| Green | Found (target located) |

---

## API Reference

### `GET /algorithms`

Returns a list of all algorithms registered in the backend registry, grouped by category.

**Response:**

```json
[
  {
    "name": "binary_search",
    "display_name": "Binary Search",
    "category": "searching",
    "time_complexity": "O(log n)",
    "description": "Divides the sorted array in half repeatedly..."
  },
  {
    "name": "bubble_sort",
    "display_name": "Bubble Sort",
    "category": "sorting",
    "time_complexity": "O(n²)",
    "description": "Repeatedly steps through the list..."
  }
]
```

### `POST /execute`

Executes the specified algorithm and returns a step-by-step execution trace.

**Request (sorting):**

```json
{
  "algorithm": "bubble_sort",
  "array": [5, 3, 1, 4, 2]
}
```

**Request (searching):**

```json
{
  "algorithm": "binary_search",
  "array": [5, 3, 1, 4, 2],
  "target": 3
}
```

**Response:**

```json
{
  "steps": [
    {
      "current_state": [1, 2, 3, 4, 5],
      "highlights": [2],
      "action_type": "found"
    }
  ],
  "metadata": {
    "execution_time_us": 2,
    "comparisons": 3,
    "operations": 3,
    "time_complexity": "O(log n)",
    "algorithm_name": "Binary Search",
    "category": "searching",
    "found_index": 2
  }
}
```

---

## Adding a New Algorithm

The backend uses a **plugin/registry pattern**. Adding a new algorithm requires creating a single file — no modifications to any existing code.

### Steps

**1.** Create a new file in `internal/algorithm/`:

```go
package algorithm

import "logicflow/internal/engine"

type HeapSort struct{}

func init() {
    engine.Register(&HeapSort{})
}

func (h *HeapSort) Name() string           { return "heap_sort" }
func (h *HeapSort) DisplayName() string    { return "Heap Sort" }
func (h *HeapSort) Category() string       { return "sorting" }       // or "searching"
func (h *HeapSort) TimeComplexity() string { return "O(n log n)" }
func (h *HeapSort) Description() string    { return "Uses a binary heap..." }

func (h *HeapSort) Execute(params engine.ExecuteParams) ([]engine.Step, int, int) {
    data := engine.CopyArray(params.Array)
    target := params.Target  // only used by searching algorithms
    steps := make([]engine.Step, 0)
    comparisons, operations := 0, 0

    // Implement the algorithm, appending Steps for each operation...

    return steps, comparisons, operations
}
```

**2.** Restart the server. The new algorithm automatically appears in the frontend dropdown under the correct category.

### How It Works

- Go's `init()` functions execute automatically when a package is imported.
- `main.go` imports `_ "logicflow/internal/algorithm"`, which triggers all `init()` registrations.
- The frontend fetches `GET /algorithms` on page load to discover all registered algorithms dynamically.
- The UI automatically adapts based on the algorithm's `category` — showing target input for searching, adjusting legend colors, etc.

---

## License

MIT