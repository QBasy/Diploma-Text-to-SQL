<script lang="ts">
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import {
        createTableAPI,
        executeQueryAPI,
        generateComplexSQLbyGPT,
        getComplexSchema
    } from "$lib/api";
    import { type Table } from "$lib/types/table"
    import Visualise from "$lib/components/Visualise.svelte";

    let activeTab = 'query';
    let sqlQuery = '';
    let naturalLanguageQuery = '';
    let queryResult: any = null;
    let queryResultType: 'success' | 'error' = 'success';
    let loading = false;
    let notification = '';
    let notificationType = 'success';
    let tables: Table[] = [];
    let showResults = false;

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

    function formatQueryResult(result: any): string {
        if (typeof result === 'string') {
            return result;
        }
        if (Array.isArray(result)) {
            if (result.length === 0) {
                return 'No results found';
            }
            return JSON.stringify(result, null, 2);
        }
        if (typeof result === 'object') {
            return JSON.stringify(result, null, 2);
        }
        return String(result);
    }

    async function executeSQLQuery() {
        loading = true;
        try {
            const result = await executeQueryAPI(sqlQuery);
            queryResult = result;
            queryResultType = 'success';
            showResults = true;
            notification = 'Query executed successfully';
            notificationType = 'success';
        } catch (error) {
            queryResult = error instanceof Error ? error.message : 'Unknown error occurred';
            queryResultType = 'error';
            showResults = true;
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
            const { sql } = await generateComplexSQLbyGPT(naturalLanguageQuery+", sqlite", schema);
            const result = await executeQueryAPI(sql);
            queryResult = { generatedSQL: sql, result: result };
            queryResultType = 'success';
            showResults = true;
            notification = 'Natural language query processed successfully';
            notificationType = 'success';
        } catch (error) {
            queryResult = error instanceof Error ? error.message : 'Unknown error occurred';
            queryResultType = 'error';
            showResults = true;
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
            queryResult = `Table "${tableForm.name}" created successfully with columns: ${tableForm.columns.map(col => `${col.name} (${col.type})`).join(', ')}`;
            queryResultType = 'success';
            showResults = true;
            notification = `Table ${tableForm.name} created successfully`;
            notificationType = 'success';
            tables = await getComplexSchema().then(schema => schema.tables);
            tableForm = { name: '', columns: [{ name: '', type: 'TEXT', isForeignKey: false, referencedTable: '', referencedColumn: '' }] };
        } catch (error) {
            queryResult = error instanceof Error ? error.message : 'Unknown error occurred';
            queryResultType = 'error';
            showResults = true;
            notification = 'Error creating table';
            notificationType = 'error';
        } finally {
            loading = false;
        }
    }

    function clearResults() {
        queryResult = null;
        showResults = false;
    }
</script>

<svelte:head>
    <title>User Database Panel</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
    <div class="container max-w-full sm:max-w-6xl mx-auto p-3 sm:p-4 space-y-4">
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
            <h1 class="text-xl sm:text-2xl font-bold text-gray-800 mb-4 text-center sm:text-left">Database Management</h1>

            <!-- Mobile-optimized tabs -->
            <div class="flex flex-wrap gap-1 mb-4 overflow-x-auto">
                <button
                        class="flex-shrink-0 px-3 py-2 text-xs sm:text-sm rounded-lg font-medium transition-colors duration-200 {activeTab === 'query' ? 'bg-blue-600 text-white shadow' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                        on:click={() => activeTab = 'query'}
                >
                    SQL Query
                </button>
                <button
                        class="flex-shrink-0 px-3 py-2 text-xs sm:text-sm rounded-lg font-medium transition-colors duration-200 {activeTab === 'natural' ? 'bg-blue-600 text-white shadow' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                        on:click={() => activeTab = 'natural'}
                >
                    Natural Language
                </button>
                <button
                        class="flex-shrink-0 px-3 py-2 text-xs sm:text-sm rounded-lg font-medium transition-colors duration-200 {activeTab === 'table' ? 'bg-blue-600 text-white shadow' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                        on:click={() => activeTab = 'table'}
                >
                    Create Table
                </button>
                <a class="flex-shrink-0 px-3 py-2 text-xs sm:text-sm rounded-lg font-medium transition-colors duration-200 bg-gray-100 text-gray-700 hover:bg-gray-200"
                   href="/profile/database/custom">
                    Custom Database
                </a>
            </div>

            {#if loading}
                <div class="flex justify-center items-center min-h-[120px]">
                    <LoadingSpinner />
                </div>
            {:else}
                {#if activeTab === 'query'}
                    <div class="space-y-3">
                        <textarea
                                bind:value={sqlQuery}
                                class="w-full h-24 sm:h-32 p-3 border border-gray-200 rounded-lg font-mono text-xs sm:text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="Enter your SQL query"
                        ></textarea>
                        <div class="flex flex-col sm:flex-row gap-2">
                            <button
                                    on:click={executeSQLQuery}
                                    class="flex-1 bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-4 text-sm rounded-lg transition-colors duration-200 shadow"
                            >
                                Execute Query
                            </button>
                            {#if showResults}
                                <button
                                        on:click={clearResults}
                                        class="sm:w-auto bg-gray-500 hover:bg-gray-600 text-white font-medium py-3 px-4 text-sm rounded-lg transition-colors duration-200"
                                >
                                    Clear Results
                                </button>
                            {/if}
                        </div>
                    </div>
                {:else if activeTab === 'natural'}
                    <div class="space-y-3">
                        <textarea
                                bind:value={naturalLanguageQuery}
                                class="w-full h-24 sm:h-32 p-3 border border-gray-200 rounded-lg text-xs sm:text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="Enter your natural language query (e.g., 'Show all users with age greater than 25')"
                        ></textarea>
                        <div class="flex flex-col sm:flex-row gap-2">
                            <button
                                    on:click={executeNaturalQuery}
                                    class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 text-sm rounded-lg transition-colors duration-200 shadow"
                            >
                                Generate and Execute
                            </button>
                            {#if showResults}
                                <button
                                        on:click={clearResults}
                                        class="sm:w-auto bg-gray-500 hover:bg-gray-600 text-white font-medium py-3 px-4 text-sm rounded-lg transition-colors duration-200"
                                >
                                    Clear Results
                                </button>
                            {/if}
                        </div>
                    </div>
                {:else if activeTab === 'table'}
                    <div class="space-y-3">
                        <input
                                type="text"
                                bind:value={tableForm.name}
                                placeholder="Table Name"
                                class="w-full p-3 text-sm border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >

                        <div class="space-y-2 max-h-48 sm:max-h-64 overflow-y-auto bg-gray-50 p-3 rounded-lg">
                            <h3 class="text-sm font-medium text-gray-700 mb-2">Columns</h3>
                            {#each tableForm.columns as column, index}
                                <div class="flex flex-col sm:flex-row gap-2 bg-white p-2 rounded border">
                                    <input
                                            type="text"
                                            bind:value={column.name}
                                            placeholder="Column Name"
                                            class="flex-1 p-2 text-sm border border-gray-200 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    >
                                    <select
                                            bind:value={column.type}
                                            class="flex-1 sm:flex-none sm:w-24 p-2 text-sm border border-gray-200 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                                    >
                                        <option value="TEXT">TEXT</option>
                                        <option value="INTEGER">INTEGER</option>
                                        <option value="REAL">REAL</option>
                                    </select>
                                    {#if tableForm.columns.length > 1}
                                        <button
                                                on:click={() => removeColumn(index)}
                                                class="w-full sm:w-auto p-2 bg-red-500 hover:bg-red-600 text-white text-sm rounded transition-colors duration-200"
                                        >
                                            Remove
                                        </button>
                                    {/if}
                                </div>
                            {/each}
                        </div>

                        <div class="flex flex-col sm:flex-row gap-2">
                            <button
                                    on:click={addColumn}
                                    class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-3 text-sm rounded-lg transition-colors duration-200 shadow"
                            >
                                Add Column
                            </button>
                            <button
                                    on:click={createTable}
                                    class="flex-1 bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-3 text-sm rounded-lg transition-colors duration-200 shadow"
                            >
                                Create Table
                            </button>
                        </div>

                        {#if showResults}
                            <button
                                    on:click={clearResults}
                                    class="w-full bg-gray-500 hover:bg-gray-600 text-white font-medium py-2 px-4 text-sm rounded-lg transition-colors duration-200"
                            >
                                Clear Results
                            </button>
                        {/if}
                    </div>
                {/if}
            {/if}
        </div>

        <!-- Results Section - Always visible when there are results -->
        {#if showResults && queryResult !== null}
            <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
                <div class="flex items-center justify-between mb-3">
                    <h2 class="text-lg font-semibold text-gray-800">Results</h2>
                    <span class="px-2 py-1 text-xs rounded-full {queryResultType === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                        {queryResultType === 'success' ? 'Success' : 'Error'}
                    </span>
                </div>

                {#if activeTab === 'natural' && queryResultType === 'success' && queryResult.generatedSQL}
                    <div class="mb-4">
                        <h3 class="text-sm font-medium text-gray-700 mb-2">Generated SQL:</h3>
                        <pre class="bg-gray-900 text-green-400 p-3 rounded-lg overflow-x-auto text-xs font-mono border">{queryResult.generatedSQL}</pre>
                    </div>
                    <div>
                        <h3 class="text-sm font-medium text-gray-700 mb-2">Query Result:</h3>
                        <pre class="bg-gray-50 p-3 rounded-lg overflow-x-auto text-xs font-mono border max-h-64 sm:max-h-96 overflow-y-auto">{formatQueryResult(queryResult.result)}</pre>
                    </div>
                {:else}
                    <pre class="bg-gray-50 p-3 rounded-lg overflow-x-auto text-xs font-mono border max-h-64 sm:max-h-96 overflow-y-auto {queryResultType === 'error' ? 'text-red-700 bg-red-50' : ''}">{formatQueryResult(queryResult)}</pre>
                {/if}
            </div>
        {/if}

        <!-- Notifications -->
        {#if notification}
            <Notification message={notification} type={notificationType} />
        {/if}

        <!-- Visualization Section -->
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
            <h2 class="text-lg font-semibold text-gray-800 mb-4">Database Visualization</h2>
            <div class="h-64 sm:h-96 w-full overflow-auto border rounded-lg bg-gray-50">
                <Visualise />
            </div>
        </div>
    </div>
</div>

<style>
    /* Custom scrollbar for better mobile experience */
    :global(.overflow-x-auto::-webkit-scrollbar) {
        height: 4px;
    }

    :global(.overflow-x-auto::-webkit-scrollbar-track) {
        background: #f1f5f9;
        border-radius: 2px;
    }

    :global(.overflow-x-auto::-webkit-scrollbar-thumb) {
        background: #cbd5e1;
        border-radius: 2px;
    }

    :global(.overflow-x-auto::-webkit-scrollbar-thumb:hover) {
        background: #94a3b8;
    }

    /* Ensure proper text wrapping on mobile */
    @media (max-width: 640px) {
        pre {
            word-wrap: break-word;
            white-space: pre-wrap;
        }
    }
</style>