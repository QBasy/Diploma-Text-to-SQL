<script lang="ts">
    type DataType = 'INTEGER' | 'TEXT' | 'BOOLEAN' | 'DATE' | 'UUID' | 'FLOAT';
    type ConstraintType = 'PRIMARY KEY' | 'FOREIGN KEY';

    interface Column {
        name: string;
        type: DataType;
        constraints: ConstraintType[];
        foreignKey?: {
            table: string;
            column: string;
        };
    }

    interface Table {
        name: string;
        columns: Column[];
        primaryKey: string;
    }

    let showFKeyModal: boolean = false;
    let referencedTable: string = '';
    let referencedColumn: string = '';

    let schema: Table[] = [];
    let selectedTableIndex = -1;
    let foreignKeyState = { tableIndex: -1, columnIndex: -1, refTable: '', refColumn: '' };

    $: selectedTable = selectedTableIndex >= 0 ? schema[selectedTableIndex] : null;

    const dataTypes = [
        { value: 'INTEGER', label: 'Integer' },
        { value: 'TEXT', label: 'Text' },
        { value: 'BOOLEAN', label: 'Boolean' },
        { value: 'DATE', label: 'Date' },
        { value: 'UUID', label: 'UUID' },
        { value: 'FLOAT', label: 'Float' }
    ];

    function addTable() {
        schema = [...schema, { name: `New Table ${schema.length + 1}`, columns: [], primaryKey: '' }];
        selectedTableIndex = schema.length - 1;
    }

    function removeTable(index: number) {
        schema = schema.filter((_, i) => i !== index);
        if (selectedTableIndex === index) selectedTableIndex = -1;
    }

    function addColumn(tableIndex: number) {
        schema = schema.map((table, i) => {
            if (i === tableIndex) {
                return {
                    ...table,
                    columns: [...table.columns, { name: `Column ${table.columns.length + 1}`, type: 'TEXT', constraints: [] }]
                };
            }
            return table;
        });
    }

    function setPrimaryKey(tableIndex: number, columnIndex: number) {
        schema = schema.map((table, i) => {
            if (i === tableIndex) {
                return {
                    ...table,
                    primaryKey: table.columns[columnIndex].name,
                    columns: table.columns.map((col, idx) => ({
                        ...col,
                        constraints: idx === columnIndex ? ['PRIMARY KEY'] : col.constraints.filter(c => c !== 'PRIMARY KEY')
                    }))
                };
            }
            return table;
        });
    }

    function handleForeignKey() {
        if (!foreignKeyState.refTable || !foreignKeyState.refColumn) return;

        schema = schema.map((table, i) => {
            if (i === foreignKeyState.tableIndex) {
                const columns = [...table.columns];
                columns[foreignKeyState.columnIndex] = {
                    ...columns[foreignKeyState.columnIndex],
                    constraints: [...columns[foreignKeyState.columnIndex].constraints, 'FOREIGN KEY'],
                    foreignKey: { table: foreignKeyState.refTable, column: foreignKeyState.refColumn }
                };
                return { ...table, columns };
            }
            return table;
        });

        showFKeyModal = false;
        foreignKeyState = { tableIndex: -1, columnIndex: -1, refTable: '', refColumn: '' };
    }
</script>

<section class="container w-10/12 mx-auto">
    <div class="grid grid-cols-4 h-screen p-4 gap-4 bg-gray-100">
        <div class="col-span-1 bg-white p-4 rounded-lg shadow-md">
            <button on:click={addTable} class="w-full bg-blue-500 text-white py-2 px-4 rounded-md mb-4">+ New Table</button>
            <div class="space-y-2 grid grid-cols-1">
                {#each schema as table, i}
                    <div class="flex justify-between items-center p-2 rounded-md cursor-pointer hover:bg-gray-200">
                        <button on:click={() => selectedTableIndex = i} class="flex-1 text-left">{table.name}</button>
                        <button on:click={() => removeTable(i)} class="text-red-500">×</button>
                    </div>
                {/each}
            </div>
        </div>

        {#if selectedTable}
            <div class="col-span-3 bg-white p-4 rounded-lg shadow-md">
                <input type="text" bind:value={selectedTable.name} class="text-xl font-bold mb-4 border p-2 rounded-md w-full" />
                <div class="space-y-2">
                    {#each selectedTable.columns as column, colIndex}
                        <div class="flex gap-2">
                            <input type="text" bind:value={column.name} class="border p-2 rounded-md w-1/3" />
                            <select bind:value={column.type} class="border p-2 rounded-md w-1/3">
                                {#each dataTypes as type}
                                    <option value={type.value}>{type.label}</option>
                                {/each}
                            </select>
                            <button on:click={() => setPrimaryKey(selectedTableIndex, colIndex)} class="bg-green-500 text-white px-2 py-1 rounded-md">
                                {column.constraints.includes('PRIMARY KEY') ? 'PK ✓' : 'Set PK'}
                            </button>
                            <button on:click={() => { foreignKeyState = { tableIndex: selectedTableIndex, columnIndex: colIndex }; showFKeyModal = true; }} class="bg-red-500 text-white px-2 py-1 rounded-md">
                                {column.constraints.includes('FOREIGN KEY') ? 'FK ✓' : 'Set FK'}
                            </button>
                        </div>
                    {/each}
                    <button on:click={() => addColumn(selectedTableIndex)} class="mt-2 bg-gray-500 text-white px-4 py-2 rounded-md">+ Add Column</button>
                </div>
            </div>
        {/if}
    </div>
</section>


{#if showFKeyModal}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
            <h3 class="text-lg font-semibold mb-4">Select Foreign Key</h3>

            <label class="block text-sm font-medium text-gray-700">Referenced Table</label>
            <select
                    bind:value={referencedTable}
                    class="w-full p-2 border rounded-lg shadow-sm"
                    on:change={() => referencedColumn = ''}
            >
                <option value="">Select Table</option>
                {#each schema as table}
                    <option value={table.name}>{table.name}</option>
                {/each}
            </select>

            {#if referencedTable}
                <label class="block text-sm font-medium text-gray-700">Referenced Column</label>
                <select
                        bind:value={referencedColumn}
                        class="w-full p-2 border rounded-lg shadow-sm"
                >
                    <option value="">Select Column</option>
                    {#each schema.find(t => t.name === referencedTable)?.columns || [] as column}
                        <option value={column.name}>{column.name}</option>
                    {/each}
                </select>
            {/if}

            <button
                    on:click={() => {
                    if (referencedTable && referencedColumn) {
                        foreignKeyState = {
                            ...foreignKeyState,
                            refTable: referencedTable,
                            refColumn: referencedColumn
                        };
                        handleForeignKey();
                    }
                }}
                    class="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
                Add Foreign Key
            </button>
            <button
                    on:click={() => showFKeyModal = false}
                    class="mt-2 w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
            >
                Cancel
            </button>
        </div>
    </div>
{/if}