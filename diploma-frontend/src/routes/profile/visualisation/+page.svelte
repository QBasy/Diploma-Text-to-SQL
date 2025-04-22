<script lang="ts">
    import type { PageData } from '../../../../.svelte-kit/types/src/routes';
    export let data: PageData;

    import Navbar from '$lib/components/Navbar.svelte';
    import Footer from '$lib/components/Footer.svelte';
    import {executeQueryVisualisation, type TableSchema, type ColumnInfo} from '$lib/api/database';
    import {
        generateComplexSQLbyGPT,
        getCustomSchemaComplex,
        listCustomDatabases,
        type TextToSQLRequest
    } from '$lib/api';
    import {
        ChevronRight,
        Loader,
        Terminal,
        BarChart3,
        RefreshCw,
        AlertCircle,
        MessageSquare,
        Database
    } from 'lucide-svelte';
    import Notification from '$lib/components/Notification.svelte';
    import type {Schema, Table, Column} from "$lib/types/table";
    import Spinner from "$lib/components/Spinner.svelte";

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

    let rowCountCheck = false;

    let notificationMessage: string = "";
    let notificationType: 'success' | 'error' = 'success';

    // Database selection related variables
    let availableDatabases: Array<{UUID: string, name: string, DBType: string, description?: string}> = [];
    let selectedDatabaseUUID: string = "";
    let isLoadingDatabases: boolean = false;

    // Convert ColumnInfo to Column
    function convertColumnInfoToColumn(columnInfo: ColumnInfo): Column {
        return {
            name: columnInfo.name,
            type: columnInfo.type,
            isForeignKey: columnInfo.isForeignKey,
            // Provide empty strings as defaults for optional fields
            referencedTable: columnInfo.referencedTable || '',
            referencedColumn: columnInfo.referencedColumn || ''
        };
    }

    function transformSchemaFormat(complexSchema: Record<string, TableSchema> | { tables: any[] } | null | undefined): Schema {
        if (!complexSchema) return { tables: [] };

        if (Array.isArray((complexSchema as any).schema)) {
            return {
                tables: (complexSchema as any).schema.map((table: any) => ({
                    name: table.name,
                    columns: table.columns.map(convertColumnInfoToColumn),
                    primaryKey: table.primaryKey || ''
                }))
            };
        }

        if (Array.isArray((complexSchema as any).tables)) {
            return {
                tables: (complexSchema as any).tables.map((table: any) => ({
                    name: table.name,
                    columns: table.columns.map(convertColumnInfoToColumn),
                    primaryKey: table.primaryKey || ''
                }))
            };
        }

        return {
            tables: Object.entries(complexSchema).map(([tableId, tableSchema]) => ({
                name: tableSchema.name,
                columns: tableSchema.columns.map(convertColumnInfoToColumn),
                primaryKey: tableSchema.primaryKey || ''
            }))
        };
    }

    let schema = transformSchemaFormat(data.schema);

    async function loadUserDatabases() {
        try {
            isLoadingDatabases = true;
            // Add the built-in database as the first option
            availableDatabases = [{
                UUID: "", // Empty UUID indicates the built-in database
                name: "Built-in Database",
                DBType: "sqlite",
                description: "Default SQLite database"
            }];

            // Fetch user's custom databases
            const response = await listCustomDatabases();
            if (!response) {
                throw new Error('Failed to fetch databases');
            }

            const customDatabases = await response;
            if (Array.isArray(customDatabases)) {
                availableDatabases = [...availableDatabases, ...customDatabases];
            }

            // Select the built-in database by default
            selectedDatabaseUUID = "";

            notificationMessage = "Databases loaded successfully";
            notificationType = 'success';
        } catch (error) {
            console.error("Error loading databases:", error);
            errorMessage = error instanceof Error ? error.message : 'Failed to load databases';

            notificationMessage = "Failed to load databases";
            notificationType = 'error';
        } finally {
            isLoadingDatabases = false;
        }
    }

    // Initialize by loading databases
    loadUserDatabases();

    // Function to handle database change
    // Fixed function to handle database change
    async function handleDatabaseChange() {
        if (stage !== 'input') {
            resetQuery();
        }

        try {
            if (selectedDatabaseUUID) {
                const customSchema = await getCustomSchemaComplex(selectedDatabaseUUID);
                schema = transformSchemaFormat(customSchema);
            } else {
                // Use the built-in database schema
                schema = transformSchemaFormat(data.schema);
            }
        } catch (error) {
            console.error("Error loading schema for database:", error);
            errorMessage = error instanceof Error ? error.message : 'Failed to load database schema';

            notificationMessage = "Failed to load database schema";
            notificationType = 'error';
        }
    }

    async function processNaturalLanguageQuery() {
        try {
            isLoading = true;
            errorMessage = null;

            if (!schema || !schema.tables || schema.tables.length === 0) {
                errorMessage = "No database schema available. Please select a valid database.";
                notificationMessage = "No schema available";
                notificationType = 'error';
                return;
            }

            console.log("Using schema for query:", schema);
            const response: TextToSQLRequest = await generateComplexSQLbyGPT(naturalLanguageQuery, schema);
            console.log("Generated SQL response:", response);
            sqlQuery = response.sql;

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

            // Pass the selected database UUID to the execution function
            queryResult = await executeQueryVisualisation(sqlQuery, selectedDatabaseUUID);

            if (sqlQuery.trim().toUpperCase().startsWith('SELECT')) {
                svgContent = queryResult.svg;

                notificationMessage = "Query executed and visualization generated";
                notificationType = 'success';
            } else {
                svgContent = null;

                notificationMessage = "Query executed successfully";
                notificationType = 'success';
            }

            stage = 'result';

            console.log(queryResult)

            rowCountCheck = queryResult && queryResult.row_count > 0
        } catch (error) {
            console.error("Error executing SQL or generating visualization:", error);
            errorMessage = error instanceof Error ? error.message : 'An unknown error occurred';

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

        notificationMessage = "Started a new query";
        notificationType = 'success';
    }

    function useExample(example: string) {
        naturalLanguageQuery = example;

        notificationMessage = "Example selected";
        notificationType = 'success';
    }

    function downloadSVGAsImage() {
        const svgElement = new DOMParser().parseFromString(svgContent, 'image/svg+xml').documentElement;

        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');

        const img = new Image();
        const svgData = new XMLSerializer().serializeToString(svgElement);
        const svgBlob = new Blob([svgData], { type: 'image/svg+xml' });
        const svgUrl = URL.createObjectURL(svgBlob);

        img.onload = function () {
            canvas.width = img.width;
            canvas.height = img.height;

            ctx.drawImage(img, 0, 0);

            const link = document.createElement('a');
            link.href = canvas.toDataURL('image/jpeg');
            link.download = 'visualization.jpg';
            link.click();
        };

        img.src = svgUrl;
    }
</script>

<svelte:head>
    <title>Visualisation of query</title>
</svelte:head>

<Notification message={notificationMessage} type={notificationType} />

<div class="container mx-auto p-4 max-w-4xl">
    {#if schema}
        <h1 class="text-3xl font-bold mb-6 text-gray-800">Data Visualization</h1>

        {#if errorMessage}
            <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-6 rounded">
                <div class="flex">
                    <AlertCircle class="h-6 w-6 text-red-500 mr-2" />
                    <p class="text-red-700">{errorMessage}</p>
                </div>
            </div>
        {/if}

        <!-- Database Selection -->
        <div class="mb-6">
            <div class="flex items-center mb-2">
                <Database class="w-5 h-5 text-blue-600 mr-2" />
                <h2 class="text-lg font-semibold text-gray-700">Select Database</h2>
            </div>

            <div class="bg-white p-4 rounded-lg shadow-md">
                {#if isLoadingDatabases}
                    <div class="flex items-center justify-center p-4">
                        <Spinner size="md" color="primary" variant="border" />
                        <span class="ml-2 text-gray-600">Loading databases...</span>
                    </div>
                {:else}
                    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4">
                        <div class="w-full sm:w-2/3">
                            <select
                                    bind:value={selectedDatabaseUUID}
                                    on:change={handleDatabaseChange}
                                    class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-400 focus:border-transparent"
                            >
                                {#each availableDatabases as db}
                                    <option value={db.UUID}>
                                        {db.name} ({db.DBType})
                                    </option>
                                {/each}
                            </select>
                        </div>
                        <button
                                class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
                                on:click={loadUserDatabases}
                        >
                            <RefreshCw class="w-4 h-4 inline mr-1" />
                            Refresh
                        </button>
                    </div>

                    {#if selectedDatabaseUUID !== ""}
                        <div class="mt-2 text-sm text-gray-500">
                            {#if availableDatabases.find(db => db.UUID === selectedDatabaseUUID)?.description}
                                {availableDatabases.find(db => db.UUID === selectedDatabaseUUID)?.description}
                            {:else}
                                Using custom database: {availableDatabases.find(db => db.UUID === selectedDatabaseUUID)?.name}
                            {/if}
                        </div>
                    {:else}
                        <div class="mt-2 text-sm text-gray-500">
                            Using built-in SQLite database
                        </div>
                    {/if}
                {/if}
            </div>
        </div>

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
                                <div class="mt-4 flex justify-end">
                                    <button
                                            class="px-4 py-2 bg-blue-600 text-white rounded-lg flex items-center"
                                            on:click={downloadSVGAsImage}
                                    >
                                        Download as Image
                                    </button>
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

                        {#if rowCountCheck}
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
                                            {#each row.values as cell}
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
    {:else}
        <Spinner size="xl" color="primary" variant="border" />
    {/if}
</div>

<Footer />
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