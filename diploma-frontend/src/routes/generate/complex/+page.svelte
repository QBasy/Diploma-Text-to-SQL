<script lang="ts">
    import { generateComplexSQL } from '$lib/api';
    import type { Table, Column } from '$lib/types/table';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import { PlusCircle, Edit, Trash2, UploadCloud } from 'lucide-svelte';
    import Navbar from "$lib/components/Navbar.svelte";

    let isLoading: boolean = false;
    let notificationMessage: string = '';
    let notificationType: 'success' | 'error' = 'success';
    let query: string = '';
    let sqlResult: string = '';
    let schemaUploadText: string = '';
    let showSchemaUploadModal: boolean = false;

    let customDataType: string = '';
    let showCustomDataTypeModal: boolean = false;
    let selectedTableIndex: number | null = null;
    let selectedColumnIndex: number | null = null;

    let showPKeyModal: boolean = false;
    let showFKeyModal: boolean = false;
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

    // Schema Upload Functions
    const openSchemaUploadModal = () => {
        showSchemaUploadModal = true;
    };

    const processSchemaUpload = () => {
        try {
            // First check if it's a table list query from information_schema
            if (schemaUploadText.includes('table_name') && !schemaUploadText.includes('column_name')) {
                processTableListUpload();
            } else if (
                schemaUploadText.includes('column_name') && schemaUploadText.includes('data_type')
            ) {
                processColumnDetailsUpload();
            } else if (
                /^[\w_]+,[\w_]+,.+/m.test(schemaUploadText) // Check for rows like: users,id,bigint
            ) {
                processCsvStyleSchemaUpload();
            } else {
                const tables = parseCustomSchemaFormat(schemaUploadText);
                if (tables.length > 0) {
                    schema = { tables };
                    notificationMessage = `Schema imported successfully with ${tables.length} tables!`;
                    notificationType = 'success';
                } else {
                    throw new Error('Unable to parse schema format');
                }
            }

            showSchemaUploadModal = false;
        } catch (err) {
            notificationMessage = `Error processing schema: ${(err as Error).message}`;
            notificationType = 'error';
        }
    };

    // New function to process the CSV-style schema format provided
    const processCsvStyleSchemaUpload = () => {
        const lines = schemaUploadText.trim().split('\n');

        // Structure to hold tables and columns
        const tablesMap = new Map<string, {
            columns: Map<string, {
                name: string,
                type: string,
                isPrimary: boolean
            }>,
            primaryKeyName: string
        }>();

        // First pass: gather all tables and columns
        for (const line of lines) {
            const parts = line.trim().split(',');
            if (parts.length >= 3) {
                const tableName = parts[0].trim();
                const columnName = parts[1].trim();
                const dataType = parts[2].trim();

                // Get or create table entry
                if (!tablesMap.has(tableName)) {
                    tablesMap.set(tableName, {
                        columns: new Map(),
                        primaryKeyName: ''
                    });
                }

                const table = tablesMap.get(tableName)!;

                // Add column if it doesn't exist
                if (!table.columns.has(columnName)) {
                    const isPrimary = columnName === 'id' ||
                        columnName === tableName + '_id' ||
                        columnName === 'uuid' ||
                        dataType.includes('primary');

                    table.columns.set(columnName, {
                        name: columnName,
                        type: mapDataType(dataType),
                        isPrimary
                    });

                    // Set potential primary key
                    if (isPrimary && !table.primaryKeyName) {
                        table.primaryKeyName = columnName;
                    }
                }
            }
        }

        // Detect foreign keys (columns ending with _id or _uuid that match table names)
        // or column names that match table names exactly
        const foreignKeys: {
            tableName: string;
            columnName: string;
            refTable: string;
            refColumn: string;
        }[] = [];

        tablesMap.forEach((tableData, tableName) => {
            tableData.columns.forEach((column, columnName) => {
                // Skip if column is primary key
                if (column.isPrimary) return;

                // Check for standard foreign key naming patterns
                if (columnName.endsWith('_id') || columnName.endsWith('_uuid')) {
                    const potentialTableName = columnName.replace(/_id$|_uuid$/, '');
                    if (tablesMap.has(potentialTableName)) {
                        const refTable = tablesMap.get(potentialTableName)!;
                        // Find matching primary key in referenced table
                        let refColumn = '';
                        refTable.columns.forEach((col, colName) => {
                            if (col.isPrimary || colName === 'id' || colName === 'uuid') {
                                refColumn = colName;
                            }
                        });

                        if (refColumn) {
                            foreignKeys.push({
                                tableName,
                                columnName,
                                refTable: potentialTableName,
                                refColumn
                            });
                        }
                    }
                }

                // Check for columns named exactly like tables (e.g. "user" column in "posts" table)
                if (tablesMap.has(columnName)) {
                    const refTable = tablesMap.get(columnName)!;
                    // Find matching primary key in referenced table
                    let refColumn = refTable.primaryKeyName;

                    if (refColumn) {
                        foreignKeys.push({
                            tableName,
                            columnName,
                            refTable: columnName,
                            refColumn
                        });
                    }
                }
            });
        });

        // Create tables array for schema
        const tables: Table[] = [];

        tablesMap.forEach((tableData, tableName) => {
            const columns: Column[] = [];

            tableData.columns.forEach((column, columnName) => {
                const foreignKey = foreignKeys.find(fk =>
                    fk.tableName === tableName && fk.columnName === columnName
                );

                columns.push({
                    name: columnName,
                    type: column.type,
                    isForeignKey: !!foreignKey,
                    referencedTable: foreignKey ? foreignKey.refTable : '',
                    referencedColumn: foreignKey ? foreignKey.refColumn : ''
                });
            });

            tables.push({
                name: tableName,
                columns,
                primaryKey: tableData.primaryKeyName
            });
        });

        schema = { tables };
        notificationMessage = `Imported ${tables.length} tables with ${lines.length} columns!`;
        notificationType = 'success';
    };

    const processTableListUpload = () => {
        // Parse output from:
        // SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
        const lines = schemaUploadText.trim().split('\n');
        const tableNames: string[] = [];

        let startParsing = false;
        for (const line of lines) {
            const trimmedLine = line.trim();

            // Skip header lines
            if (trimmedLine === 'table_name' || trimmedLine.includes('----')) {
                startParsing = true;
                continue;
            }

            if (startParsing && trimmedLine && !trimmedLine.includes('rows)')) {
                tableNames.push(trimmedLine);
            }
        }

        if (tableNames.length === 0) {
            throw new Error('No table names found');
        }

        // Create tables with empty columns
        const tables: Table[] = tableNames.map(name => ({
            name,
            columns: [],
            primaryKey: ''
        }));

        schema = { tables };
        notificationMessage = `Imported ${tables.length} tables. Now you need to add columns.`;
        notificationType = 'success';
    };

    const processColumnDetailsUpload = () => {
        // Parse output from:
        // SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_schema = 'public';
        const lines = schemaUploadText.trim().split('\n');
        const columnData: {tableName: string, columnName: string, dataType: string}[] = [];

        // Skip header line if present
        let startIndex = 0;
        if (lines[0].includes('table_name') || lines[0].includes('column_name')) {
            startIndex = 1;
        }

        // Process each line (skipping headers)
        for (let i = startIndex; i < lines.length; i++) {
            const line = lines[i].trim();
            if (!line || line.includes('rows)')) continue;

            // Split by comma
            const parts = line.split(',');
            if (parts.length >= 3) {
                columnData.push({
                    tableName: parts[0].trim(),
                    columnName: parts[1].trim(),
                    dataType: parts[2].trim()
                });
            }
        }

        if (columnData.length === 0) {
            throw new Error('No column data found');
        }

        // Group by table name
        const tableMap = new Map<string, Column[]>();

        for (const col of columnData) {
            if (!tableMap.has(col.tableName)) {
                tableMap.set(col.tableName, []);
            }

            tableMap.get(col.tableName)?.push({
                name: col.columnName,
                type: mapDataType(col.dataType),
                isForeignKey: false,
                referencedTable: '',
                referencedColumn: ''
            });
        }

        // Determine primary keys (usually 'id' or 'uuid')
        const tables: Table[] = Array.from(tableMap.entries()).map(([name, columns]) => {
            let primaryKey = '';

            // First look for 'id' column
            const idColumn = columns.find(c => c.name === 'id');
            if (idColumn) {
                primaryKey = 'id';
            } else {
                // Then look for 'uuid' column
                const uuidColumn = columns.find(c => c.name === 'uuid');
                if (uuidColumn) {
                    primaryKey = 'uuid';
                } else {
                    // Then look for '<table>_id' pattern
                    const tableIdColumn = columns.find(c => c.name === `${name}_id`);
                    if (tableIdColumn) {
                        primaryKey = tableIdColumn.name;
                    }
                }
            }

            return {
                name,
                columns,
                primaryKey
            };
        });

        // Look for potential foreign keys
        tables.forEach(table => {
            table.columns.forEach(column => {
                // Skip if it's the primary key
                if (column.name === table.primaryKey) return;

                // Check if this column name ends with '_id' or '_uuid'
                if (column.name.endsWith('_id') || column.name.endsWith('_uuid')) {
                    const potentialTableName = column.name.replace(/_id$|_uuid$/, '');

                    // Look for a matching table
                    const referencedTable = tables.find(t => t.name === potentialTableName);
                    if (referencedTable && referencedTable.primaryKey) {
                        column.isForeignKey = true;
                        column.referencedTable = referencedTable.name;
                        column.referencedColumn = referencedTable.primaryKey;
                    }
                }
                // Also check for exact table name matches
                else {
                    const referencedTable = tables.find(t => t.name === column.name);
                    if (referencedTable && referencedTable.primaryKey) {
                        column.isForeignKey = true;
                        column.referencedTable = referencedTable.name;
                        column.referencedColumn = referencedTable.primaryKey;
                    }
                }
            });
        });

        schema = { tables };
        notificationMessage = `Imported ${tables.length} tables with ${columnData.length} columns!`;
        notificationType = 'success';
    };

    const mapDataType = (dbType: string): string => {
        // Map database-specific types to our simplified type system
        const type = dbType.toLowerCase();

        if (type.includes('int') || type.includes('serial') || type.includes('numeric') || type.includes('bigint')) {
            return 'INT';
        } else if (type.includes('varchar') || type.includes('text') || type.includes('char') || type.includes('jsonb') || type.includes('json')) {
            return 'VARCHAR';
        } else if (type.includes('date') || type.includes('time')) {
            return 'TIMESTAMP';
        } else if (type.includes('bool')) {
            return 'BOOLEAN';
        } else if (type.includes('uuid')) {
            return 'UUID';
        } else {
            // Return the original type if no mapping is found
            return type.toUpperCase();
        }
    };

    const parseCustomSchemaFormat = (text: string): Table[] => {
        // This is a fallback parser for custom schema formats
        // Could be extended based on common schema description formats
        const tables: Table[] = [];
        const lines = text.trim().split('\n');

        let currentTable: Table | null = null;

        for (const line of lines) {
            const trimmedLine = line.trim();

            // Skip empty lines
            if (!trimmedLine) continue;

            // Check if line defines a table
            if (trimmedLine.match(/^(CREATE\s+TABLE\s+)?["`']?(\w+)["`']?/i)) {
                // Extract table name
                const match = trimmedLine.match(/^(CREATE\s+TABLE\s+)?["`']?(\w+)["`']?/i);
                if (match && match[2]) {
                    // Save previous table if exists
                    if (currentTable) {
                        tables.push(currentTable);
                    }

                    // Create new table
                    currentTable = {
                        name: match[2],
                        columns: [],
                        primaryKey: ''
                    };
                }
            }
            // Check if line defines a column
            else if (currentTable && trimmedLine.match(/^["`']?(\w+)["`']?\s+(\w+)/)) {
                const match = trimmedLine.match(/^["`']?(\w+)["`']?\s+(\w+)/);
                if (match && match[1] && match[2]) {
                    // Add column to current table
                    currentTable.columns.push({
                        name: match[1],
                        type: mapDataType(match[2]),
                        isForeignKey: false,
                        referencedTable: '',
                        referencedColumn: ''
                    });

                    // Check if this is primary key
                    if (trimmedLine.toLowerCase().includes('primary key')) {
                        currentTable.primaryKey = match[1];
                    }
                }
            }
        }

        // Add the last table
        if (currentTable) {
            tables.push(currentTable);
        }

        return tables;
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

    // Custom Data Type handlers
    const openCustomDataTypeModal = (tableIndex: number, columnIndex: number) => {
        selectedTableIndex = tableIndex;
        selectedColumnIndex = columnIndex;
        customDataType = schema.tables[tableIndex].columns[columnIndex].type || '';
        showCustomDataTypeModal = true;
    };

    const handleSetCustomDataType = () => {
        if (selectedTableIndex === null || selectedColumnIndex === null) return;

        updateColumn(selectedTableIndex, selectedColumnIndex, {
            type: customDataType
        });

        showCustomDataTypeModal = false;
        customDataType = '';
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

    const commonDataTypes = [
        'INT', 'VARCHAR', 'DATE', 'BOOLEAN', 'UUID', 'TEXT', 'TIMESTAMP',
        'NUMERIC', 'FLOAT', 'DECIMAL', 'JSONB', 'JSON'
    ];
</script>

<svelte:head>
    <title>Complex SQL Query</title>
</svelte:head>
<Navbar />
<div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <header class="text-center mb-8">
            <h1 class="text-4xl font-bold text-green-600">Complex SQL Generator</h1>
            <p class="mt-2 text-lg text-gray-600">Enter your query and define the schema with primary keys, foreign keys, and data types.</p>
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
                    <div class="flex justify-between items-center mb-4">
                        <h2 class="text-xl font-semibold text-gray-900">Schema</h2>
                        <button
                                type="button"
                                on:click={openSchemaUploadModal}
                                class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 flex items-center space-x-2"
                        >
                            <UploadCloud class="h-5 w-5" />
                            <span>Upload Schema</span>
                        </button>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
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
                                        <div class="mt-1 p-2 bg-gray-100 rounded flex justify-between items-center">
                                            <span>{table.primaryKey}</span>
                                            <button
                                                    type="button"
                                                    on:click={() => openPKeyModal(tableIndex)}
                                                    class="text-blue-600 hover:text-blue-800"
                                            >
                                                <Edit class="h-4 w-4" />
                                            </button>
                                        </div>
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
                                                <div class="flex items-center">
                                                    <select
                                                            value={column.type}
                                                            on:change={(e) => updateColumn(tableIndex, columnIndex, { type: e.currentTarget.value })}
                                                            class="mt-1 block w-full p-2 border border-gray-300 rounded-l-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                                                    >
                                                        <option value="">Select type</option>
                                                        {#each commonDataTypes as dataType}
                                                            <option value={dataType}>{dataType}</option>
                                                        {/each}
                                                        <option value="CUSTOM">Custom...</option>
                                                    </select>
                                                    <button
                                                            type="button"
                                                            on:click={() => openCustomDataTypeModal(tableIndex, columnIndex)}
                                                            class="mt-1 px-3 py-2 bg-gray-200 border border-gray-300 rounded-r-md"
                                                    >
                                                        <Edit class="h-4 w-4" />
                                                    </button>
                                                </div>
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

    <!-- Schema Upload Modal -->
    {#if showSchemaUploadModal}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-2xl">
                <h3 class="text-lg font-semibold mb-2">Upload Database Schema</h3>
                <p class="text-gray-600 mb-4">
                    Paste the output of database schema queries like:
                    <code class="bg-gray-100 p-1 text-sm rounded">SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';</code>
                    or
                    <code class="bg-gray-100 p-1 text-sm rounded">SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_schema = 'public';</code>
                    or CSV-like format: table_name,column_name,data_type
                </p>

                <textarea
                        bind:value={schemaUploadText}
                        rows="10"
                        class="w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
                        placeholder="Paste database schema output here..."
                ></textarea>

                <div class="mt-6 flex justify-end space-x-3">
                    <button
                            on:click={() => showSchemaUploadModal = false}
                            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
                    >
                        Cancel
                    </button>
                    <button
                            on:click={processSchemaUpload}
                            disabled={!schemaUploadText.trim()}
                            class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        Process Schema
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- Primary Key Modal -->
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

    <!-- Foreign Key Modal -->
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

    <!-- Custom Data Type Modal -->
    {#if showCustomDataTypeModal && selectedTableIndex !== null && selectedColumnIndex !== null}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div class="bg-white rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-semibold mb-4">
                    Set Custom Data Type for {schema.tables[selectedTableIndex].columns[selectedColumnIndex].name}
                </h3>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Data Type</label>
                    <input
                            type="text"
                            bind:value={customDataType}
                            class="w-full p-2 border border-gray-300 rounded-md"
                            placeholder="e.g., TIMESTAMP WITH TIME ZONE"
                    />
                </div>
                <div class="mt-6 flex justify-end space-x-3">
                    <button
                            on:click={() => showCustomDataTypeModal = false}
                            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
                    >
                        Cancel
                    </button>
                    <button
                            on:click={handleSetCustomDataType}
                            disabled={!customDataType.trim()}
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