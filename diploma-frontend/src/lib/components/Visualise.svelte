<script lang="ts">
    import { onMount } from 'svelte';
    import { getSchema } from '$lib/api';
    import type { Schema } from '$lib/types/table';

    let schema: Schema | null = null;
    let error: string | null = null;
    let draggingTable: string | null = null;
    let initialMousePosition = { x: 0, y: 0 };
    let initialTablePosition = { x: 0, y: 0 };
    let lines: Array<{start: any, end: any, label: string}> = [];

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
        Object.entries(schema).forEach(([tableName, tableInfo]) => {
            console.log(`Checking table: ${tableName}`); // Отладочное сообщение
            if (!tableInfo.foreignKeys?.length) {
                console.log(`No foreign keys found for ${tableName}`); // Отладочное сообщение
                return;
            }

            const sourceElement = document.getElementById(`table-${tableName}`);
            if (!sourceElement) return;

            tableInfo.foreignKeys.forEach(fk => {
                const targetElement = document.getElementById(`table-${fk.Table}`);
                if (!targetElement) return;

                const sourcePos = getTablePosition(sourceElement);
                const targetPos = getTablePosition(targetElement);

                if (!sourcePos || !targetPos) return;

                const start = {
                    x: sourcePos.left + sourcePos.width,
                    y: sourcePos.top + (sourcePos.height / 2)
                };

                const end = {
                    x: targetPos.left,
                    y: targetPos.top + (targetPos.height / 2)
                };

                console.log(`Adding line from ${start.x},${start.y} to ${end.x},${end.y}`); // Отладочное сообщение

                newLines.push({
                    start,
                    end,
                    label: `${fk.From} → ${fk.To}`
                });
            });
        });

        lines = newLines;
        console.log("Lines updated:", lines);
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

                // Recalculate lines when table is moved
                requestAnimationFrame(calculateLines);
            }
        }
    }

    function stopDrag() {
        draggingTable = null;
        calculateLines();
    }

    onMount(async () => {
        try {
            const response = await getSchema();
            if (response.schema) {
                schema = response.schema;
                // Initial line calculation after schema is loaded
                setTimeout(calculateLines, 100); // Small delay to ensure DOM is ready
            } else {
                error = "Schema data is missing.";
            }
        } catch (err) {
            error = `Error: ${err.message}`;
        }

        // Add scroll event listener to update lines when container is scrolled
        const container = document.getElementById('schema-container');
        if (container) {
            container.addEventListener('scroll', calculateLines);
        }
    });
</script>

<div class="flex flex-col h-screen">
    <div class="flex-grow relative bg-gray-100 p-4 min-h-[800px] overflow-auto"
         id="schema-container"
         on:mousemove={onDragMove}
         on:mouseup={stopDrag}
         on:mouseleave={stopDrag}
    >
        {#if error}
            <div class="text-red-500">{error}</div>
        {:else if schema}

            <svg class="connections" width="100%" height="100%">
                {#each lines as line}
                    <g>
                        <path
                                d="M {line.start.x} {line.start.y}
                               Q {(line.start.x + line.end.x) / 2} {line.start.y},
                                 {(line.start.x + line.end.x) / 2} {(line.start.y + line.end.y) / 2}
                               Q {(line.start.x + line.end.x) / 2} {line.end.y},
                                 {line.end.x} {line.end.y}"
                                stroke="#666"
                                stroke-width="2"
                                fill="none"
                                marker-end="url(#arrowhead)"
                        />

                        <text
                                x={(line.start.x + line.end.x) / 2}
                                y={(line.start.y + line.end.y) / 2 - 10}
                                text-anchor="middle"
                                class="relationship-label"
                        >
                            {line.label}
                        </text>
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
                        <polygon points="0 0, 10 3.5, 0 7" fill="#666" />
                    </marker>
                </defs>
            </svg>


            {#each Object.entries(schema) as [tableName, tableInfo]}
                <div class="table-item rounded-lg shadow-md"
                     id="table-{tableName}"
                     style="top: {Object.keys(schema).indexOf(tableName) * 220}px; left: {Object.keys(schema).indexOf(tableName) * 220}px;"
                     on:mousedown={(e) => startDrag(e, `table-${tableName}`)}>
                    <div class="bg-white p-4 table-box">
                        <h2 class="text-xl font-semibold mb-2">{tableName}</h2>
                        <p class="text-gray-700">Primary Key: {tableInfo.primaryKey}</p>
                        <p class="text-gray-700">Columns:</p>
                        <ul class="list-disc list-inside mb-2">
                            {#each tableInfo.columns as column}
                                <li>{column}</li>
                            {/each}
                        </ul>
                        <p class="text-gray-700">Foreign Keys:</p>
                        <ul class="list-disc list-inside">
                            {#each tableInfo.foreignKeys as fk}
                                <li>
                                    {fk.From} → <strong>{fk.Table}.{fk.To}</strong>
                                </li>
                            {/each}
                        </ul>
                    </div>
                </div>
            {/each}
        {/if}
    </div>
</div>

<style>
    .table-item {
        position: absolute;
        cursor: move;
        width: 280px;
        min-height: 200px;
        user-select: none;
        background: white;
        z-index: 1;
    }

    .table-box {
        padding: 1rem;
        height: 100%;
        border-radius: 0.5rem;
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

    .relationship-label {
        font-size: 12px;
        fill: #666;
        pointer-events: none;
        background: white;
        padding: 2px;
    }

    :global(body) {
        margin: 0;
        padding: 0;
    }
</style>
