<script lang="ts">
    import { onMount } from 'svelte';
    import cytoscape from 'cytoscape';
    import type { Schema } from '$lib/api/database';
    import {getComplexSchema, getSchema} from '$lib/api/database';

    let cy: any;

    const fetchSchema = async () => {
        const schemaResponse = await getSchema();
        if (schemaResponse.status === 'success' && schemaResponse.schema) {
            formatSchema(schemaResponse.schema);
        } else {
            console.error('Error fetching schema:', schemaResponse.message);
        }
    };


    function formatSchema(schema: Schema): any[] {
        let elements: any[] = [];

        if (!schema.tables) {
            console.error("No tables found in schema");
            return elements;
        }

        for (const table of schema.tables) {
            elements.push({
                data: { id: table.name, label: table.name },
            });

            for (const column of table.columns) {
                elements.push({
                    data: { id: `${table.name}.${column.name}`, label: column.name, parent: table.name },
                });
            }

            for (const column of table.columns) {
                if (column.isForeignKey && column.referencedTable) {
                    elements.push({
                        data: {
                            id: `edge-${table.name}-${column.referencedTable}`,
                            source: table.name,
                            target: column.referencedTable,
                            label: `${column.name} → ${column.referencedColumn}`,
                        },
                    });
                }
            }
        }
        return elements;
    }


    onMount(async () => {
        const elements = await fetchSchema();

        if (!cy) {
            cy = cytoscape({
                container: document.getElementById('cy'),
                elements: elements,
                style: [
                    {
                        selector: 'node',
                        style: {
                            'label': 'data(label)',
                            'text-valign': 'center',
                            'color': '#fff',
                            'background-color': '#007bff',
                            'width': '120px',
                            'height': '50px',
                            'text-wrap': 'wrap',
                            'border-width': 2,
                            'border-color': '#333',
                        },
                    },
                    {
                        selector: 'edge',
                        style: {
                            'curve-style': 'bezier',
                            'target-arrow-shape': 'triangle',
                            'line-color': '#ccc',
                            'target-arrow-color': '#ccc',
                            'label': 'data(label)',
                            'text-background-color': '#fff',
                            'text-background-opacity': 1,
                        },
                    },
                ],
                layout: {
                    name: 'cose',
                },
            });
        } else {
            cy.json({ elements });
        }
    });
</script>

<div id="cy" class="w-full h-[500px] border border-gray-300 rounded"></div>

<style>
    #cy {
        width: 100%;
        height: 500px;
        border: 1px solid gray;
        border-radius: 5px;
    }
</style>
