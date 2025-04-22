<script lang="ts">
    import { onMount } from 'svelte';
    import { getCustomSchemaComplex } from '$lib/api/customDatabase';
    import { Database, RefreshCw, Table2, Key, Link, ArrowRightLeft, ChevronDown, ChevronRight } from 'lucide-svelte';

    export let databaseUUID: string;

    type ColumnInfo = {
        name: string;
        type: string;
        isForeignKey?: boolean;
        referencedTable?: string;
        referencedColumn?: string;
    };

    type TableInfo = {
        name: string;
        columns: ColumnInfo[];
        primaryKey: string;
    };


    let isLoading = true;
    let error: string | null = null;
    let schema: TableInfo[] = [];
    let expandedTables: Record<string, boolean> = {};

    function toggleTable(tableName: string) {
        expandedTables[tableName] = !expandedTables[tableName];
        expandedTables = { ...expandedTables };
    }

    async function fetchSchema() {
        isLoading = true;
        error = null;

        try {
            const response = await getCustomSchemaComplex(databaseUUID);
            console.log(response);
            if (response.status === "success") {
                schema = response.schema;
                schema.forEach(table => {
                    expandedTables[table.name] = true;
                });
            } else {
                error = 'Failed to fetch schema data';
            }
        } catch (err: any) {
            error = err?.response?.data?.error || 'An error occurred while fetching schema';
            console.error(err);
        } finally {
            isLoading = false;
        }
    }

    function isPrimaryKey(table: TableInfo, columnName: string): boolean {
        return table.primaryKey === columnName;
    }

    onMount(fetchSchema);
</script>


<div class="w-full">
    <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
            <Database class="w-5 h-5 text-blue-600" />
            Database Schema
        </h3>

        <button
                on:click={fetchSchema}
                class="text-blue-600 hover:text-blue-800 flex items-center gap-1 text-sm"
                disabled={isLoading}
        >
            <RefreshCw class={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
        </button>
    </div>

    {#if isLoading}
        <div class="flex justify-center p-12">
            <RefreshCw class="w-8 h-8 text-blue-600 animate-spin" />
        </div>
    {:else if error}
        <div class="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
            {error}
        </div>
    {:else if schema.length === 0}
        <div class="bg-gray-50 p-8 rounded-lg text-center border">
            <Table2 class="w-12 h-12 text-gray-400 mx-auto mb-3" />
            <h3 class="text-lg font-medium text-gray-700">No Tables Found</h3>
            <p class="text-gray-500 mt-1">This database doesn't have any tables or the connection failed.</p>
        </div>
    {:else}
        <div class="bg-white border rounded-lg shadow-sm">
            {#each schema as table (table.name)}
                <div class="border-b last:border-b-0">
                    <div
                            class="flex items-center justify-between p-3 cursor-pointer hover:bg-gray-50"
                            on:click={() => toggleTable(table.name)}
                    >
                        <div class="flex items-center gap-2">
                            {#if expandedTables[table.name]}
                                <ChevronDown class="w-4 h-4 text-gray-500" />
                            {:else}
                                <ChevronRight class="w-4 h-4 text-gray-500" />
                            {/if}
                            <Table2 class="w-5 h-5 text-blue-600" />
                            <span class="font-medium">{table.name}</span>
                            <span class="text-xs text-gray-500">({table.columns.length} columns)</span>
                        </div>
                    </div>

                    {#if expandedTables[table.name]}
                        <div class="px-4 pb-3 pt-1">
                            <table class="min-w-full divide-y divide-gray-200 text-sm">
                                <thead>
                                <tr>
                                    <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Column</th>
                                    <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                                    <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Properties</th>
                                </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                {#each table.columns as column (table.name + column.name)} <!-- Modified key here -->
                                    <tr class="hover:bg-gray-50">
                                        <td class="px-3 py-2 whitespace-nowrap">
                                            <div class="flex items-center gap-1">
                                                {#if isPrimaryKey(table, column.name)}
                                                    <Key class="w-3 h-3 text-amber-500" />
                                                {/if}
                                                {column.name}
                                            </div>
                                        </td>
                                        <td class="px-3 py-2 whitespace-nowrap text-gray-600">{column.type}</td>
                                        <td class="px-3 py-2 whitespace-nowrap">
                                            {#if isPrimaryKey(table, column.name)}
                                                <span class="px-2 py-1 text-xs rounded-full bg-amber-100 text-amber-800">Primary Key</span>
                                            {/if}

                                            {#if column.isForeignKey && column.referencedTable}
                                                <div class="flex items-center gap-1 text-xs text-blue-600 mt-1">
                                                    <Link class="w-3 h-3" />
                                                    <span>References</span>
                                                    <ArrowRightLeft class="w-3 h-3" />
                                                    <span class="font-medium">{column.referencedTable}.{column.referencedColumn}</span>
                                                </div>
                                            {/if}
                                        </td>
                                    </tr>
                                {/each}
                                </tbody>
                            </table>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>