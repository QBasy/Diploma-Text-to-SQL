<script lang="ts">
    import { onMount } from 'svelte';
    import { getSchema } from '$lib/api';
    import type { Schema } from '$lib/types/table';
    import { ArrowRightCircle, Database, Shield, Key, Link2, ZoomIn, ZoomOut, Move, Grid3X3, Eye, EyeOff } from 'lucide-svelte';

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
    let isMobile = false;
    let zoom = 1;
    let minZoom = 0.3;
    let maxZoom = 2;
    let showRelations = true;
    let selectedTable: string | null = null;
    let viewMode: 'compact' | 'normal' | 'detailed' = 'compact';

    // Touch handling for mobile
    let touchStartPos = { x: 0, y: 0 };
    let touchStartTime = 0;
    let isTouch = false;

    // Check if device is mobile
    function checkMobile() {
        isMobile = window.innerWidth < 768;
        if (isMobile) {
            zoom = 0.5; // Even smaller zoom on mobile for many tables
            viewMode = 'compact';
        }
    }

    // Calculate better positions for tables in a grid layout
    function calculateInitialPositions() {
        if (!schema) return;

        const schemaKeys = Object.keys(schema);
        const tableCount = schemaKeys.length;

        // Adjust spacing based on table count and view mode
        let spacing;
        let startOffset;

        if (tableCount > 20) {
            // Many tables - compact layout
            spacing = isMobile ? 120 : 160;
            startOffset = isMobile ? 10 : 20;
        } else if (tableCount > 10) {
            // Medium number of tables
            spacing = isMobile ? 140 : 180;
            startOffset = isMobile ? 15 : 30;
        } else {
            // Few tables - spacious layout
            spacing = isMobile ? 160 : 220;
            startOffset = isMobile ? 20 : 50;
        }

        // Calculate grid dimensions - more columns for many tables
        let gridCols;
        if (tableCount > 25) {
            gridCols = isMobile ? 3 : 6;
        } else if (tableCount > 15) {
            gridCols = isMobile ? 3 : 5;
        } else if (tableCount > 8) {
            gridCols = isMobile ? 2 : 4;
        } else {
            gridCols = isMobile ? 2 : 3;
        }

        const gridRows = Math.ceil(tableCount / gridCols);

        // Calculate container size needed
        const minContainerWidth = gridCols * spacing + startOffset * 2;
        const minContainerHeight = gridRows * spacing + startOffset * 2;

        // Position tables in grid
        schemaKeys.forEach((tableName, index) => {
            const row = Math.floor(index / gridCols);
            const col = index % gridCols;

            tablePositions[tableName] = {
                x: col * spacing + startOffset,
                y: row * spacing + startOffset
            };
        });
    }

    function getTablePosition(element: HTMLElement) {
        const rect = element.getBoundingClientRect();
        const container = document.getElementById('schema-container');
        const containerRect = container?.getBoundingClientRect();

        if (!containerRect) return null;

        return {
            left: (rect.left - containerRect.left + container.scrollLeft) / zoom,
            top: (rect.top - containerRect.top + container.scrollTop) / zoom,
            width: rect.width / zoom,
            height: rect.height / zoom
        };
    }

    function calculateLines() {
        if (!schema || !showRelations) {
            lines = [];
            return;
        }

        const newLines: typeof lines = [];

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

                // Calculate connection points
                let startPoint, endPoint;

                const sourceCenterX = sourcePos.left + (sourcePos.width / 2);
                const sourceCenterY = sourcePos.top + (sourcePos.height / 2);
                const targetCenterX = targetPos.left + (targetPos.width / 2);
                const targetCenterY = targetPos.top + (targetPos.height / 2);

                // Determine best connection points
                if (Math.abs(sourceCenterX - targetCenterX) > Math.abs(sourceCenterY - targetCenterY)) {
                    // Connect horizontally
                    if (sourceCenterX < targetCenterX) {
                        startPoint = { x: sourcePos.left + sourcePos.width, y: sourceCenterY };
                        endPoint = { x: targetPos.left, y: targetCenterY };
                    } else {
                        startPoint = { x: sourcePos.left, y: sourceCenterY };
                        endPoint = { x: targetPos.left + targetPos.width, y: targetCenterY };
                    }
                } else {
                    // Connect vertically
                    if (sourceCenterY < targetCenterY) {
                        startPoint = { x: sourceCenterX, y: sourcePos.top + sourcePos.height };
                        endPoint = { x: targetCenterX, y: targetPos.top };
                    } else {
                        startPoint = { x: sourceCenterX, y: sourcePos.top };
                        endPoint = { x: targetCenterX, y: targetPos.top + targetPos.height };
                    }
                }

                newLines.push({
                    start: startPoint,
                    end: endPoint,
                    label: `${fk.from} → ${fk.to}`,
                    startTable: tableName,
                    endTable: fk.Table
                });
            });
        });

        lines = newLines;
    }

    function handleTouchStart(event: TouchEvent, tableId: string) {
        if (event.touches.length === 1) {
            isTouch = true;
            touchStartTime = Date.now();
            const touch = event.touches[0];
            touchStartPos = { x: touch.clientX, y: touch.clientY };
            startDrag({ clientX: touch.clientX, clientY: touch.clientY } as MouseEvent, tableId);
        }
    }

    function handleTouchMove(event: TouchEvent) {
        if (event.touches.length === 1 && isTouch) {
            event.preventDefault();
            const touch = event.touches[0];
            onDragMove({ clientX: touch.clientX, clientY: touch.clientY } as MouseEvent);
        }
    }

    function handleTouchEnd(event: TouchEvent) {
        const touchDuration = Date.now() - touchStartTime;
        const touchDistance = Math.sqrt(
            Math.pow(event.changedTouches[0].clientX - touchStartPos.x, 2) +
            Math.pow(event.changedTouches[0].clientY - touchStartPos.y, 2)
        );

        if (touchDuration < 300 && touchDistance < 10) {
            const tableId = (event.target as HTMLElement).closest('.table-item')?.id;
            if (tableId) {
                const tableName = tableId.replace('table-', '');
                selectTable(tableName);
            }
        }

        isTouch = false;
        stopDrag();
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
                const deltaX = (event.clientX - initialMousePosition.x) / zoom;
                const deltaY = (event.clientY - initialMousePosition.y) / zoom;

                const newLeft = Math.max(0, initialTablePosition.x + deltaX);
                const newTop = Math.max(0, initialTablePosition.y + deltaY);

                tableElement.style.left = `${newLeft}px`;
                tableElement.style.top = `${newTop}px`;

                const tableName = draggingTable.replace('table-', '');
                tablePositions[tableName] = { x: newLeft, y: newTop };

                requestAnimationFrame(calculateLines);
            }
        }
    }

    function stopDrag() {
        draggingTable = null;
        calculateLines();
    }

    function selectTable(tableName: string) {
        selectedTable = selectedTable === tableName ? null : tableName;
    }

    function toggleGrid() {
        isGridVisible = !isGridVisible;
    }

    function toggleRelations() {
        showRelations = !showRelations;
        calculateLines();
    }

    function zoomIn() {
        zoom = Math.min(maxZoom, zoom + 0.1);
        setTimeout(calculateLines, 50);
    }

    function zoomOut() {
        zoom = Math.max(minZoom, zoom - 0.1);
        setTimeout(calculateLines, 50);
    }

    function applyAutoLayout() {
        calculateInitialPositions();

        Object.entries(tablePositions).forEach(([tableName, position]) => {
            const tableElement = document.getElementById(`table-${tableName}`);
            if (tableElement) {
                tableElement.style.left = `${position.x}px`;
                tableElement.style.top = `${position.y}px`;
            }
        });

        setTimeout(calculateLines, 100);
    }

    function toggleViewMode() {
        const modes = ['compact', 'normal', 'detailed'];
        const currentIndex = modes.indexOf(viewMode);
        viewMode = modes[(currentIndex + 1) % modes.length] as typeof viewMode;
    }

    onMount(async () => {
        checkMobile();

        try {
            const response = await getSchema();
            if (response.schema) {
                schema = response.schema;

                // Set initial zoom based on table count
                const tableCount = Object.keys(schema).length;
                if (tableCount > 25) {
                    zoom = isMobile ? 0.4 : 0.6;
                } else if (tableCount > 15) {
                    zoom = isMobile ? 0.5 : 0.7;
                } else if (tableCount > 8) {
                    zoom = isMobile ? 0.6 : 0.8;
                }

                const container = document.getElementById('schema-container');
                if (container) {
                    containerWidth = container.clientWidth;
                    containerHeight = container.clientHeight;

                    calculateInitialPositions();

                    setTimeout(() => {
                        if (autoLayout) applyAutoLayout();
                        setTimeout(calculateLines, 200);
                    }, 100);

                    container.addEventListener('scroll', calculateLines);

                    const resizeObserver = new ResizeObserver(() => {
                        containerWidth = container.clientWidth;
                        containerHeight = container.clientHeight;
                        checkMobile();
                        calculateLines();
                    });
                    resizeObserver.observe(container);
                }
                console.log(schema);
            } else {
                error = "Schema data is missing.";
            }
        } catch (err) {
            error = `Error: ${err.message}`;
        }
    });

    function getTableSize() {
        const tableCount = schema ? Object.keys(schema).length : 0;

        if (viewMode === 'compact' || tableCount > 20) {
            return isMobile ? 'w-28' : 'w-36';
        } else if (viewMode === 'normal' || tableCount > 10) {
            return isMobile ? 'w-32' : 'w-44';
        } else {
            return isMobile ? 'w-40' : 'w-60';
        }
    }

    function getMaxColumns(tableInfo: any) {
        const tableCount = schema ? Object.keys(schema).length : 0;

        if (viewMode === 'compact' || tableCount > 20) {
            return 2;
        } else if (viewMode === 'normal' || tableCount > 10) {
            return 4;
        } else {
            return tableInfo.columns?.length || 0;
        }
    }

    function getMaxForeignKeys(tableInfo: any) {
        const tableCount = schema ? Object.keys(schema).length : 0;

        if (viewMode === 'compact' || tableCount > 20) {
            return 1;
        } else if (viewMode === 'normal' || tableCount > 10) {
            return 2;
        } else {
            return tableInfo.foreignKeys?.length || 0;
        }
    }
</script>

<svelte:window on:resize={checkMobile} />

<div class="flex flex-col h-full bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm p-2 sm:p-3 flex-shrink-0">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-2 sm:space-y-0">
            <div class="flex items-center">
                <Database class="text-blue-600 mr-2" size={isMobile ? 18 : 20} />
                <h1 class="text-base sm:text-lg font-semibold text-gray-800">Schema Visualizer</h1>
                {#if schema}
                    <span class="ml-2 text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded-full">
                        {Object.keys(schema).length} tables
                    </span>
                {/if}
            </div>

            <!-- Controls -->
            <div class="flex flex-wrap gap-1 text-xs">
                <button
                        class="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition flex items-center"
                        on:click={applyAutoLayout}
                >
                    <ArrowRightCircle size={12} class="mr-1" />
                    Auto
                </button>

                <button
                        class="px-2 py-1 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition flex items-center"
                        on:click={toggleGrid}
                >
                    <Grid3X3 size={12} class="mr-1" />
                    Grid
                </button>

                <button
                        class="px-2 py-1 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition flex items-center"
                        on:click={toggleRelations}
                >
                    {#if showRelations}
                        <Eye size={12} class="mr-1" />
                    {:else}
                        <EyeOff size={12} class="mr-1" />
                    {/if}
                    FK
                </button>

                <button
                        class="px-2 py-1 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition"
                        on:click={toggleViewMode}
                >
                    {viewMode === 'compact' ? 'Compact' : viewMode === 'normal' ? 'Normal' : 'Detailed'}
                </button>

                <div class="flex items-center">
                    <button
                            class="px-1 py-1 bg-gray-200 text-gray-800 rounded-l hover:bg-gray-300 transition"
                            on:click={zoomOut}
                    >
                        <ZoomOut size={12} />
                    </button>
                    <span class="text-xs text-gray-600 px-2 bg-gray-100">{Math.round(zoom * 100)}%</span>
                    <button
                            class="px-1 py-1 bg-gray-200 text-gray-800 rounded-r hover:bg-gray-300 transition"
                            on:click={zoomIn}
                    >
                        <ZoomIn size={12} />
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- Main Container -->
    <div class="flex-grow relative bg-gray-100 overflow-auto"
         id="schema-container"
         on:mousemove={onDragMove}
         on:mouseup={stopDrag}
         on:mouseleave={stopDrag}
         on:touchmove={handleTouchMove}
         on:touchend={handleTouchEnd}
         class:grid-bg={isGridVisible}
         style="transform-origin: top left; transform: scale({zoom});"
    >
        {#if error}
            <div class="m-4 p-4 text-red-500 bg-red-50 rounded-lg shadow text-sm">{error}</div>
        {:else if schema}
            <!-- SVG Connections -->
            {#if showRelations}
                <svg class="connections" width="100%" height="100%">
                    {#each lines as line}
                        <g>
                            <path
                                    d="M {line.start.x} {line.start.y}
                                   C {line.start.x + ((line.end.x - line.start.x) / 3)} {line.start.y},
                                     {line.end.x - ((line.end.x - line.start.x) / 3)} {line.end.y},
                                     {line.end.x} {line.end.y}"
                                    stroke="#4f46e5"
                                    stroke-width="1.5"
                                    fill="none"
                                    marker-end="url(#arrowhead)"
                                    class="relation-path"
                            />

                            {#if !isMobile && zoom > 0.7}
                                <g class="relation-label">
                                    <rect
                                            x={(line.start.x + line.end.x) / 2 - 30}
                                            y={(line.start.y + line.end.y) / 2 - 8}
                                            width="60"
                                            height="16"
                                            rx="3"
                                            fill="white"
                                            stroke="#e5e7eb"
                                    />
                                    <text
                                            x={(line.start.x + line.end.x) / 2}
                                            y={(line.start.y + line.end.y) / 2 + 3}
                                            text-anchor="middle"
                                            class="relation-text"
                                    >
                                        {line.label}
                                    </text>
                                </g>
                            {/if}
                        </g>
                    {/each}
                    <defs>
                        <marker
                                id="arrowhead"
                                markerWidth="6"
                                markerHeight="4"
                                refX="5"
                                refY="2"
                                orient="auto"
                        >
                            <polygon points="0 0, 6 2, 0 4" fill="#4f46e5" />
                        </marker>
                    </defs>
                </svg>
            {/if}

            <!-- Tables -->
            {#each Object.entries(schema) as [tableName, tableInfo]}
                <div class="table-item rounded-lg shadow-md bg-white border border-gray-200 {getTableSize()} {selectedTable === tableName ? 'selected' : ''}"
                     id="table-{tableName}"
                     style="top: {tablePositions[tableName]?.y || Object.keys(schema).indexOf(tableName) * (isMobile ? 120 : 160)}px;
                            left: {tablePositions[tableName]?.x || Object.keys(schema).indexOf(tableName) * (isMobile ? 120 : 160)}px;"
                     on:mousedown={(e) => startDrag(e, `table-${tableName}`)}
                     on:touchstart={(e) => handleTouchStart(e, `table-${tableName}`)}
                     on:click={() => selectTable(tableName)}
                >
                    <!-- Table Header -->
                    <div class="table-header bg-indigo-600 text-white p-1.5 sm:p-2 rounded-t-lg flex items-center justify-between">
                        <h2 class="text-xs sm:text-sm font-semibold truncate" title={tableName}>
                            {tableName.length > 12 ? tableName.substring(0, 12) + '...' : tableName}
                        </h2>
                        <span class="text-xs bg-indigo-800 px-1 py-0.5 rounded-full flex-shrink-0">
                            {tableInfo.columns?.length || 0}
                        </span>
                    </div>

                    <!-- Table Content -->
                    <div class="p-1.5 sm:p-2 space-y-1">
                        <!-- Primary Key -->
                        {#if tableInfo.primaryKey}
                            <div class="flex items-center">
                                <Shield size={10} class="text-yellow-500 mr-1 flex-shrink-0" />
                                <span class="text-xs px-1 py-0.5 bg-yellow-100 text-yellow-800 rounded truncate">
                                    {tableInfo.primaryKey.length > 8 ? tableInfo.primaryKey.substring(0, 8) + '...' : tableInfo.primaryKey}
                                </span>
                            </div>
                        {/if}

                        <!-- Columns -->
                        {#if tableInfo.columns?.length}
                            <div>
                                <div class="text-xs text-gray-600 mb-1">Columns:</div>
                                <div class="grid grid-cols-1 gap-0.5 max-h-16 overflow-y-auto">
                                    {#each tableInfo.columns.slice(0, getMaxColumns(tableInfo)) as column}
                                        <div class="text-xs py-0.5 px-1 bg-gray-50 rounded truncate" title={column}>
                                            {column.length > 12 ? column.substring(0, 12) + '...' : column}
                                        </div>
                                    {/each}
                                    {#if tableInfo.columns.length > getMaxColumns(tableInfo)}
                                        <div class="text-xs text-gray-500 px-1">
                                            +{tableInfo.columns.length - getMaxColumns(tableInfo)} more
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {/if}

                        <!-- Foreign Keys -->
                        {#if tableInfo.foreignKeys?.length}
                            <div>
                                <div class="text-xs text-gray-600 mb-1 flex items-center">
                                    <Link2 size={10} class="text-blue-500 mr-1 flex-shrink-0" />
                                    FK:
                                </div>
                                <div class="grid grid-cols-1 gap-0.5">
                                    {#each tableInfo.foreignKeys.slice(0, getMaxForeignKeys(tableInfo)) as fk}
                                        <div class="text-xs py-0.5 px-1 bg-blue-50 rounded text-blue-800 truncate" title="{fk.from} → {fk.table}.{fk.to}">
                                            {fk.from} → {fk.table}.{fk.to}
                                        </div>
                                    {/each}
                                    {#if tableInfo.foreignKeys.length > getMaxForeignKeys(tableInfo)}
                                        <div class="text-xs text-gray-500 px-1">
                                            +{tableInfo.foreignKeys.length - getMaxForeignKeys(tableInfo)} more
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        {:else}
            <div class="flex justify-center items-center h-full">
                <div class="text-center p-4 bg-white rounded-lg shadow m-4">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500 mx-auto mb-4"></div>
                    <p class="text-gray-600 text-sm">Loading database schema...</p>
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .table-item {
        position: absolute;
        cursor: move;
        min-height: 80px;
        user-select: none;
        z-index: 1;
        transition: box-shadow 0.2s ease, border-color 0.2s ease;
    }

    .table-item:hover {
        box-shadow: 0 4px 8px -2px rgba(0, 0, 0, 0.1);
    }

    .table-item.selected {
        border-color: #4f46e5;
        border-width: 2px;
        box-shadow: 0 0 0 2px rgba(79, 70, 229, 0.1);
    }

    .table-header {
        cursor: move;
        touch-action: none;
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

    .relation-text {
        font-size: 8px;
        fill: #4b5563;
    }

    .grid-bg {
        background-image: linear-gradient(to right, rgba(0, 0, 0, 0.05) 1px, transparent 1px),
        linear-gradient(to bottom, rgba(0, 0, 0, 0.05) 1px, transparent 1px);
        background-size: 15px 15px;
    }

    @media (max-width: 767px) {
        .table-item * {
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            user-select: none;
        }
    }
</style>