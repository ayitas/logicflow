/**
 * LogicFlow — Algorithm Visualizer
 * Frontend logic: Fetch API, D3.js v7 animations, async playback
 * Supports both sorting and searching algorithm categories.
 *
 * D3.js is loaded globally via CDN in index.html
 */

// ============================================================
// State
// ============================================================
const state = {
    array: [],
    steps: [],
    currentStep: 0,
    isRunning: false,
    animationId: null,      // for cancellation
    algorithms: [],
    speed: 50,              // ms per step
    arraySize: 30,
    currentCategory: 'sorting',
    eliminatedIndices: new Set(),  // track eliminated bars during search
};

// ============================================================
// DOM Elements
// ============================================================
const dom = {
    select: document.getElementById('algorithm-select'),
    targetGroup: document.getElementById('target-group'),
    targetInput: document.getElementById('target-input'),
    sizeSlider: document.getElementById('array-size-slider'),
    sizeValue: document.getElementById('array-size-value'),
    speedSlider: document.getElementById('speed-slider'),
    speedValue: document.getElementById('speed-value'),
    btnGenerate: document.getElementById('btn-generate'),
    btnStart: document.getElementById('btn-start'),
    btnStop: document.getElementById('btn-stop'),
    svg: d3.select('#viz-svg'),
    stepCurrent: document.getElementById('step-current'),
    stepTotal: document.getElementById('step-total'),
    valComparisons: document.getElementById('val-comparisons'),
    valOperations: document.getElementById('val-operations'),
    operationsLabel: document.getElementById('operations-label'),
    valComplexity: document.getElementById('val-complexity'),
    valTime: document.getElementById('val-time'),
    algoDescText: document.getElementById('algo-desc-text'),
    algoCategoryBadge: document.getElementById('algo-category-badge'),
    metricsPanel: document.getElementById('metrics-panel'),
    legendItems: document.getElementById('legend-items'),
    searchResult: document.getElementById('search-result'),
    searchResultIcon: document.getElementById('search-result-icon'),
    searchResultText: document.getElementById('search-result-text'),
};


// ============================================================
// Color Mapping for Action Types
// ============================================================
function getActionColors() {
    const cs = getComputedStyle(document.documentElement);
    return {
        compare:   cs.getPropertyValue('--color-compare').trim(),
        swap:      cs.getPropertyValue('--color-swap').trim(),
        partition: cs.getPropertyValue('--color-partition').trim(),
        merge:     cs.getPropertyValue('--color-merge').trim(),
        sorted:    cs.getPropertyValue('--color-sorted').trim(),
        insert:    cs.getPropertyValue('--color-insert').trim(),
        shift:     cs.getPropertyValue('--color-shift').trim(),
        check:     cs.getPropertyValue('--color-check').trim(),
        found:     cs.getPropertyValue('--color-found').trim(),
        not_found: cs.getPropertyValue('--color-not-found').trim(),
        eliminate: cs.getPropertyValue('--color-eliminate').trim(),
        jump:      cs.getPropertyValue('--color-jump').trim(),
        default:   cs.getPropertyValue('--color-default').trim(),
    };
}

let actionColors = {};


// ============================================================
// Legend Definitions
// ============================================================
const legendConfig = {
    sorting: [
        { color: '--color-default', label: 'Default' },
        { color: '--color-compare', label: 'Comparing' },
        { color: '--color-swap', label: 'Swap / Move' },
        { color: '--color-partition', label: 'Partition / Pivot' },
        { color: '--color-sorted', label: 'Sorted' },
    ],
    searching: [
        { color: '--color-default', label: 'Default' },
        { color: '--color-check', label: 'Checking' },
        { color: '--color-jump', label: 'Jump / Range' },
        { color: '--color-eliminate', label: 'Eliminated' },
        { color: '--color-found', label: 'Found' },
    ],
};


// ============================================================
// Initialize
// ============================================================
async function init() {
    actionColors = getActionColors();
    await loadAlgorithms();
    bindEvents();
    generateArray();
    updateLegend();
}

/**
 * Fetch available algorithms from the backend registry.
 * Groups them by category using <optgroup>.
 */
async function loadAlgorithms() {
    try {
        const res = await fetch('/algorithms');
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        state.algorithms = await res.json();

        // Group by category
        const grouped = {};
        state.algorithms.forEach(algo => {
            if (!grouped[algo.category]) grouped[algo.category] = [];
            grouped[algo.category].push(algo);
        });

        // Populate dropdown with optgroups
        dom.select.innerHTML = '';
        const categoryOrder = ['searching', 'sorting'];
        categoryOrder.forEach(cat => {
            const algos = grouped[cat];
            if (!algos) return;
            const optgroup = document.createElement('optgroup');
            optgroup.label = cat.charAt(0).toUpperCase() + cat.slice(1);
            algos.forEach(algo => {
                const opt = document.createElement('option');
                opt.value = algo.name;
                opt.dataset.category = algo.category;
                opt.textContent = `${algo.display_name}  (${algo.time_complexity})`;
                optgroup.appendChild(opt);
            });
            dom.select.appendChild(optgroup);
        });

        // Select first and show description
        if (state.algorithms.length > 0) {
            dom.select.selectedIndex = 0;
            updateAlgorithmUI();
        }
    } catch (err) {
        console.error('Failed to load algorithms:', err);
        dom.select.innerHTML = '<option disabled>Error loading algorithms</option>';
    }
}


// ============================================================
// Event Binding
// ============================================================
function bindEvents() {
    dom.btnGenerate.addEventListener('click', () => {
        if (state.isRunning) stopAnimation();
        generateArray();
    });

    dom.btnStart.addEventListener('click', startVisualization);
    dom.btnStop.addEventListener('click', stopAnimation);

    dom.sizeSlider.addEventListener('input', (e) => {
        state.arraySize = parseInt(e.target.value);
        dom.sizeValue.textContent = state.arraySize;
    });

    dom.sizeSlider.addEventListener('change', () => {
        if (!state.isRunning) generateArray();
    });

    dom.speedSlider.addEventListener('input', (e) => {
        state.speed = parseInt(e.target.value);
        dom.speedValue.textContent = `${state.speed}ms`;
    });

    dom.select.addEventListener('change', () => {
        updateAlgorithmUI();
        if (!state.isRunning) {
            generateArray();
        }
    });

    // Responsive SVG resize
    window.addEventListener('resize', () => {
        if (!state.isRunning) renderBars(state.array, [], 'default');
    });
}


// ============================================================
// Algorithm UI Updates
// ============================================================
function updateAlgorithmUI() {
    const selected = dom.select.value;
    const algo = state.algorithms.find(a => a.name === selected);
    if (!algo) return;

    state.currentCategory = algo.category;

    // Update description
    dom.algoDescText.textContent = algo.description;

    // Update category badge
    dom.algoCategoryBadge.textContent = algo.category;
    dom.algoCategoryBadge.className = `category-badge ${algo.category}`;

    // Show/hide target input
    dom.targetGroup.style.display = algo.category === 'searching' ? '' : 'none';

    // Update operations label
    dom.operationsLabel.textContent = algo.category === 'searching' ? 'Checks' : 'Swaps / Moves';

    // Update legend
    updateLegend();

    // Hide search result
    dom.searchResult.style.display = 'none';
}

function updateLegend() {
    const items = legendConfig[state.currentCategory] || legendConfig.sorting;
    dom.legendItems.innerHTML = items.map(item =>
        `<div class="legend-item"><span class="legend-color" style="background: var(${item.color})"></span> ${item.label}</div>`
    ).join('');
}


// ============================================================
// Array Generation
// ============================================================
function generateArray() {
    const size = state.arraySize;
    state.array = Array.from({ length: size }, () =>
        Math.floor(Math.random() * 95) + 5  // values between 5 and 99
    );
    state.steps = [];
    state.currentStep = 0;
    state.eliminatedIndices = new Set();
    resetMetrics();
    renderBars(state.array, [], 'default');
    dom.stepCurrent.textContent = '0';
    dom.stepTotal.textContent = '0';
    dom.searchResult.style.display = 'none';

    // For searching: pre-fill a random target from the array (50% chance it exists)
    if (state.currentCategory === 'searching') {
        if (Math.random() > 0.3) {
            // Pick an existing value
            dom.targetInput.value = state.array[Math.floor(Math.random() * state.array.length)];
        } else {
            // Pick a value that might not exist
            dom.targetInput.value = Math.floor(Math.random() * 95) + 5;
        }
    }
}


// ============================================================
// D3.js Visualization
// ============================================================
function renderBars(array, highlights, actionType) {
    const svg = dom.svg;
    const container = document.getElementById('svg-container');
    const width = container.clientWidth;
    const height = parseInt(getComputedStyle(document.getElementById('viz-svg')).height);
    const padding = { top: 20, bottom: 10, left: 4, right: 4 };

    const chartWidth = width - padding.left - padding.right;
    const chartHeight = height - padding.top - padding.bottom;

    svg.attr('viewBox', `0 0 ${width} ${height}`)
       .attr('preserveAspectRatio', 'xMidYMid meet');

    const n = array.length;
    const gap = Math.max(1, Math.min(3, Math.floor(chartWidth / n * 0.1)));
    const barWidth = Math.max(2, (chartWidth - gap * (n - 1)) / n);

    const maxVal = Math.max(...array, 1);
    const yScale = d3.scaleLinear()
        .domain([0, maxVal])
        .range([0, chartHeight]);

    const highlightSet = new Set(highlights);

    // Data join with key function for stable identity
    const bars = svg.selectAll('rect.bar')
        .data(array, (d, i) => i);

    // EXIT
    bars.exit().remove();

    // ENTER + UPDATE
    bars.join(
        enter => enter.append('rect')
            .attr('class', 'bar')
            .attr('rx', Math.min(3, barWidth / 3))
            .attr('ry', Math.min(3, barWidth / 3))
            .attr('x', (d, i) => padding.left + i * (barWidth + gap))
            .attr('y', height)
            .attr('width', barWidth)
            .attr('height', 0)
            .attr('fill', (d, i) => getBarColor(i, highlightSet, actionType))
            .attr('opacity', (d, i) => getBarOpacity(i, highlightSet, actionType))
            .call(enter => enter.transition()
                .duration(300)
                .attr('y', (d) => padding.top + chartHeight - yScale(d))
                .attr('height', (d) => yScale(d))
            ),
        update => update
            .transition()
            .duration(state.isRunning ? Math.max(20, state.speed * 0.4) : 300)
            .attr('x', (d, i) => padding.left + i * (barWidth + gap))
            .attr('y', (d) => padding.top + chartHeight - yScale(d))
            .attr('width', barWidth)
            .attr('height', (d) => yScale(d))
            .attr('fill', (d, i) => getBarColor(i, highlightSet, actionType))
            .attr('opacity', (d, i) => getBarOpacity(i, highlightSet, actionType))
            .attr('rx', Math.min(3, barWidth / 3))
            .attr('ry', Math.min(3, barWidth / 3))
    );

    // Value labels (only show for smaller arrays)
    if (n <= 40 && barWidth > 14) {
        const labels = svg.selectAll('text.bar-label')
            .data(array, (d, i) => i);

        labels.exit().remove();

        labels.join(
            enter => enter.append('text')
                .attr('class', 'bar-label')
                .attr('text-anchor', 'middle')
                .attr('fill', '#e8e8f0')
                .attr('font-size', Math.min(11, barWidth * 0.5) + 'px')
                .attr('font-family', 'var(--font-mono)')
                .attr('font-weight', '500')
                .attr('x', (d, i) => padding.left + i * (barWidth + gap) + barWidth / 2)
                .attr('y', (d) => padding.top + chartHeight - yScale(d) - 5)
                .text(d => d),
            update => update
                .transition()
                .duration(state.isRunning ? Math.max(20, state.speed * 0.4) : 300)
                .attr('x', (d, i) => padding.left + i * (barWidth + gap) + barWidth / 2)
                .attr('y', (d) => padding.top + chartHeight - yScale(d) - 5)
                .text(d => d)
        );
    } else {
        svg.selectAll('text.bar-label').remove();
    }
}

/**
 * Determine bar color based on highlights, action type, and eliminated state.
 */
function getBarColor(index, highlightSet, actionType) {
    if (highlightSet.has(index)) {
        return actionColors[actionType] || actionColors.compare;
    }
    // For searching: show eliminated bars in gray
    if (state.currentCategory === 'searching' && state.eliminatedIndices.has(index)) {
        return actionColors.eliminate;
    }
    return actionColors.default;
}

/**
 * Determine bar opacity — dimmed for eliminated indices during search.
 */
function getBarOpacity(index, highlightSet, actionType) {
    if (highlightSet.has(index)) return 1;
    if (state.currentCategory === 'searching' && state.eliminatedIndices.has(index)) return 0.35;
    return 1;
}


// ============================================================
// Visualization Execution
// ============================================================
async function startVisualization() {
    const algorithm = dom.select.value;
    if (!algorithm) {
        alert('Please select an algorithm first.');
        return;
    }

    const algo = state.algorithms.find(a => a.name === algorithm);

    // Validate target for searching
    if (algo && algo.category === 'searching') {
        const targetVal = dom.targetInput.value;
        if (targetVal === '' || isNaN(parseInt(targetVal))) {
            alert('Please enter a target value to search for.');
            dom.targetInput.focus();
            return;
        }
    }

    // Disable controls during running
    setRunningState(true);
    resetMetrics();
    state.eliminatedIndices = new Set();
    dom.searchResult.style.display = 'none';

    try {
        // Build request body
        const body = {
            algorithm: algorithm,
            array: [...state.array],
        };

        // Add target for searching algorithms
        if (algo && algo.category === 'searching') {
            body.target = parseInt(dom.targetInput.value);
        }

        // Fetch execution trace from Go backend
        const res = await fetch('/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });

        if (!res.ok) {
            const err = await res.json().catch(() => ({ error: 'Unknown error' }));
            throw new Error(err.error || `HTTP ${res.status}`);
        }

        const data = await res.json();
        state.steps = data.steps;
        state.currentStep = 0;

        // Set final metrics immediately
        dom.valComplexity.textContent = data.metadata.time_complexity;
        dom.valTime.textContent = `${data.metadata.execution_time_us}µs`;
        dom.stepTotal.textContent = state.steps.length;

        // Animate through steps
        await animateSteps(data.metadata.comparisons, data.metadata.operations);

        // Show search result if applicable
        if (data.metadata.category === 'searching') {
            showSearchResult(data.metadata.found_index, body.target);
        }

    } catch (err) {
        console.error('Execution error:', err);
        alert(`Error: ${err.message}`);
    } finally {
        setRunningState(false);
    }
}

/**
 * Iterate through execution trace steps with configurable delay.
 * Uses a cancellation token pattern for the stop button.
 */
async function animateSteps(totalComparisons, totalOperations) {
    const runId = Symbol('run');
    state.animationId = runId;

    let compCount = 0;
    let opCount = 0;

    for (let i = 0; i < state.steps.length; i++) {
        // Check if animation was cancelled
        if (state.animationId !== runId) return;

        const step = state.steps[i];
        state.currentStep = i + 1;

        // Update counters based on action type
        if (step.action_type === 'compare' || step.action_type === 'check') compCount++;
        if (['swap', 'merge', 'insert', 'shift', 'check', 'jump'].includes(step.action_type)) opCount++;

        // Track eliminated indices for searching visualization
        if (step.action_type === 'eliminate' && state.currentCategory === 'searching') {
            step.highlights.forEach(idx => state.eliminatedIndices.add(idx));
        }

        // Update DOM
        dom.stepCurrent.textContent = state.currentStep;
        dom.valComparisons.textContent = compCount;
        dom.valOperations.textContent = opCount;

        // Render current state
        renderBars(step.current_state, step.highlights, step.action_type);

        // Wait before next step
        await sleep(state.speed);
    }

    // Final state — use server-side counts for accuracy
    dom.valComparisons.textContent = totalComparisons;
    dom.valOperations.textContent = totalOperations;
}

function stopAnimation() {
    state.animationId = null; // cancels the running loop
    setRunningState(false);
}

/**
 * Show search result banner after search completes.
 */
function showSearchResult(foundIndex, target) {
    dom.searchResult.style.display = '';
    if (foundIndex >= 0) {
        dom.searchResult.className = 'glass-card found';
        dom.searchResultIcon.textContent = '✓';
        dom.searchResultText.textContent = `Target ${target} found at index ${foundIndex}`;
    } else {
        dom.searchResult.className = 'glass-card not-found';
        dom.searchResultIcon.textContent = '✗';
        dom.searchResultText.textContent = `Target ${target} was not found in the array`;
    }
}


// ============================================================
// UI State Management
// ============================================================
function setRunningState(running) {
    state.isRunning = running;
    dom.btnStart.disabled = running;
    dom.btnGenerate.disabled = running;
    dom.btnStop.disabled = !running;
    dom.select.disabled = running;
    dom.sizeSlider.disabled = running;
    dom.targetInput.disabled = running;

    if (running) {
        dom.metricsPanel.classList.add('is-running');
        dom.btnStart.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none"><path d="M4 2.5v11l9-5.5z" fill="currentColor"/></svg>
            Running...`;
    } else {
        dom.metricsPanel.classList.remove('is-running');
        dom.btnStart.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none"><path d="M4 2.5v11l9-5.5z" fill="currentColor"/></svg>
            Start`;
    }
}

function resetMetrics() {
    dom.valComparisons.textContent = '0';
    dom.valOperations.textContent = '0';
    dom.valComplexity.textContent = '—';
    dom.valTime.textContent = '—';
    dom.stepCurrent.textContent = '0';
    dom.stepTotal.textContent = '0';
}


// ============================================================
// Utilities
// ============================================================
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}


// ============================================================
// Boot
// ============================================================
document.addEventListener('DOMContentLoaded', init);
