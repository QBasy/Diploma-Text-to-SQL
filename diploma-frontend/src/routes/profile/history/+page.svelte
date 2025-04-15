<script lang="ts">
    import { onMount } from 'svelte';
    import {
        History,
        Trash2,
        Database,
        Search,
        Calendar,
        Code,
        RefreshCcw
    } from 'lucide-svelte';
    import { getHistory, clearHistory} from '$lib/api'
    import type {QueryHistory} from '$lib/api'

    let histories: QueryHistory[] = [];
    let loading: boolean = true;
    let error: any = null;

    onMount(async () => {
        await fetchHistory();
    });

    async function fetchHistory() {
        loading = true;
        error = null;
        try {
            histories = await getHistory();
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    async function clear() {
        if (!confirm('Are you sure you want to clear all history?')) return;

        try {
            const response = await clearHistory();
            if (response) {
                histories = [];
            }

        } catch (err: any) {
            error = err.message;
        }
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
                    on:click={clearHistory}
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
                You haven't run any queries yet. Start from the homepage to generate your first one.
            </p>
        </div>
    {:else}
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
                        </div>

                        <div>
                            <h3 class="text-blue-600 font-medium text-sm mb-1">Query:</h3>
                            <div class="bg-slate-50 p-3 rounded-lg border text-sm font-mono overflow-x-auto">
                                {history.query}
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
    {/if}
</main>
