<script lang="ts">
    import { onMount } from "svelte";
    import { executeQueryAPI, getSchemaVisualisationSvg } from '$lib/api/database';
    import {
        ChevronRight,
        Loader,
        Terminal,
        BarChart3,
        RefreshCw,
        AlertCircle,
        MessageSquare
    } from 'lucide-svelte';
    import Notification from '$lib/components/Notification.svelte';

    // State variables
    let naturalLanguageQuery: string = "";
    let sqlQuery: string = "";
    let svgContent: string | null = null;
    let queryResult: any = null;
    let isLoading: boolean = false;
    let isProcessing: boolean = false;
    let errorMessage: string | null = null;
    let stage: 'input' | 'sql' | 'result' = 'input';
    let examples: string[] = [
        "Show me all users and their email addresses",
        "How many orders were placed last month?",
        "What are the top 5 products by sales?",
        "List all customers from New York"
    ];

    let notificationMessage: string = "";
    let notificationType: 'success' | 'error' = 'success';

    async function processNaturalLanguageQuery() {
        try {
            isLoading = true;
            errorMessage = null;

            const response = await fetch('/text-to-sql/gpt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query: naturalLanguageQuery })
            });

            if (!response.ok) {
                throw new Error('Failed to convert natural language to SQL');
            }

            const data = await response.json();
            sqlQuery = data.query;
            stage = 'sql';

            notificationMessage = "Successfully converted to SQL query";
            notificationType = 'success';

        } catch (error) {
            console.error("Error in natural language processing:", error);
            errorMessage = error instanceof Error ? error.message : 'An unknown error occurred';

            notificationMessage = "Failed to convert query";
            notificationType = 'error';
        } finally {
            isLoading = false;
        }
    }

    async function executeSQL() {
        try {
            isProcessing = true;
            errorMessage = null;

            // Step 2: Execute SQL query
            const result = await executeQueryAPI(sqlQuery);
            queryResult = result.result;

            // Step 3: If it's a SELECT query, get visualization
            if (sqlQuery.trim().toUpperCase().startsWith('SELECT')) {
                const visualResult = await getSchemaVisualisationSvg(sqlQuery);
                svgContent = visualResult.svg;

                // Notification for visualization
                notificationMessage = "Query executed and visualization generated";
                notificationType = 'success';
            } else {
                svgContent = null;

                // Notification for non-SELECT queries
                notificationMessage = "Query executed successfully";
                notificationType = 'success';
            }

            stage = 'result';

        } catch (error) {
            console.error("Error executing SQL or generating visualization:", error);
            errorMessage = error instanceof Error ? error.message : 'An unknown error occurred';

            // Show error notification
            notificationMessage = "Failed to execute query";
            notificationType = 'error';
        } finally {
            isProcessing = false;
        }
    }

    function resetQuery() {
        stage = 'input';
        queryResult = null;
        svgContent = null;
        sqlQuery = "";
        errorMessage = null;

        // Optional notification for reset
        notificationMessage = "Started a new query";
        notificationType = 'success';
    }

    function useExample(example: string) {
        naturalLanguageQuery = example;

        // Subtle feedback
        notificationMessage = "Example selected";
        notificationType = 'success';
    }
</script>

<Notification message={notificationMessage} type={notificationType} />

<div class="container mx-auto p-4 max-w-4xl">
    <h1 class="text-3xl font-bold mb-6 text-gray-800">Data Visualization</h1>

    {#if errorMessage}
        <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-6 rounded">
            <div class="flex">
                <AlertCircle class="h-6 w-6 text-red-500 mr-2" />
                <p class="text-red-700">{errorMessage}</p>
            </div>
        </div>
    {/if}

    <div class="mb-8">
        <div class="flex items-center mb-2">
            <span class="w-8 h-8 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mr-2">1</span>
            <h2 class="text-xl font-semibold text-gray-700">Natural Language Query</h2>
        </div>

        <div class="bg-white p-6 rounded-lg shadow-md mb-4">
            {#if stage === 'input'}
                <div class="mb-4">
                    <div class="flex items-start">
                        <MessageSquare class="w-5 h-5 text-gray-500 mt-1 mr-2" />
                        <textarea
                                bind:value={naturalLanguageQuery}
                                class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400 focus:border-transparent"
                                rows="3"
                                placeholder="Describe what data you want to see in plain English..."
                                disabled={isLoading}
                        ></textarea>
                    </div>
                </div>

                <div class="mb-4">
                    <p class="text-sm text-gray-500 mb-2">Examples:</p>
                    <div class="flex flex-wrap gap-2">
                        {#each examples as example}
                            <button
                                    class="text-xs bg-gray-100 hover:bg-gray-200 px-3 py-1 rounded-full text-gray-700"
                                    on:click={() => useExample(example)}
                            >
                                {example}
                            </button>
                        {/each}
                    </div>
                </div>

                <div class="flex justify-end">
                    <button
                            class="px-4 py-2 bg-blue-600 text-white rounded-lg flex items-center disabled:bg-blue-300"
                            on:click={processNaturalLanguageQuery}
                            disabled={!naturalLanguageQuery.trim() || isLoading}
                    >
                        {#if isLoading}
                            <Loader class="w-4 h-4 mr-2 animate-spin" />
                            Processing...
                        {:else}
                            Convert to SQL
                            <ChevronRight class="w-4 h-4 ml-1" />
                        {/if}
                    </button>
                </div>
            {:else}
                <div class="flex items-center justify-between bg-gray-50 p-3 rounded mb-2">
                    <div class="text-gray-700">{naturalLanguageQuery}</div>
                    <button
                            class="text-blue-600 text-sm hover:underline flex items-center"
                            on:click={resetQuery}
                    >
                        <RefreshCw class="w-3 h-3 mr-1" />
                        Change
                    </button>
                </div>
            {/if}
        </div>
    </div>

    {#if stage === 'sql' || stage === 'result'}
        <div class="mb-8">
            <div class="flex items-center mb-2">
                <span class="w-8 h-8 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mr-2">2</span>
                <h2 class="text-xl font-semibold text-gray-700">Generated SQL</h2>
            </div>

            <div class="bg-white p-6 rounded-lg shadow-md mb-4">
                {#if stage === 'sql'}
                    <div class="mb-4">
                        <div class="flex items-start">
                            <Terminal class="w-5 h-5 text-gray-500 mt-1 mr-2" />
                            <textarea
                                    bind:value={sqlQuery}
                                    class="w-full p-3 border border-gray-300 rounded-lg font-mono text-sm bg-gray-50 focus:ring-2 focus:ring-blue-400 focus:border-transparent"
                                    rows="4"
                                    placeholder="SQL query will appear here..."
                                    disabled={isProcessing}
                            ></textarea>
                        </div>
                    </div>

                    <div class="flex justify-end">
                        <button
                                class="px-4 py-2 bg-blue-600 text-white rounded-lg flex items-center disabled:bg-blue-300"
                                on:click={executeSQL}
                                disabled={!sqlQuery.trim() || isProcessing}
                        >
                            {#if isProcessing}
                                <Loader class="w-4 h-4 mr-2 animate-spin" />
                                Executing...
                            {:else}
                                Execute SQL
                                <ChevronRight class="w-4 h-4 ml-1" />
                            {/if}
                        </button>
                    </div>
                {:else}
                    <div class="flex items-center justify-between bg-gray-50 p-3 rounded mb-2 font-mono text-sm overflow-x-auto">
                        <div class="text-gray-700">{sqlQuery}</div>
                        <button
                                class="text-blue-600 text-sm hover:underline flex items-center ml-3"
                                on:click={() => stage = 'sql'}
                        >
                            <RefreshCw class="w-3 h-3 mr-1" />
                            Edit
                        </button>
                    </div>
                {/if}
            </div>
        </div>
    {/if}

    {#if stage === 'result'}
        <div>
            <div class="flex items-center mb-2">
                <span class="w-8 h-8 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mr-2">3</span>
                <h2 class="text-xl font-semibold text-gray-700">Results</h2>
            </div>

            <div class="bg-white p-6 rounded-lg shadow-md">
                {#if sqlQuery.trim().toUpperCase().startsWith('SELECT')}
                    <div class="mb-6">
                        <h3 class="text-lg font-medium mb-2 flex items-center">
                            <BarChart3 class="w-5 h-5 mr-2 text-blue-600" />
                            Visualization
                        </h3>

                        {#if svgContent}
                            <div class="svg-container border border-gray-200 rounded-lg p-4 bg-gray-50 overflow-x-auto">
                                {@html svgContent}
                            </div>
                        {:else}
                            <div class="text-gray-500 italic text-center p-6 border border-dashed border-gray-300 rounded-lg">
                                <Loader class="w-6 h-6 mx-auto mb-2 text-gray-400 animate-spin" />
                                <p>Generating visualization...</p>
                            </div>
                        {/if}
                    </div>
                {/if}

                <div>
                    <h3 class="text-lg font-medium mb-2">Data Results</h3>

                    {#if queryResult && queryResult.row_count > 0}
                        <div class="overflow-x-auto">
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead class="bg-gray-50">
                                <tr>
                                    {#each queryResult.columns as column}
                                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                            {column}
                                        </th>
                                    {/each}
                                </tr>
                                </thead>
                                <tbody class="bg-white divide-y divide-gray-200">
                                {#each queryResult.result as row}
                                    <tr>
                                        {#each Object.values(row) as cell}
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                {cell}
                                            </td>
                                        {/each}
                                    </tr>
                                {/each}
                                </tbody>
                            </table>
                        </div>
                        <div class="mt-2 text-sm text-gray-500">
                            Showing {queryResult.result.length} of {queryResult.row_count} rows
                        </div>
                    {:else if sqlQuery.trim().toUpperCase().startsWith('SELECT')}
                        <div class="text-gray-500 italic text-center p-6 border border-dashed border-gray-300 rounded-lg">
                            No data returned from this query
                        </div>
                    {:else}
                        <div class="text-green-600 p-4 bg-green-50 rounded-lg flex items-center">
                            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                            </svg>
                            Query executed successfully
                        </div>
                    {/if}
                </div>

                <div class="mt-6 flex justify-end">
                    <button
                            class="px-4 py-2 bg-blue-600 text-white rounded-lg flex items-center"
                            on:click={resetQuery}
                    >
                        New Query
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .svg-container {
        width: 100%;
        height: auto;
        max-width: 100%;
    }

    .svg-container svg {
        width: 100%;
        height: auto;
    }
</style>