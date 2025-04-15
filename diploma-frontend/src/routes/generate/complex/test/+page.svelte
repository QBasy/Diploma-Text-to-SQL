<script lang="ts">
    import { generateComplexSQL } from '$lib/api';
    import type { Table, Column } from '$lib/types/table';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

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

    let showTableModal: boolean = false;
    let editingTable: Table | null = null;

    const openTableModal = (tableIndex: number | null = null) => {
        if (tableIndex !== null) {
            editingTable = { ...schema.tables[tableIndex] };
        } else {
            editingTable = { name: '', columns: [], primaryKey: '' };
        }
        selectedTableIndex = tableIndex;
        showTableModal = true;
    };

    const handleSaveTable = () => {
        if (editingTable) {
            if (selectedTableIndex !== null) {
                schema.tables[selectedTableIndex] = editingTable;
            } else {
                schema.tables.push(editingTable);
            }
            showTableModal = false;
            editingTable = null;
        }
    };

    const handleCancelTable = () => {
        showTableModal = false;
        editingTable = null;
    };
</script>

<svelte:head>
    <title>Complex SQL Query Generator</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 py-8">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <header class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-900">Complex SQL Query Generator</h1>
            <p class="mt-2 text-sm text-gray-600">Define your schema and generate complex SQL queries.</p>
        </header>

        <form on:submit|preventDefault={submitQuery} class="bg-white shadow-lg rounded-lg p-6">
            <div class="space-y-6">
                <div>
                    <label for="query" class="block text-sm font-medium text-gray-700">SQL Query:</label>
                    <textarea
                            id="query"
                            bind:value={query}
                            rows="4"
                            class="mt-1 block w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
                            placeholder="Enter your SQL query here..."
                    ></textarea>
                </div>

                <div>
                    <h2 class="text-xl font-semibold text-gray-900 mb-4">Schema Definition</h2>
                    <div class="space-y-4">
                        {#each schema.tables as table, tableIndex (tableIndex)}
                            <div class="border border-gray-200 rounded-lg p-4">
                                <div class="mb-4">
                                    <label class="block text-sm font-medium text-gray-700">Table Name</label>
                                    <input
                                            type="text"
                                            bind:value={table.name}
                                            on:input={(e) => updateTableName(tableIndex, e.currentTarget.value)}
                                            class="mt-1 block w-full p-2 border border-gray-300 rounded-md"
                                            placeholder="Table Name"
                                    />
                                </div>
                                <div class="mb-4">
                                    <label class="block text-sm font-medium text-gray-700">Primary Key</label>
                                    {#if table.primaryKey}
                                        <div class="mt-1 p-2 bg-gray-100 rounded">{table.primaryKey}</div>
                                    {:else}
                                        <button
                                                type="button"
                                                on:click={() => openPKeyModal(tableIndex)}
                                                class="w-full mt-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
                                        >
                                            Add Primary Key
                                        </button>
                                    {/if}
                                </div>
                                <div class="space-y-2">
                                    {#each table.columns as column, columnIndex (columnIndex)}
                                        <div class="border border-gray-100 rounded-md p-3">
                                            <div class="mb-2">
                                                <label class="block text-sm font-medium text-gray-700">Column Name</label>
                                                <input
                                                        type="text"
                                                        bind:value={column.name}
                                                        on:input={(e) => updateColumn(tableIndex, columnIndex, { name: e.currentTarget.value })}
                                                        class="mt-1 block w-full p-2 border border-gray-300 rounded-md"
                                                        placeholder="Column Name"
                                                />
                                            </div>
                                            <div class="mb-2">
                                                <label class="block text-sm font-medium text-gray-700">Data Type</label>
                                                <select
                                                        bind:value={column.type}
                                                        on:change={(e) => updateColumn(tableIndex, columnIndex, { type: e.currentTarget.value })}
                                                        class="mt-1 block w-full p-2 border border-gray-300 rounded-md"
                                                >
                                                    <option value="">Select Type</option>
                                                    <option value="INT">INT</option>
                                                    <option value="VARCHAR">VARCHAR</option>
                                                    <option value="DATE">DATE</option>
                                                    <option value="BOOLEAN">BOOLEAN</option>
                                                </select>
                                            </div>
                                            <div>
                                                {#if column.isForeignKey}
                                                    <div class="text-sm">
                                                        <div class="text-gray-700">References {column.referencedTable}.{column.referencedColumn}</div>
                                                        <div class="mt-2 space-x-2">
                                                            <button on:click={() => openFKeyModal(tableIndex, columnIndex)} class="text-blue-600 hover:text-blue-800 text-sm">Edit</button>
                                                            <button on:click={() => removeForeignKey(tableIndex, columnIndex)} class="text-red-600 hover:text-red-800 text-sm">Remove</button>
                                                        </div>
                                                    </div>
                                                {:else}
                                                    <button type="button" on:click={() => openFKeyModal(tableIndex, columnIndex)} class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm">Add Foreign Key</button>
                                                {/if}
                                            </div>
                                        </div>
                                    {/each}
                                    <button type="button" on:click={() => addColumn(tableIndex)} class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200">+ Add Column</button>
                                </div>
                            </div>
                        {/each}
                        <button type="button" on:click={addTable} class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200">+ Add Table</button>
                    </div>
                </div>

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

        {#if sqlResult}
            <div class="mt-8 bg-white shadow-lg rounded-lg p-6">
                <h2 class="text-xl font-semibold text-gray-900">Generated SQL</h2>
                <pre class="mt-4 p-4 bg-gray-100 rounded-md text-sm font-mono text-gray-700 overflow-x-auto">{sqlResult}</pre>
            </div>
        {/if}
    </div>

    {#if showPKeyModal && selectedTableIndex !== null}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">Select Primary Key for {schema.tables[selectedTableIndex].name}</h3>
                <div class="space-y-2">
                    {#each schema.tables[selectedTableIndex].columns as column}
                        <button on:click={() => handleSetPrimaryKey(column.name)} class="w-full p-3 text-left rounded-md bg-gray-50 hover:bg-gray-100 transition-colors">{column.name}</button>
                    {/each}
                </div>
                <div class="mt-6 flex justify-end">
                    <button on:click={() => showPKeyModal = false} class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md">Cancel</button>
                </div>
            </div>
        </div>
    {/if}

    {#if showFKeyModal && selectedTableIndex !== null && selectedColumnIndex !== null}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">Set Foreign Key for {schema.tables[selectedTableIndex].columns[selectedColumnIndex].name}</h3>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Referenced Table</label>
                        <select bind:value={selectedRefTable} class="w-full p-2 border border-gray-300 rounded-md">
                            <option value="">Select a table</option>
                            {#each getAvailableTables(schema.tables[selectedTableIndex].name) as table}
                                <option value={table.name}>{table.name}</option>
                            {/each}
                        </select>
                    </div>
                    {#if selectedRefTable}
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">Referenced Column</label>
                            <select bind:value={selectedRefColumn} class="w-full p-2 border border-gray-300 rounded-md">
                                <option value="">Select a column</option>
                                {#each getAvailableColumns(selectedRefTable) as column}
                                    <option value={column.name}>{column.name}</option>
                                {/each}
                            </select>
                        </div>
                    {/if}
                </div>
                <div class="mt-6 flex justify-end space-x-3">
                    <button on:click={() => showFKeyModal = false} class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md">Cancel</button>
                    <button on:click={handleSetForeignKey} disabled={!selectedRefTable || !selectedRefColumn} class="px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50 disabled:cursor-not-allowed">Confirm</button>
                </div>
            </div>
        </div>
    {/if}

    <Notification message={notificationMessage} type={notificationType} />
</div>

<div class="min-h-screen bg-gray-100 py-8 flex">
    <aside class="w-64 bg-gray-200 p-4 space-y-2">
        <h2 class="text-lg font-semibold mb-4">Таблицы</h2>
        <button on:click={() => openTableModal()} class="w-full p-2 bg-blue-500 text-white rounded">+ Добавить таблицу</button>
        {#each schema.tables as table, index (index)}
            <button on:click={() => selectedTableIndex = index} class="w-full p-2 rounded hover:bg-gray-300" class:bg-gray-300={selectedTableIndex === index}>
                {table.name}
            </button>
        {/each}
    </aside>

    <div class="flex-1 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        {#if selectedTableIndex !== null}
            <div class="mt-8 bg-white shadow-lg rounded-lg p-6">
                <h2 class="text-xl font-semibold text-gray-900 mb-4">{schema.tables[selectedTableIndex].name}</h2>
                <button on:click={() => openTableModal(selectedTableIndex)} class="mt-4 px-4 py-2 bg-blue-500 text-white rounded">Редактировать таблицу</button>
            </div>
        {/if}
    </div>

    {#if showTableModal && editingTable}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">Редактировать таблицу</h3>
                <div class="mt-6 flex justify-end space-x-3">
                    <button on:click={handleCancelTable} class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md">Отмена</button>
                    <button on:click={handleSaveTable} class="px-4 py-2 bg-blue-600 text-white rounded-md">Сохранить</button>
                </div>
            </div>
        </div>
    {/if}
</div>
