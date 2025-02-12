<script lang="ts">
    import {onMount} from 'svelte';
    import {getSchema} from "$lib/api";
    import {goto} from "$app/navigation";

    let schema: Record<string, string[]> = {};
    let loading = true;
    let error: string | null = null;

    onMount(async () => {
        try {
            const data = await getSchema();
            schema = data.reduce((acc, table) => {
                acc[table.name] = table.columns.map((col) => col.name);
                return acc;
            }, {} as Record<string, string[]>);
        } catch (err: any) {
            if (err.message === "Unauthorized") {
                goto('/login');
            }
            error = err.message;
        } finally {
            loading = false;
        }
    });
</script>

<main class="p-6 bg-gray-100 min-h-screen">
    <h1 class="text-3xl font-bold mb-6 text-gray-800">Database Schema</h1>

    {#if loading}
        <div class="text-center text-gray-600">Loading schema...</div>
    {:else if error}
        <div class="text-center text-red-600">Error: {error}</div>
    {:else}
        <div class="space-y-4">
            {#each Object.entries(schema) as [table, columns]}
                <div class="bg-white shadow-md rounded-lg p-4">
                    <h2 class="text-xl font-semibold text-gray-700 mb-3">{table}</h2>
                    <ul class="space-y-2">
                        {#each columns as column}
                            <li class="bg-gray-50 p-2 rounded-md text-gray-600">{column}</li>
                        {/each}
                    </ul>
                </div>
            {/each}
        </div>
    {/if}
</main>
