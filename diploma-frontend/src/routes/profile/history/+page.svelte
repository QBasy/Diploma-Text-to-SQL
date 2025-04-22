<script lang="ts">
    import { onMount } from 'svelte';
    import {
        History,
        Trash2,
        Database,
        Search,
        Calendar,
        Code,
        RefreshCcw,
        Filter,
        ChevronLeft,
        ChevronRight
    } from 'lucide-svelte';
    import { getHistory, clearHistory } from '$lib/api'
    import type { QueryHistory } from '$lib/api'

    let histories: QueryHistory[] = [];
    let loading: boolean = true;
    let error: any = null;

    let currentPage = 1;
    let perPage = 10;
    let totalItems = 0;
    let totalPages = 0;

    let databaseUUID = "";
    let queryType = "text-to-sql";
    let startDate = "";
    let endDate = "";
    let searchQuery = "";
    let sortBy = "timestamp";
    let sortDir = "DESC";

    let databaseOptions: string[] = [];
    let queryTypeOptions: string[] = [];

    onMount(async () => {
        await fetchHistory();
        await fetchFilterOptions();
    });

    async function fetchFilterOptions() {
        try {
            // In Future when we will be able to add many user databases const dbResponse = await getDatabases();
            // databaseOptions = dbResponse.map(db => db.uuid);

            queryTypeOptions = ["text-to-sql", "execute-sql"];
        } catch (err: any) {
            console.error("Failed to fetch filter options:", err);
        }
    }

    async function fetchHistory() {
        loading = true;
        error = null;
        try {
            const queryParams = new URLSearchParams();

            queryParams.append('page', currentPage.toString());
            queryParams.append('per_page', perPage.toString());

            if (databaseUUID) queryParams.append('database_uuid', databaseUUID);
            if (queryType) queryParams.append('query_type', queryType);
            if (startDate) queryParams.append('start_date', startDate);
            if (endDate) queryParams.append('end_date', endDate);
            if (searchQuery) queryParams.append('search', searchQuery);

            queryParams.append('sort_by', sortBy);
            queryParams.append('sort_dir', sortDir);

            const response = await getHistory(queryParams.toString());
            histories = response.data;
            totalItems = response.total;
            totalPages = response.last_page;
            currentPage = response.page;
            console.log(response)
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    function resetFilters() {
        databaseUUID = "";
        queryType = "text-to-sql";
        startDate = "";
        endDate = "";
        searchQuery = "";
        sortBy = "timestamp";
        sortDir = "DESC";
        currentPage = 1;
        fetchHistory();
    }

    async function clear() {
        if (!confirm('Are you sure you want to clear all history?')) return;

        try {
            const response = await clearHistory();
            if (response) {
                histories = [];
                totalItems = 0;
                totalPages = 0;
            }
        } catch (err: any) {
            error = err.message;
        }
    }

    function changePage(newPage: number) {
        if (newPage < 1 || newPage > totalPages) return;
        currentPage = newPage;
        fetchHistory();
    }

    function changeSort(field: string) {
        if (sortBy === field) {
            sortDir = sortDir === "ASC" ? "DESC" : "ASC";
        } else {
            sortBy = field;
            sortDir = "DESC";
        }
        currentPage = 1;
        fetchHistory();
    }

    function applyFilters() {
        currentPage = 1;
        fetchHistory();
    }

    function formatDate(dateString: Date | string): string {
        return new Date(dateString).toLocaleString();
    }
</script>

<svelte:head>
    <title>Query History | Text-To-SQL</title>
</svelte:head>

<main class="flex-grow md:w-11/12 container mx-auto px-4 py-12">
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
        <div class="flex items-center gap-3">
            <History class="w-7 h-7 text-blue-600" />
            <h1 class="text-3xl font-bold text-blue-600 tracking-tight">Query History</h1>
        </div>

        <div class="flex gap-3">
            <button
                    on:click={fetchHistory}
                    class="flex items-center gap-2 bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-xl transition"
            >
                <RefreshCcw class="w-5 h-5" />
                <span>Refresh</span>
            </button>

            <button
                    on:click={clear}
                    class="flex items-center gap-2 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-xl transition"
            >
                <Trash2 class="w-5 h-5" />
                <span>Clear</span>
            </button>
        </div>
    </div>

    <p class="mb-8 text-lg text-gray-600 max-w-3xl">
        Review and reuse your previously generated queries. Your history helps you track and improve your SQL generation workflow.
    </p>

    <div class="bg-white p-5 rounded-2xl border shadow-sm mb-8">
        <div class="flex items-center gap-2 mb-4">
            <Filter class="w-5 h-5 text-blue-600" />
            <h2 class="text-lg font-semibold">Filters</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1" for="search">
                    Search Query
                </label>
                <div class="relative">
                    <input
                            id="search"
                            type="text"
                            bind:value={searchQuery}
                            placeholder="Search in queries..."
                            class="w-full p-2 pl-9 border rounded-lg focus:ring-2 focus:ring-blue-600 focus:border-blue-600"
                    />
                    <Search class="w-4 h-4 text-gray-400 absolute left-3 top-3" />
                </div>
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1" for="database">
                    Database
                </label>
                <select
                        id="database"
                        bind:value={databaseUUID}
                        class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-600 focus:border-blue-600"
                >
                    <option value="">All Databases</option>
                    {#each databaseOptions as db}
                        <option value={db}>{db}</option>
                    {/each}
                </select>
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1" for="queryType">
                    Query Type
                </label>
                <select
                        id="queryType"
                        bind:value={queryType}
                        class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-600 focus:border-blue-600"
                >
                    <option value="">All Types</option>
                    {#each queryTypeOptions as type}
                        <option value={type}>{type}</option>
                    {/each}
                </select>
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1" for="startDate">
                    From Date
                </label>
                <input
                        id="startDate"
                        type="datetime-local"
                        bind:value={startDate}
                        class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-600 focus:border-blue-600"
                />
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1" for="endDate">
                    To Date
                </label>
                <input
                        id="endDate"
                        type="datetime-local"
                        bind:value={endDate}
                        class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-600 focus:border-blue-600"
                />
            </div>

            <div class="flex items-end gap-3">
                <button
                        on:click={applyFilters}
                        class="flex-1 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition"
                >
                    Apply Filters
                </button>
                <button
                        on:click={resetFilters}
                        class="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg transition"
                >
                    Reset
                </button>
            </div>
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-600"></div>
        </div>
    {:else if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            <p>Error: {error}</p>
            <p class="text-sm mt-1">Please try again later.</p>
        </div>
    {:else if histories.length === 0}
        <div class="bg-slate-50 border border-slate-200 rounded-xl p-8 text-center">
            <History class="w-12 h-12 mx-auto text-slate-400 mb-3" />
            <h2 class="text-xl font-semibold text-slate-700 mb-2">No Query History</h2>
            <p class="text-slate-500 max-w-md mx-auto">
                {searchQuery || databaseUUID || queryType || startDate || endDate
                    ? "No results match your filter criteria. Try adjusting your filters."
                    : "You haven't run any queries yet. Start from the homepage to generate your first one."}
            </p>
        </div>
    {:else}
        <div class="flex justify-end mb-4">
            <div class="flex gap-2 text-sm">
                <button
                        class="flex items-center gap-1 px-3 py-1 rounded {sortBy === 'timestamp' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100'}"
                        on:click={() => changeSort('timestamp')}
                >
                    Date {sortBy === 'timestamp' && (sortDir === 'ASC' ? '↑' : '↓')}
                </button>
                <button
                        class="flex items-center gap-1 px-3 py-1 rounded {sortBy === 'query_type' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100'}"
                        on:click={() => changeSort('query_type')}
                >
                    Type {sortBy === 'query_type' && (sortDir === 'ASC' ? '↑' : '↓')}
                </button>
                <button
                        class="flex items-center gap-1 px-3 py-1 rounded {sortBy === 'success' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100'}"
                        on:click={() => changeSort('success')}
                >
                    Status {sortBy === 'success' && (sortDir === 'ASC' ? '↑' : '↓')}
                </button>
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
            {#each histories as history}
                <div class="bg-white p-5 rounded-2xl border shadow-sm hover:shadow-lg transition-all flex flex-col justify-between">
                    <div class="space-y-4">
                        <div class="flex flex-wrap gap-3 text-sm text-gray-600">
                            <span class="flex items-center gap-1">
                                <Database class="w-4 h-4" /> {history.database_uuid}
                            </span>
                            <span class="flex items-center gap-1">
                                <Search class="w-4 h-4" /> {history.query_type}
                            </span>
                            <span class="flex items-center gap-1">
                                <Calendar class="w-4 h-4" /> {formatDate(history.timestamp)}
                            </span>
                            <span class={`px-2 py-0.5 rounded-full text-xs ${history.success ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                                {history.success ? 'Success' : 'Failed'}
                            </span>
                        </div>

                        <div>
                            <h3 class="text-blue-600 font-medium text-sm mb-1">Query:</h3>
                            <div class="bg-slate-50 p-3 rounded-lg border text-sm font-mono overflow-x-auto">
                                {#if history.query === ""}
                                    Database-Schema
                                    {:else}
                                    {history.query}
                                {/if}
                            </div>
                        </div>

                        <div>
                            <h3 class="text-green-600 font-medium text-sm mb-1">Result:</h3>
                            <div class="bg-slate-50 p-3 rounded-lg border text-sm font-mono overflow-x-auto">
                                <pre>{JSON.stringify(JSON.parse(history.result), null, 2)}</pre>
                            </div>
                        </div>
                    </div>

                    <div class="mt-4 flex justify-end">
                        <button
                                class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-xl transition"
                        >
                            <Code class="w-4 h-4" />
                            <span>Reuse Query</span>
                        </button>
                    </div>
                </div>
            {/each}
        </div>

        {#if totalPages > 1}
            <div class="mt-8 flex justify-center">
                <div class="flex items-center gap-2">
                    <!-- Prev -->
                    <button
                            class="p-2 rounded-lg border {currentPage === 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'hover:bg-gray-100'}"
                            disabled={currentPage === 1}
                            on:click={() => changePage(currentPage - 1)}
                    >
                        <ChevronLeft class="w-5 h-5" />
                    </button>

                    {#if totalPages <= 7}
                        {#each Array(totalPages) as _, i}
                            <button
                                    class="w-10 h-10 rounded-lg {currentPage === i + 1 ? 'bg-blue-600 text-white' : 'bg-white hover:bg-gray-100 border'}"
                                    on:click={() => changePage(i + 1)}
                            >
                                {i + 1}
                            </button>
                        {/each}
                    {:else}
                        <!-- First page -->
                        <button
                                class="w-10 h-10 rounded-lg {currentPage === 1 ? 'bg-blue-600 text-white' : 'bg-white hover:bg-gray-100 border'}"
                                on:click={() => changePage(1)}
                        >
                            1
                        </button>

                        {#if currentPage > 3}
                            <span class="text-gray-500">...</span>
                        {/if}

                        {#each Array(3) as _, i}
                            {@const pageNum = currentPage - 1 + i}
                            {#if pageNum > 1 && pageNum < totalPages}
                                <button
                                        class="w-10 h-10 rounded-lg {currentPage === pageNum ? 'bg-blue-600 text-white' : 'bg-white hover:bg-gray-100 border'}"
                                        on:click={() => changePage(pageNum)}
                                >
                                    {pageNum}
                                </button>
                            {/if}
                        {/each}

                        {#if currentPage < totalPages - 2}
                            <span class="text-gray-500">...</span>
                        {/if}

                        <!-- Last page -->
                        <button
                                class="w-10 h-10 rounded-lg {currentPage === totalPages ? 'bg-blue-600 text-white' : 'bg-white hover:bg-gray-100 border'}"
                                on:click={() => changePage(totalPages)}
                        >
                            {totalPages}
                        </button>
                    {/if}

                    <!-- Next -->
                    <button
                            class="p-2 rounded-lg border {currentPage === totalPages ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'hover:bg-gray-100'}"
                            disabled={currentPage === totalPages}
                            on:click={() => changePage(currentPage + 1)}
                    >
                        <ChevronRight class="w-5 h-5" />
                    </button>
                </div>
            </div>
        {/if}

        <div class="mt-4 text-center text-sm text-gray-600">
            Showing {(currentPage - 1) * perPage + 1} to {Math.min(currentPage * perPage, totalItems)} of {totalItems} results
        </div>
    {/if}
</main>