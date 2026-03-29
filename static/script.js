/**
 * LogicFlow — Algorithm Visualizer
 * Frontend logic: Fetch API, D3.js v7 animations, async playback
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
    animationId: null,  // for cancellation
    algorithms: [],
    speed: 50,          // ms per step
    arraySize: 30,
};

// ============================================================
// DOM Elements
// ============================================================
const dom = {
    select: document.getElementById('algorithm-select'),
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
    valSwaps: document.getElementById('val-swaps'),
    valComplexity: document.getElementById('val-complexity'),
    valTime: document.getElementById('val-time'),
    algoDescText: document.getElementById('algo-desc-text'),
    metricsPanel: document.getElementById('metrics-panel'),
};


// ============================================================
// Color Mapping for Action Types
// ============================================================
const actionColors = {
    compare:   getComputedStyle(document.documentElement).getPropertyValue('--color-compare').trim(),
    swap:      getComputedStyle(document.documentElement).getPropertyValue('--color-swap').trim(),
    partition: getComputedStyle(document.documentElement).getPropertyValue('--color-partition').trim(),
    merge:     getComputedStyle(document.documentElement).getPropertyValue('--color-merge').trim(),
    sorted:    getComputedStyle(document.documentElement).getPropertyValue('--color-sorted').trim(),
    insert:    getComputedStyle(document.documentElement).getPropertyValue('--color-insert').trim(),
    shift:     getComputedStyle(document.documentElement).getPropertyValue('--color-shift').trim(),
    default:   getComputedStyle(document.documentElement).getPropertyValue('--color-default').trim(),
};


// ============================================================
// Initialize
// ============================================================
async function init() {
    await loadAlgorithms();
    bindEvents();
    generateArray();
}

/**
 * Fetch available algorithms from the backend registry.
 * The dropdown is dynamically populated — no hardcoded values.
 */
async function loadAlgorithms() {
    try {
        const res = await fetch('/algorithms');
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        state.algorithms = await res.json();

        // Populate dropdown
        dom.select.innerHTML = '';
        state.algorithms.forEach((algo, i) => {
            const opt = document.createElement('option');
            opt.value = algo.name;
            opt.textContent = `${algo.display_name}  (${algo.time_complexity})`;
            dom.select.appendChild(opt);
        });

        // Select first and show description
        if (state.algorithms.length > 0) {
            dom.select.selectedIndex = 0;
            updateDescription();
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

    dom.select.addEventListener('change', updateDescription);

    // Responsive SVG resize
    window.addEventListener('resize', () => {
        if (!state.isRunning) renderBars(state.array, [], 'default');
    });
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
    resetMetrics();
    renderBars(state.array, [], 'default');
    dom.stepCurrent.textContent = '0';
    dom.stepTotal.textContent = '0';
}


// ============================================================
// Algorithm Description
// ============================================================
function updateDescription() {
    const selected = dom.select.value;
    const algo = state.algorithms.find(a => a.name === selected);
    if (algo) {
        dom.algoDescText.textContent = algo.description;
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
            .attr('y', height)  // start from bottom for entry animation
            .attr('width', barWidth)
            .attr('height', 0)
            .attr('fill', (d, i) => getBarColor(i, highlightSet, actionType))
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
 * Determine bar color based on whether it's highlighted and the action type.
 */
function getBarColor(index, highlightSet, actionType) {
    if (highlightSet.has(index)) {
        return actionColors[actionType] || actionColors.compare;
    }
    return actionColors.default;
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

    // Disable controls during running
    setRunningState(true);
    resetMetrics();

    try {
        // Fetch execution trace from Go backend
        const res = await fetch('/sort', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                algorithm: algorithm,
                array: [...state.array],
            }),
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
        await animateSteps(data.metadata.comparisons, data.metadata.swaps_moves);

    } catch (err) {
        console.error('Sort error:', err);
        alert(`Error: ${err.message}`);
    } finally {
        setRunningState(false);
    }
}

/**
 * Iterate through execution trace steps with configurable delay.
 * Uses a cancellation token pattern for the stop button.
 */
async function animateSteps(totalComparisons, totalSwaps) {
    const runId = Symbol('run');
    state.animationId = runId;

    let compCount = 0;
    let swapCount = 0;

    for (let i = 0; i < state.steps.length; i++) {
        // Check if animation was cancelled
        if (state.animationId !== runId) return;

        const step = state.steps[i];
        state.currentStep = i + 1;

        // Update counters based on action type
        if (step.action_type === 'compare') compCount++;
        if (['swap', 'merge', 'insert', 'shift'].includes(step.action_type)) swapCount++;

        // Update DOM
        dom.stepCurrent.textContent = state.currentStep;
        dom.valComparisons.textContent = compCount;
        dom.valSwaps.textContent = swapCount;

        // Render current state
        renderBars(step.current_state, step.highlights, step.action_type);

        // Wait before next step
        await sleep(state.speed);
    }

    // Final state — use server-side counts for accuracy
    dom.valComparisons.textContent = totalComparisons;
    dom.valSwaps.textContent = totalSwaps;
}

function stopAnimation() {
    state.animationId = null; // cancels the running loop
    setRunningState(false);
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
    dom.valSwaps.textContent = '0';
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
