<script lang="ts">
    import { onMount } from 'svelte';
    import { getSchema } from '$lib/api';
    import type { Schema } from '$lib/types/table';
    import { ArrowRightCircle, Database, Shield, Key, Link2 } from 'lucide-svelte';

    let schema: Schema | null = null;
    let error: string | null = null;
    let draggingTable: string | null = null;
    let initialMousePosition = { x: 0, y: 0 };
    let initialTablePosition = { x: 0, y: 0 };
    let lines: Array<{start: any, end: any, label: string, startTable: string, endTable: string}> = [];
    let isGridVisible = true;
    let autoLayout = true;
    let tablePositions = {};
    let containerWidth = 0;
    let containerHeight = 0;

    // Calculate better positions for tables in a grid layout
    function calculateInitialPositions() {
        if (!schema) return;

        const schemaKeys = Object.keys(schema);
        const tableCount = schemaKeys.length;

        // Calculate grid dimensions
        const gridCols = Math.ceil(Math.sqrt(tableCount));
        const gridRows = Math.ceil(tableCount / gridCols);

        // Calculate cell size based on container size
        const cellWidth = Math.floor(containerWidth / gridCols);
        const cellHeight = Math.floor(containerHeight / gridRows);

        // Position tables in grid
        schemaKeys.forEach((tableName, index) => {
            const row = Math.floor(index / gridCols);
            const col = index % gridCols;

            tablePositions[tableName] = {
                x: col * cellWidth + 50, // Add some margin
                y: row * cellHeight + 50 // Add some margin
            };
        });
    }

    function getTablePosition(element: HTMLElement) {
        const rect = element.getBoundingClientRect();
        const container = document.getElementById('schema-container');
        const containerRect = container?.getBoundingClientRect();

        if (!containerRect) return null;

        return {
            left: rect.left - containerRect.left + container.scrollLeft,
            top: rect.top - containerRect.top + container.scrollTop,
            width: rect.width,
            height: rect.height
        };
    }

    function calculateLines() {
        if (!schema) return;

        const newLines: typeof lines = [];
        const tablesWithForeignKeys = new Set();

        // First pass: collect all tables with foreign keys
        Object.entries(schema).forEach(([tableName, tableInfo]) => {
            if (tableInfo.foreignKeys?.length) {
                tableInfo.foreignKeys.forEach(fk => {
                    tablesWithForeignKeys.add(tableName);
                    tablesWithForeignKeys.add(fk.Table);
                });
            }
        });

        Object.entries(schema).forEach(([tableName, tableInfo]) => {
            if (!tableInfo.foreignKeys?.length) return;

            const sourceElement = document.getElementById(`table-${tableName}`);
            if (!sourceElement) return;

            tableInfo.foreignKeys.forEach(fk => {
                const targetElement = document.getElementById(`table-${fk.Table}`);
                if (!targetElement) return;

                const sourcePos = getTablePosition(sourceElement);
                const targetPos = getTablePosition(targetElement);

                if (!sourcePos || !targetPos) return;

                // Calculate points for better line routing
                const start = {
                    x: sourcePos.left + (sourcePos.width / 2),
                    y: sourcePos.top + (sourcePos.height / 2)
                };

                const end = {
                    x: targetPos.left + (targetPos.width / 2),
                    y: targetPos.top + (targetPos.height / 2)
                };

                // Determine connection points on the table borders
                let startPoint, endPoint;

                // Horizontal positioning
                if (start.x < end.x) {
                    // Source is to the left of target
                    startPoint = { x: sourcePos.left + sourcePos.width, y: start.y };
                    endPoint = { x: targetPos.left, y: end.y };
                } else {
                    // Source is to the right of target
                    startPoint = { x: sourcePos.left, y: start.y };
                    endPoint = { x: targetPos.left + targetPos.width, y: end.y };
                }

                // Adjust for vertical alignment if needed
                if (Math.abs(start.y - end.y) > Math.abs(start.x - end.x)) {
                    if (start.y < end.y) {
                        // Source is above target
                        startPoint = { x: start.x, y: sourcePos.top + sourcePos.height };
                        endPoint = { x: end.x, y: targetPos.top };
                    } else {
                        // Source is below target
                        startPoint = { x: start.x, y: sourcePos.top };
                        endPoint = { x: end.x, y: targetPos.top + targetPos.height };
                    }
                }

                newLines.push({
                    start: startPoint,
                    end: endPoint,
                    label: `${fk.From} → ${fk.To}`,
                    startTable: tableName,
                    endTable: fk.Table
                });
            });
        });

        lines = newLines;
    }

    function startDrag(event: MouseEvent, tableId: string) {
        draggingTable = tableId;
        const tableElement = document.getElementById(tableId);
        if (tableElement) {
            initialMousePosition = {
                x: event.clientX,
                y: event.clientY
            };
            initialTablePosition = {
                x: parseInt(tableElement.style.left) || 0,
                y: parseInt(tableElement.style.top) || 0
            };

            // Save the position in our state
            const tableName = tableId.replace('table-', '');
            tablePositions[tableName] = {
                x: initialTablePosition.x,
                y: initialTablePosition.y
            };

            event.preventDefault();
        }
    }

    function onDragMove(event: MouseEvent) {
        if (draggingTable) {
            const tableElement = document.getElementById(draggingTable);
            if (tableElement) {
                const deltaX = event.clientX - initialMousePosition.x;
                const deltaY = event.clientY - initialMousePosition.y;

                const newLeft = initialTablePosition.x + deltaX;
                const newTop = initialTablePosition.y + deltaY;

                tableElement.style.left = `${newLeft}px`;
                tableElement.style.top = `${newTop}px`;

                // Update position in our state
                const tableName = draggingTable.replace('table-', '');
                tablePositions[tableName] = {
                    x: newLeft,
                    y: newTop
                };

                // Recalculate lines when table is moved
                requestAnimationFrame(calculateLines);
            }
        }
    }

    function stopDrag() {
        draggingTable = null;
        calculateLines();
    }

    function toggleGrid() {
        isGridVisible = !isGridVisible;
    }

    function applyAutoLayout() {
        calculateInitialPositions();

        // Apply positions to DOM elements
        Object.entries(tablePositions).forEach(([tableName, position]) => {
            const tableElement = document.getElementById(`table-${tableName}`);
            if (tableElement) {
                tableElement.style.left = `${position.x}px`;
                tableElement.style.top = `${position.y}px`;
            }
        });

        setTimeout(calculateLines, 100);
    }

    onMount(async () => {
        try {
            const response = await getSchema();
            if (response.schema) {
                schema = response.schema;

                // Get container dimensions
                const container = document.getElementById('schema-container');
                if (container) {
                    containerWidth = container.clientWidth;
                    containerHeight = container.clientHeight;

                    // Initialize positions
                    calculateInitialPositions();

                    // Apply initial positions after a delay
                    setTimeout(() => {
                        if (autoLayout) applyAutoLayout();
                        setTimeout(calculateLines, 200);
                    }, 100);

                    // Add scroll event listener and resize observer
                    container.addEventListener('scroll', calculateLines);

                    const resizeObserver = new ResizeObserver(() => {
                        containerWidth = container.clientWidth;
                        containerHeight = container.clientHeight;
                        calculateLines();
                    });
                    resizeObserver.observe(container);
                }
            } else {
                error = "Schema data is missing.";
            }
        } catch (err) {
            error = `Error: ${err.message}`;
        }
    });
</script>

<div class="flex flex-col h-screen">
    <div class="bg-white shadow-sm p-4 flex justify-between items-center">
        <div class="flex items-center">
            <Database class="text-blue-600 mr-2" size={24} />
            <h1 class="text-xl font-semibold text-gray-800">Database Schema Visualizer</h1>
        </div>
        <div class="flex space-x-4">
            <button
                    class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition flex items-center"
                    on:click={applyAutoLayout}
            >
                <ArrowRightCircle size={18} class="mr-2" />
                Auto Layout
            </button>
            <button
                    class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition"
                    on:click={toggleGrid}
            >
                {isGridVisible ? 'Hide Grid' : 'Show Grid'}
            </button>
        </div>
    </div>

    <div class="flex-grow relative bg-gray-100 p-4 min-h-[800px] overflow-auto"
         id="schema-container"
         on:mousemove={onDragMove}
         on:mouseup={stopDrag}
         on:mouseleave={stopDrag}
         class:grid-bg={isGridVisible}
    >
        {#if error}
            <div class="text-red-500 bg-red-50 p-4 rounded-lg shadow">{error}</div>
        {:else if schema}
            <svg class="connections" width="100%" height="100%">
                {#each lines as line}
                    <g>
                        <path
                                d="M {line.start.x} {line.start.y}
                               C {line.start.x + ((line.end.x - line.start.x) / 3)} {line.start.y},
                                 {line.end.x - ((line.end.x - line.start.x) / 3)} {line.end.y},
                                 {line.end.x} {line.end.y}"
                                stroke="#4f46e5"
                                stroke-width="2"
                                fill="none"
                                marker-end="url(#arrowhead)"
                                class="relation-path"
                        />

                        <path
                                d="M {line.start.x} {line.start.y}
                               C {line.start.x + ((line.end.x - line.start.x) / 3)} {line.start.y},
                                 {line.end.x - ((line.end.x - line.start.x) / 3)} {line.end.y},
                                 {line.end.x} {line.end.y}"
                                stroke="transparent"
                                stroke-width="12"
                                fill="none"
                                class="relation-hover"
                        />

                        <g class="relation-label">
                            <rect
                                    x={(line.start.x + line.end.x) / 2 - 50}
                                    y={(line.start.y + line.end.y) / 2 - 12}
                                    width="100"
                                    height="24"
                                    rx="4"
                                    fill="white"
                                    stroke="#e5e7eb"
                            />
                            <text
                                    x={(line.start.x + line.end.x) / 2}
                                    y={(line.start.y + line.end.y) / 2 + 5}
                                    text-anchor="middle"
                                    class="relation-text"
                            >
                                {line.label}
                            </text>
                        </g>
                    </g>
                {/each}
                <defs>
                    <marker
                            id="arrowhead"
                            markerWidth="10"
                            markerHeight="7"
                            refX="9"
                            refY="3.5"
                            orient="auto"
                    >
                        <polygon points="0 0, 10 3.5, 0 7" fill="#4f46e5" />
                    </marker>
                </defs>
            </svg>

            {#each Object.entries(schema) as [tableName, tableInfo]}
                <div class="table-item rounded-lg shadow-lg bg-white border border-gray-200"
                     id="table-{tableName}"
                     style="top: {tablePositions[tableName]?.y || Object.keys(schema).indexOf(tableName) * 220}px;
                            left: {tablePositions[tableName]?.x || Object.keys(schema).indexOf(tableName) * 220}px;"
                     on:mousedown={(e) => startDrag(e, `table-${tableName}`)}>
                    <div class="table-header bg-indigo-600 text-white p-3 rounded-t-lg flex items-center justify-between">
                        <h2 class="text-lg font-semibold">{tableName}</h2>
                        <div class="flex items-center">
                            <span class="text-xs bg-indigo-800 px-2 py-1 rounded-full">
                                {tableInfo.columns?.length || 0} columns
                            </span>
                        </div>
                    </div>
                    <div class="p-4">
                        <div class="mb-3 flex items-center">
                            <Shield size={16} class="text-yellow-500 mr-2" />
                            <span class="text-sm font-medium text-gray-700">Primary Key:</span>
                            <span class="ml-2 px-2 py-1 bg-yellow-100 text-yellow-800 text-xs rounded-full">
                                {tableInfo.primaryKey}
                            </span>
                        </div>

                        <div class="mb-3">
                            <div class="font-medium text-gray-700 mb-1">Columns:</div>
                            <div class="grid grid-cols-1 gap-1 max-h-32 overflow-y-auto pr-1">
                                {#each tableInfo.columns as column}
                                    <div class="flex items-center text-sm py-1 px-2 bg-gray-50 rounded">
                                        <span>{column}</span>
                                    </div>
                                {/each}
                            </div>
                        </div>

                        {#if tableInfo.foreignKeys?.length}
                            <div>
                                <div class="font-medium text-gray-700 mb-1 flex items-center">
                                    <Link2 size={16} class="text-blue-500 mr-2" />
                                    Foreign Keys:
                                </div>
                                <div class="grid grid-cols-1 gap-1">
                                    {#each tableInfo.foreignKeys as fk}
                                        <div class="flex items-center text-sm py-1 px-2 bg-blue-50 rounded text-blue-800">
                                            <span class="font-medium">{fk.From}</span>
                                            <span class="mx-1">→</span>
                                            <span>{fk.Table}.{fk.To}</span>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        {:else}
            <div class="flex justify-center items-center h-full">
                <div class="text-center p-8 bg-white rounded-lg shadow">
                    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500 mx-auto mb-4"></div>
                    <p class="text-gray-600">Loading database schema...</p>
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .table-item {
        position: absolute;
        cursor: move;
        width: 300px;
        min-height: 200px;
        user-select: none;
        z-index: 1;
        transition: box-shadow 0.2s ease;
    }

    .table-item:hover {
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    }

    .table-header {
        cursor: move;
    }

    .connections {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        pointer-events: none;
        z-index: 0;
    }

    .relation-path {
        transition: stroke 0.2s ease;
    }

    .relation-hover:hover + .relation-path {
        stroke: #6366f1;
        stroke-width: 3;
    }

    .relation-hover {
        pointer-events: all;
        cursor: pointer;
    }

    .relation-label {
        pointer-events: none;
    }

    .relation-text {
        font-size: 12px;
        fill: #4b5563;
    }

    .grid-bg {
        background-image: linear-gradient(to right, rgba(0, 0, 0, 0.05) 1px, transparent 1px),
        linear-gradient(to bottom, rgba(0, 0, 0, 0.05) 1px, transparent 1px);
        background-size: 20px 20px;
    }

    :global(body) {
        margin: 0;
        padding: 0;
    }
</style>
