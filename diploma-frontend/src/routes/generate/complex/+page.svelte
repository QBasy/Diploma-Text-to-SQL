<script lang="ts">
    import { generateComplexSQL } from '$lib/api';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

    let isLoading: boolean = false;
    let notificationMessage: string = '';
    let notificationType: 'success' | 'error' = 'success';
    let query: string = '';
    let schema: any = { tables: [] };
    let sqlResult: string = '';

    const addTable = () => {
        schema.tables = [...schema.tables, { name: '', columns: [] }];
    };

    const addColumn = (tableIndex: number) => {
        schema.tables[tableIndex].columns = [...schema.tables[tableIndex].columns, { name: '' }];
    };

    const submitQuery = async () => {
        isLoading = true;
        try {
            const response = await generateComplexSQL(query, schema);
            sqlResult = response.sql;
            notificationMessage = 'SQL generated successfully!';
            notificationType = 'success';
        } catch (err: any) {
            notificationMessage = `Error: ${err.message}`;
            notificationType = 'error';
        } finally {
            isLoading = false;
        }
    };
</script>

<svelte:head>
    <title>Complex Query</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header -->
        <header class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-900">Complex SQL Query Generator</h1>
            <p class="mt-2 text-sm text-gray-600">Enter your query and define the schema to generate SQL.</p>
        </header>

        <!-- Query Input Form -->
        <form on:submit|preventDefault={submitQuery} class="bg-white shadow-sm rounded-lg p-6">
            <div class="space-y-6">
                <!-- Query Textarea -->
                <div>
                    <label for="query" class="block text-sm font-medium text-gray-700">Enter your query:</label>
                    <textarea
                            id="query"
                            bind:value={query}
                            rows="4"
                            class="mt-1 block w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
                            placeholder="e.g., Select all users from the database"
                    ></textarea>
                </div>

                <!-- Schema Section -->
                <div>
                    <h2 class="text-xl font-semibold text-gray-900 mb-4">Schema</h2>
                    {#each schema.tables as table, tableIndex}
                        <div class="mb-6 p-4 border border-gray-200 rounded-lg">
                            <!-- Table Name Input -->
                            <div class="mb-4">
                                <label class="block text-sm font-medium text-gray-700">Table Name:</label>
                                <input
                                        type="text"
                                        bind:value={table.name}
                                        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="e.g., users"
                                />
                            </div>

                            <!-- Columns Section -->
                            <div class="space-y-4">
                                {#each table.columns as column, columnIndex}
                                    <div>
                                        <label class="block text-sm font-medium text-gray-700">Column Name:</label>
                                        <input
                                                type="text"
                                                bind:value={column.name}
                                                class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                                placeholder="e.g., id, name"
                                        />
                                    </div>
                                {/each}

                                <!-- Add Column Button -->
                                <button
                                        type="button"
                                        on:click={() => addColumn(tableIndex)}
                                        class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
                                >
                                    + Add Column
                                </button>
                            </div>
                        </div>
                    {/each}

                    <!-- Add Table Button -->
                    <button
                            type="button"
                            on:click={addTable}
                            class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
                    >
                        + Add Table
                    </button>
                </div>

                <!-- Submit Button -->
                <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full flex justify-center items-center px-4 py-2 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {#if isLoading}
                        <LoadingSpinner size={6} color="white" />
                        <span class="ml-2">Generating...</span>
                    {:else}
                        Generate SQL
                    {/if}
                </button>
            </div>
        </form>

        <!-- Generated SQL Result -->
        {#if sqlResult}
            <div class="mt-8 bg-white shadow-sm rounded-lg p-6">
                <h2 class="text-xl font-semibold text-gray-900">Generated SQL:</h2>
                <pre class="mt-4 p-4 bg-gray-100 rounded-md text-sm font-mono text-gray-700 overflow-x-auto">{sqlResult}</pre>
            </div>
        {/if}
    </div>

    <!-- Notification -->
    <Notification {notificationMessage} {notificationType}/>
</div>
