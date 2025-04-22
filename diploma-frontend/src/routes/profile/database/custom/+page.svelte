<script lang="ts">
    import { onMount } from 'svelte';
    import { listCustomDatabases, deleteCustomDatabase, type CustomDatabase } from '$lib/api/customDatabase';
    import { Plus, Trash2, Database, RefreshCw, Eye } from 'lucide-svelte';
    import DatabaseForm from './DatabaseForm.svelte';
    import SchemaViewer from './SchemaViewer.svelte';

    let customDatabases: CustomDatabase[] = [];
    let isLoading = true;
    let showAddForm = false;
    let error: string | null = null;
    let success: string | null = null;
    let selectedDatabaseUUID: string | null = null;
    let showSchema = false;

    async function fetchCustomDatabases() {
        isLoading = true;
        try {
            customDatabases = await listCustomDatabases();
            console.log(customDatabases);
        } catch (err) {
            error = 'Failed to load databases';
            console.error(err);
        } finally {
            isLoading = false;
            console.log(isLoading);
        }
    }

    async function deleteDatabase(uuid: string) {
        try {
            await deleteCustomDatabase(uuid);
            success = 'Database deleted successfully';
            fetchCustomDatabases();
        } catch (err) {
            error = 'Failed to delete database';
            console.error(err);
        }
    }

    function handleDatabaseAdded(event: CustomEvent) {
        success = 'Database added successfully';
        showAddForm = false;
        fetchCustomDatabases();
    }

    function viewSchema(uuid: string) {
        selectedDatabaseUUID = uuid;
        showSchema = true;
        console.log(selectedDatabaseUUID, showSchema);
    }

    function clearAlerts() {
        setTimeout(() => {
            error = null;
            success = null;
        }, 3000);
    }

    $: if (error || success) clearAlerts();

    onMount(fetchCustomDatabases);
</script>


<svelte:head>
    <title>Custom Database Management</title>
    <meta name="description" content="Manage and connect to your custom databases" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
    <div class="mb-8 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-800 flex items-center gap-2">
            <Database class="w-6 h-6 text-blue-600" />
            Custom Database Management
        </h1>

        <button
                on:click={() => showAddForm = !showAddForm}
                class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md transition-colors"
        >
            {#if showAddForm}
                Close Form
            {:else}
                <Plus size={18} />
                Add Database
            {/if}
        </button>
    </div>

    {#if error}
        <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-4 text-red-700">
            {error}
        </div>
    {/if}

    {#if success}
        <div class="bg-green-50 border-l-4 border-green-500 p-4 mb-4 text-green-700">
            {success}
        </div>
    {/if}

    {#if showAddForm}
        <div class="bg-gray-50 p-6 rounded-lg shadow-sm border mb-8">
            <h2 class="text-lg font-medium mb-4">Add New Database Connection</h2>
            <DatabaseForm on:databaseAdded={handleDatabaseAdded} />
        </div>
    {/if}

    {#if isLoading}
        <div class="flex justify-center p-12">
            <RefreshCw class="w-8 h-8 text-blue-600 animate-spin" />
        </div>
    {:else if customDatabases.length === 0}
        <div class="bg-gray-50 p-8 rounded-lg text-center border">
            <Database class="w-12 h-12 text-gray-400 mx-auto mb-3" />
            <h3 class="text-lg font-medium text-gray-700">No Custom Databases</h3>
            <p class="text-gray-500 mt-1">Add your first database connection to get started.</p>
        </div>
    {:else}
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {#each customDatabases as db}
                <div class="bg-white rounded-lg shadow border hover:shadow-md transition-shadow">
                    <div class="p-5">
                        <div class="flex justify-between items-start">
                            <h3 class="font-medium text-lg text-gray-800">{db.name}</h3>
                            <span class="px-2 py-1 rounded text-xs uppercase bg-blue-100 text-blue-800">{db.db_type}</span>
                        </div>

                        <p class="text-gray-600 text-sm mt-2 line-clamp-2">{db.description || 'No description provided'}</p>

                        <div class="mt-3 text-sm text-gray-500">
                            <p>Host: {db.host}:{db.port}</p>
                            <p>Database: {db.database}</p>
                        </div>
                    </div>

                    <div class="border-t px-5 py-3 flex justify-between items-center bg-gray-50 rounded-b-lg">
                        <button
                                on:click={() => viewSchema(db.UUID)}
                                class="text-blue-600 hover:text-blue-800 flex items-center gap-1 text-sm font-medium"
                        >
                            <Eye size={16} />
                            View Schema
                        </button>

                        <button
                                on:click={() => deleteDatabase(db.uuid)}
                                class="text-red-600 hover:text-red-800 flex items-center gap-1 text-sm font-medium"
                        >
                            <Trash2 size={16} />
                            Delete
                        </button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

{#if showSchema && selectedDatabaseUUID }
    <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-lg shadow-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
            <div class="p-4 border-b flex justify-between items-center">
                <h2 class="text-xl font-semibold">Database Schema</h2>
                <button
                        on:click={() => { showSchema = false; selectedDatabaseUUID = null; }}
                        class="p-2 text-gray-500 hover:text-gray-700"
                >
                    &times;
                </button>
            </div>

            <div class="flex-1 overflow-auto p-4">
                <SchemaViewer databaseUUID={selectedDatabaseUUID} />
            </div>
        </div>
    </div>
{/if}