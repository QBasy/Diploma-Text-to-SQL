<script lang="ts">
    import { generateComplexSQL } from '$lib/api';
    import type { Table, Column } from '$lib/types/table';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import { PlusCircle, Edit, Trash2 } from 'lucide-svelte'; // Added Lucide icons

    let isLoading: boolean = false;
    let notificationMessage: string = '';
    let notificationType: 'success' | 'error' = 'success';
    let query: string = '';
    let sqlResult: string = '';

    let customDataType: string = '';

    let showPKeyModal: boolean = false;
    let showFKeyModal: boolean = false;
    let selectedTableIndex: number | null = null;
    let selectedColumnIndex: number | null = null;
    let selectedRefTable: string = '';
    let selectedRefColumn: string = '';

    let schema: { tables: Table[] } = { tables: [] };

    $: schema.tables.forEach((table, index) => {
        if (table.primaryKey && !table.columns.some(c => c.name === table.primaryKey)) {
            updatePrimaryKey(index, '');
        }
    });

    const addTable = (): void => {
        schema = {
            tables: [...schema.tables, { name: '', columns: [], primaryKey: '' }]
        };
    };

    const addColumn = (tableIndex: number): void => {
        const newColumn: Column = {
            name: '',
            type: '',
            isForeignKey: false,
            referencedTable: '',
            referencedColumn: ''
        };

        const updatedTables = [...schema.tables];
        updatedTables[tableIndex].columns = [...updatedTables[tableIndex].columns, newColumn];
        schema = { tables: updatedTables };
    };

    const updateColumn = (tableIndex: number, columnIndex: number, updates: Partial<Column>): void => {
        const updatedTables = [...schema.tables];
        updatedTables[tableIndex].columns[columnIndex] = {
            ...updatedTables[tableIndex].columns[columnIndex],
            ...updates
        };
        schema = { tables: updatedTables };
    };

    const updateTableName = (tableIndex: number, name: string): void => {
        const updatedTables = [...schema.tables];
        updatedTables[tableIndex].name = name;
        schema = { tables: updatedTables };
    };

    const updatePrimaryKey = (tableIndex: number, primaryKey: string): void => {
        const updatedTables = [...schema.tables];
        updatedTables[tableIndex].primaryKey = primaryKey;
        schema = { tables: updatedTables };
    };

    const openPKeyModal = (tableIndex: number) => {
        selectedTableIndex = tableIndex;
        showPKeyModal = true;
    };

    const handleSetPrimaryKey = (columnName: string) => {
        if (selectedTableIndex !== null) {
            updatePrimaryKey(selectedTableIndex, columnName);
            showPKeyModal = false;
        }
    };

    // Foreign Key Modal handlers
    const openFKeyModal = (tableIndex: number, columnIndex: number) => {
        selectedTableIndex = tableIndex;
        selectedColumnIndex = columnIndex;
        const column = schema.tables[tableIndex].columns[columnIndex];
        selectedRefTable = column.referencedTable;
        selectedRefColumn = column.referencedColumn;
        showFKeyModal = true;
    };

    const handleSetForeignKey = () => {
        if (selectedTableIndex === null || selectedColumnIndex === null) return;

        updateColumn(selectedTableIndex, selectedColumnIndex, {
            isForeignKey: true,
            referencedTable: selectedRefTable,
            referencedColumn: selectedRefColumn
        });

        showFKeyModal = false;
        selectedRefTable = '';
        selectedRefColumn = '';
    };

    const removeForeignKey = (tableIndex: number, columnIndex: number) => {
        updateColumn(tableIndex, columnIndex, {
            isForeignKey: false,
            referencedTable: '',
            referencedColumn: ''
        });
    };

    const submitQuery = async (): Promise<void> => {
        isLoading = true;
        try {
            console.log(schema, query)
            const response = await generateComplexSQL(query, schema);
            sqlResult = response.sql;
            console.log(sqlResult);
            notificationMessage = 'SQL generated successfully!';
            notificationType = 'success';
        } catch (err: unknown) {
            notificationMessage = `Error: ${(err as Error).message}`;
            notificationType = 'error';
        } finally {
            isLoading = false;
        }
    };

    const getAvailableColumns = (tableName: string): Column[] => {
        return schema.tables.find(t => t.name === tableName)?.columns || [];
    };

    const getAvailableTables = (currentTableName: string): Table[] => {
        return schema.tables.filter(t => t.name && t.name !== currentTableName);
    };
</script>

<svelte:head>
    <title>Complex SQL Query</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <header class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-900">Complex SQL Query Generator</h1>
            <p class="mt-2 text-sm text-gray-600">Enter your query and define the schema with primary keys, foreign keys, and data types.</p>
        </header>

        <form on:submit|preventDefault={submitQuery} class="bg-white shadow-sm rounded-lg p-6">
            <div class="space-y-6">
                <!-- Query Input -->
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
                    <div class="grid grid-cols-3 gap-4">
                        {#each schema.tables as table, tableIndex}
                            <div class="p-4 border border-gray-200 rounded-lg">
                                <!-- Table Name -->
                                <div class="mb-4">
                                    <label class="block text-sm font-medium text-gray-700">Table Name:</label>
                                    <input
                                            type="text"
                                            value={table.name}
                                            on:input={(e) => updateTableName(tableIndex, e.currentTarget.value)}
                                            class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                            placeholder="e.g., users"
                                    />
                                </div>

                                <!-- Primary Key -->
                                <div class="mb-4">
                                    <label class="block text-sm font-medium text-gray-700">Primary Key:</label>
                                    {#if table.primaryKey}
                                        <div class="mt-1 p-2 bg-gray-100 rounded">{table.primaryKey}</div>
                                    {:else}
                                        <button
                                                type="button"
                                                on:click={() => openPKeyModal(tableIndex)}
                                                class="w-full mt-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 flex items-center justify-between"
                                        >
                                            <span class="text-sm">+ Add Primary Key</span>
                                            <PlusCircle class="h-5 w-5 text-gray-700" />
                                        </button>
                                    {/if}
                                </div>

                                <!-- Columns -->
                                <div class="space-y-4">
                                    {#each table.columns as column, columnIndex}
                                        <div class="p-4 border border-gray-100 rounded-md">
                                            <!-- Column Name -->
                                            <div class="mb-4">
                                                <label class="block text-sm font-medium text-gray-700">Column Name:</label>
                                                <input
                                                        type="text"
                                                        value={column.name}
                                                        on:input={(e) => updateColumn(tableIndex, columnIndex, { name: e.currentTarget.value })}
                                                        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                                        placeholder="e.g., id"
                                                />
                                            </div>

                                            <!-- Data Type -->
                                            <div class="mb-4">
                                                <label class="block text-sm font-medium text-gray-700">Data Type:</label>
                                                <select
                                                        value={column.type}
                                                        on:change={(e) => updateColumn(tableIndex, columnIndex, { type: e.currentTarget.value })}
                                                        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                                >
                                                    <option value="">Select type</option>
                                                    <option value="INT">INT</option>
                                                    <option value="VARCHAR">VARCHAR</option>
                                                    <option value="DATE">DATE</option>
                                                    <option value="BOOLEAN">BOOLEAN</option>
                                                </select>
                                            </div>

                                            <!-- Foreign Key -->
                                            <div class="mb-4">
                                                {#if column.isForeignKey}
                                                    <div class="text-sm">
                                                        <div class="text-gray-700">References {column.referencedTable}.{column.referencedColumn}</div>
                                                        <div class="mt-2 space-x-2">
                                                            <button
                                                                    type="button"
                                                                    on:click={() => openFKeyModal(tableIndex, columnIndex)}
                                                                    class="text-blue-600 hover:text-blue-800 text-sm flex items-center"
                                                            >
                                                                <Edit class="h-4 w-4 mr-1" />
                                                                Edit
                                                            </button>
                                                            <button
                                                                    type="button"
                                                                    on:click={() => removeForeignKey(tableIndex, columnIndex)}
                                                                    class="text-red-600 hover:text-red-800 text-sm flex items-center"
                                                            >
                                                                <Trash2 class="h-4 w-4 mr-1" />
                                                                Remove
                                                            </button>
                                                        </div>
                                                    </div>
                                                {:else}
                                                    <button
                                                            type="button"
                                                            on:click={() => openFKeyModal(tableIndex, columnIndex)}
                                                            class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm flex items-center justify-between"
                                                    >
                                                        <span class="text-sm">+ Add Foreign Key</span>
                                                        <Edit class="h-5 w-5 text-gray-700" />
                                                    </button>
                                                {/if}
                                            </div>
                                        </div>
                                    {/each}
                                </div>

                                <!-- Add Column Button -->
                                <button
                                        type="button"
                                        on:click={() => addColumn(tableIndex)}
                                        class="w-full mt-4 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 flex items-center justify-between"
                                >
                                    <span class="text-sm">+ Add Column</span>
                                    <PlusCircle class="h-5 w-5" />
                                </button>
                            </div>
                        {/each}
                    </div>

                    <!-- Add Table Button -->
                    <button
                            type="button"
                            on:click={addTable}
                            class="w-full mt-8 px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 flex items-center justify-between"
                    >
                        <span class="text-sm">+ Add Table</span>
                        <PlusCircle class="h-5 w-5" />
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

        <!-- SQL Result -->
        {#if sqlResult}
            <div class="mt-8 bg-white shadow-sm rounded-lg p-6">
                <h2 class="text-xl font-semibold text-gray-900">Generated SQL:</h2>
                <pre class="mt-4 p-4 bg-gray-100 rounded-md text-sm font-mono text-gray-700 overflow-x-auto">{sqlResult}</pre>
            </div>
        {/if}
    </div>

    <!-- Modals -->
    {#if showPKeyModal && selectedTableIndex !== null}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">
                    Select Primary Key for {schema.tables[selectedTableIndex].name}
                </h3>
                <div class="space-y-2">
                    {#each schema.tables[selectedTableIndex].columns as column}
                        <button
                                on:click={() => handleSetPrimaryKey(column.name)}
                                class="w-full p-3 text-left rounded-md bg-gray-50 hover:bg-gray-100 transition-colors"
                        >
                            {column.name}
                        </button>
                    {/each}
                </div>
                <div class="mt-6 flex justify-end">
                    <button
                            on:click={() => showPKeyModal = false}
                            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    {/if}

    {#if showFKeyModal && selectedTableIndex !== null && selectedColumnIndex !== null}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">
                    Set Foreign Key for {schema.tables[selectedTableIndex].columns[selectedColumnIndex].name}
                </h3>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Referenced Table</label>
                        <select
                                bind:value={selectedRefTable}
                                class="w-full p-2 border border-gray-300 rounded-md"
                        >
                            <option value="">Select a table</option>
                            {#each getAvailableTables(schema.tables[selectedTableIndex].name) as table}
                                <option value={table.name}>{table.name}</option>
                            {/each}
                        </select>
                    </div>

                    {#if selectedRefTable}
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">Referenced Column</label>
                            <select
                                    bind:value={selectedRefColumn}
                                    class="w-full p-2 border border-gray-300 rounded-md"
                            >
                                <option value="">Select a column</option>
                                {#each getAvailableColumns(selectedRefTable) as column}
                                    <option value={column.name}>{column.name}</option>
                                {/each}
                            </select>
                        </div>
                    {/if}
                </div>
                <div class="mt-6 flex justify-end space-x-3">
                    <button
                            on:click={() => showFKeyModal = false}
                            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
                    >
                        Cancel
                    </button>
                    <button
                            on:click={handleSetForeignKey}
                            disabled={!selectedRefTable || !selectedRefColumn}
                            class="px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <Notification message={notificationMessage} type={notificationType} />
</div>