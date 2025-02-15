<script lang="ts">
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import {createTableAPI, executeQueryAPI} from "$lib/api";
    import {goto} from "$app/navigation";
    import Visualise from "$lib/components/Visualise.svelte";

    let activeTab = 'query';
    let sqlQuery = '';
    let naturalLanguageQuery = '';
    let queryResult: string | null = null;
    let loading = false;
    let notification = '';
    let notificationType = 'success';

    let tableForm = {
        name: '',
        columns: [{ name: '', type: 'TEXT' }]
    };

    function addColumn() {
        tableForm.columns = [...tableForm.columns, { name: '', type: 'TEXT' }];
    }

    function removeColumn(index: number) {
        tableForm.columns = tableForm.columns.filter((_, i) => i !== index);
    }

    async function executeQuery() {
        loading = true;
        try {
            const result = await executeQueryAPI(sqlQuery)
            queryResult = 'Query executed successfully';
            notification = 'Query processed successfully: ' + result;
            notificationType = 'success';
        } catch (error) {
            notification = 'Error executing query';
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
            const result = await createTableAPI(tableForm)
            notification = `Table ${tableForm.name} created successfully`;
            notificationType = 'success';
            tableForm = { name: '', columns: [{ name: '', type: 'TEXT' }] };
        } catch (error) {
            notification = 'Error creating table';
            notificationType = 'error';
        } finally {
            loading = false;
        }
    }
</script>

<div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6 text-gray-800">Database Management</h1>

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
                class="px-4 py-2 mr-2 {activeTab === 'table' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                on:click={() => activeTab = 'table'}
        >
            Create Table
        </button>
        <button
                class="px-4 py-2 {activeTab === 'visualize' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                on:click={() => goto('/profile/database/visualize')}
        >
            Visualization
        </button>
    </div>

    {#if loading}
        <LoadingSpinner />
    {:else}
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
        {:else if activeTab === 'table'}
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
    {/if}

    {#if queryResult}
        <div class="mt-4 bg-white shadow-md rounded-lg p-6">
            <h2 class="text-xl font-semibold mb-4">Query Result</h2>
            <pre class="bg-gray-100 p-4 rounded-md overflow-x-auto">{queryResult}</pre>
        </div>
    {/if}

    {#if notification}
        <Notification message={notification} type={notificationType} />
    {/if}

    <div class="mx-auto border border-l-gray-200 rounded-sm w-full my-5">
        <Visualise />
    </div>
</div>
