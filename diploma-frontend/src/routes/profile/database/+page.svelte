<script lang="ts">
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import { createTableAPI, executeQueryAPI, generateComplexSQL, getComplexSchema } from "$lib/api";
    import { type Table } from "$lib/types/table"
    import Visualise from "$lib/components/Visualise.svelte";
    import DatabaseGraph from "$lib/components/DatabaseGraph.svelte";

    let activeTab = 'query';
    let sqlQuery = '';
    let naturalLanguageQuery = '';
    let queryResult: string | null = null;
    let loading = false;
    let notification = '';
    let notificationType = 'success';
    let tables: Table[] = [];

    let tableForm = {
        name: '',
        columns: [{ name: '', type: 'TEXT', isForeignKey: false, referencedTable: '', referencedColumn: '' }]
    };

    function addColumn() {
        tableForm.columns = [...tableForm.columns, { name: '', type: 'TEXT', isForeignKey: false, referencedTable: '', referencedColumn: '' }];
    }

    function removeColumn(index: number) {
        tableForm.columns = tableForm.columns.filter((_, i) => i !== index);
    }

    async function executeSQLQuery() {
        loading = true;
        try {
            const result = await executeQueryAPI(sqlQuery);
            queryResult = 'Query executed successfully';
            notification = 'Query processed successfully: ' + result;
            notificationType = 'success';
        } catch (error) {
            notification = 'Error executing query';
            notificationType = 'error';
        } finally {
            loading = false;
        }
    }

    async function executeNaturalQuery() {
        loading = true;
        try {
            const schema = await getComplexSchema();
            tables = schema.tables;
            const { sql } = await generateComplexSQL(naturalLanguageQuery, [schema]);
            const result = await executeQueryAPI(sql);
            queryResult = 'Query executed successfully, Result: ' + result;
            notification = 'Query processed successfully: ' + result;
            notificationType = 'success';
        } catch (error) {
            notification = 'Error executing natural language query';
            notificationType = 'error';
        } finally {
            loading = false;
        }
    }

    async function createTable() {
        if (!tableForm.name.trim()) {
            notification = 'Error: Table name is required';
            notificationType = 'error';
            return;
        }

        if (tableForm.columns.some(col => !col.name.trim())) {
            notification = 'Error: Column names cannot be empty';
            notificationType = 'error';
            return;
        }

        loading = true;
        try {
            await createTableAPI(tableForm);
            notification = `Table ${tableForm.name} created successfully`;
            notificationType = 'success';
            tables = await getComplexSchema().then(schema => schema.tables); // Обновляем таблицы
            tableForm = { name: '', columns: [{ name: '', type: 'TEXT', isForeignKey: false, referencedTable: '', referencedColumn: '' }] };
        } catch (error) {
            notification = 'Error creating table';
            notificationType = 'error';
        } finally {
            loading = false;
        }
    }
</script>

<div class="container w-4/5 mx-auto p-8 space-y-6">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">Database Management</h1>

    <div class="flex gap-2 mb-8">
        <button
                class="px-6 py-3 rounded-lg font-medium transition-colors duration-200 {activeTab === 'query' ? 'bg-blue-600 text-white shadow-lg' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                on:click={() => activeTab = 'query'}
        >
            SQL Query
        </button>
        <button
                class="px-6 py-3 rounded-lg font-medium transition-colors duration-200 {activeTab === 'natural' ? 'bg-blue-600 text-white shadow-lg' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                on:click={() => activeTab = 'natural'}
        >
            Natural Language Query
        </button>
        <button
                class="px-6 py-3 rounded-lg font-medium transition-colors duration-200 {activeTab === 'table' ? 'bg-blue-600 text-white shadow-lg' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                on:click={() => activeTab = 'table'}
        >
            Create Table
        </button>
    </div>

    {#if loading}
        <div class="flex justify-center items-center min-h-[200px]">
            <LoadingSpinner />
        </div>
    {:else}
        {#if activeTab === 'query'}
            <div class="bg-white shadow-lg rounded-xl p-8 space-y-4">
                <textarea
                        bind:value={sqlQuery}
                        class="w-full h-48 p-4 border border-gray-200 rounded-lg font-mono text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="Enter your SQL query"
                ></textarea>
                <button
                        on:click={executeSQLQuery}
                        class="w-full bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200 shadow-md"
                >
                    Execute Query
                </button>
            </div>
        {:else if activeTab === 'natural'}
            <div class="bg-white shadow-lg rounded-xl p-8 space-y-4">
                <textarea
                        bind:value={naturalLanguageQuery}
                        class="w-full h-48 p-4 border border-gray-200 rounded-lg text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="Enter your natural language query"
                ></textarea>
                <button
                        on:click={executeNaturalQuery}
                        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200 shadow-md"
                >
                    Generate and Execute
                </button>
            </div>
        {:else if activeTab === 'table'}
            <div class="bg-white shadow-lg rounded-xl p-8 space-y-6">
                <input
                        type="text"
                        bind:value={tableForm.name}
                        placeholder="Table Name"
                        class="w-full p-4 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                <div class="space-y-4">
                    {#each tableForm.columns as column, index}
                        <div class="flex gap-4 items-center">
                            <input
                                    type="text"
                                    bind:value={column.name}
                                    placeholder="Column Name"
                                    class="flex-1 p-4 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            >
                            <select
                                    bind:value={column.type}
                                    class="p-4 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                            >
                                <option value="TEXT">TEXT</option>
                                <option value="INTEGER">INTEGER</option>
                                <option value="REAL">REAL</option>
                            </select>
                            {#if tableForm.columns.length > 1}
                                <button
                                        on:click={() => removeColumn(index)}
                                        class="p-4 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors duration-200"
                                >
                                    Remove
                                </button>
                            {/if}
                        </div>
                    {/each}
                </div>
                <div class="flex gap-4">
                    <button
                            on:click={addColumn}
                            class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200 shadow-md"
                    >
                        Add Column
                    </button>
                    <button
                            on:click={createTable}
                            class="flex-1 bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200 shadow-md"
                    >
                        Create Table
                    </button>
                </div>
            </div>
        {/if}
    {/if}

    {#if queryResult}
        <div class="mt-8 bg-white shadow-lg rounded-xl p-8">
            <h2 class="text-xl font-semibold text-gray-800 mb-4">Query Result</h2>
            <pre class="bg-gray-50 p-6 rounded-lg overflow-x-auto font-mono text-sm">{queryResult}</pre>
        </div>
    {/if}

    {#if notification}
        <Notification message={notification} type={notificationType} />
    {/if}

    <div class="mt-8 bg-white shadow-lg rounded-xl p-8">
        <Visualise />
    </div>

    <div class="mt-8 bg-white shadow-lg rounded-xl p-8">
        <DatabaseGraph />
    </div>
</div>