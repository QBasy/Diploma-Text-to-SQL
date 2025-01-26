<script>
    import { onMount } from 'svelte';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

    let activeTab = 'query';
    let sqlQuery = '';
    let naturalLanguageQuery = '';
    let queryResult = null;
    let loading = false;
    let notification = '';

    let tableForm = {
        name: '',
        columns: [{ name: '', type: 'TEXT' }]
    };

    function addColumn() {
        tableForm.columns = [...tableForm.columns, { name: '', type: 'TEXT' }];
    }

    function removeColumn(index) {
        tableForm.columns.splice(index, 1);
        tableForm.columns = tableForm.columns;
    }

    async function executeQuery() {
        loading = true;
        try {
            // TODO: Implement actual query execution logic
            queryResult = 'Query executed successfully';
            notification = 'Query processed';
        } catch (error) {
            notification = 'Error executing query';
        } finally {
            loading = false;
        }
    }

    async function createTable() {
        loading = true;
        try {
            // TODO: Implement table creation logic
            notification = `Table ${tableForm.name} created successfully`;
            tableForm = { name: '', columns: [{ name: '', type: 'TEXT' }] };
        } catch (error) {
            notification = 'Error creating table';
        } finally {
            loading = false;
        }
    }
</script>

<div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">Database Management</h1>

    <div class="flex mb-4">
        <button
                class="px-4 py-2 mr-2 {activeTab === 'query' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                on:click={() => activeTab = 'query'}
        >
            SQL Query
        </button>
        <button
                class="px-4 py-2 mr-2 {activeTab === 'natural' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                on:click={() => activeTab = 'natural'}
        >
            Natural Language Query
        </button>
        <button
                class="px-4 py-2 {activeTab === 'table' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                on:click={() => activeTab = 'table'}
        >
            Create Table
        </button>
    </div>

    {#if activeTab === 'query'}
        <div class="bg-white shadow-md rounded-lg p-6">
            <textarea
                    bind:value={sqlQuery}
                    class="w-full h-40 p-2 border rounded-md"
                    placeholder="Enter your SQL query"
            ></textarea>
            <button
                    on:click={executeQuery}
                    class="mt-4 bg-green-500 text-white px-4 py-2 rounded-md"
            >
                Execute Query
            </button>
        </div>
    {:else if activeTab === 'natural'}
        <div class="bg-white shadow-md rounded-lg p-6">
            <textarea
                    bind:value={naturalLanguageQuery}
                    class="w-full h-40 p-2 border rounded-md"
                    placeholder="Enter your natural language query"
            ></textarea>
            <button
                    on:click={executeQuery}
                    class="mt-4 bg-blue-500 text-white px-4 py-2 rounded-md"
            >
                Generate and Execute
            </button>
        </div>
    {:else}
        <div class="bg-white shadow-md rounded-lg p-6">
            <input
                    type="text"
                    bind:value={tableForm.name}
                    placeholder="Table Name"
                    class="w-full p-2 border rounded-md mb-4"
            >

            {#each tableForm.columns as column, index}
                <div class="flex mb-2">
                    <input
                            type="text"
                            bind:value={column.name}
                            placeholder="Column Name"
                            class="flex-grow p-2 border rounded-md mr-2"
                    >
                    <select
                            bind:value={column.type}
                            class="p-2 border rounded-md"
                    >
                        <option value="TEXT">TEXT</option>
                        <option value="INTEGER">INTEGER</option>
                        <option value="REAL">REAL</option>
                    </select>
                    {#if tableForm.columns.length > 1}
                        <button
                                on:click={() => removeColumn(index)}
                                class="ml-2 bg-red-500 text-white px-2 py-1 rounded-md"
                        >
                            Remove
                        </button>
                    {/if}
                </div>
            {/each}

            <button
                    on:click={addColumn}
                    class="mt-2 bg-blue-500 text-white px-4 py-2 rounded-md mr-2"
            >
                Add Column
            </button>
            <button
                    on:click={createTable}
                    class="mt-2 bg-green-500 text-white px-4 py-2 rounded-md"
            >
                Create Table
            </button>
        </div>
    {/if}

    {#if loading}
        <LoadingSpinner />
    {/if}

    {#if queryResult}
        <div class="mt-4 bg-white shadow-md rounded-lg p-6">
            <h2 class="text-xl font-semibold mb-4">Query Result</h2>
            <pre class="bg-gray-100 p-4 rounded-md overflow-x-auto">{queryResult}</pre>
        </div>
    {/if}

    {#if notification}
        <Notification
                message={notification}
                type={notification.includes('successfully') ? 'success' : 'error'}
        />
    {/if}
</div>