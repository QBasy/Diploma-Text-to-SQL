<script lang="ts">
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';
    import { textToSqlService, customQueryService, itemService } from '$lib/api';

    interface Item {
        id: number;
        name: string;
        type: string;
        owner_id: number;
    }

    interface SQLResult {
        columns: string[];
        data: Record<string, any>[];
    }

    const naturalLanguageQuery = writable('');
    const sqlQuery = writable('');
    const queryResult = writable<SQLResult | null>(null);
    const error = writable<string | null>(null);

    let items: Item[] = [];
    let newItem = { name: '', type: '' };
    let fetchError: string | null = null;

    async function fetchItems() {
        try {
            items = await itemService.getItems();
        } catch (err) {
            fetchError = err instanceof Error ? err.message : 'Failed to fetch items';
        }
    }

    async function createItem() {
        try {
            const createdItem = await itemService.createItem(newItem);
            items = [...items, createdItem];
            newItem = { name: '', type: '' }; // Сброс формы
        } catch (err) {
            fetchError = err instanceof Error ? err.message : 'Failed to create item';
        }
    }

    async function handleNaturalLanguageSubmit() {
        try {
            const { sql } = await textToSqlService.convertQuery($naturalLanguageQuery);
            sqlQuery.set(sql);

            const result = await customQueryService.executeQuery(sql);
            queryResult.set(result);
            error.set(null);
        } catch (err) {
            error.set('An error occurred while executing the query');
            console.error(err);
        }
    }

    onMount(() => {
        fetchItems();
    });
</script>

<div class="container mx-auto p-6 space-y-8">
    <h1 class="text-3xl font-bold">Integrated Management System</h1>

    <section class="bg-gray-100 p-4 rounded shadow">
        <h2 class="text-2xl font-semibold mb-4">Items Management</h2>

        <div class="mb-4">
            <h3 class="font-semibold mb-2">Create New Item:</h3>
            <form on:submit|preventDefault={createItem} class="space-y-4">
                <div>
                    <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
                    <input
                            type="text"
                            id="name"
                            bind:value={newItem.name}
                            required
                            class="w-full p-2 border rounded"
                    />
                </div>
                <div>
                    <label for="type" class="block text-sm font-medium text-gray-700">Type</label>
                    <input
                            type="text"
                            id="type"
                            bind:value={newItem.type}
                            required
                            class="w-full p-2 border rounded"
                    />
                </div>
                <button
                        type="submit"
                        class="bg-indigo-500 text-white px-4 py-2 rounded hover:bg-indigo-600"
                >
                    Create Item
                </button>
            </form>
        </div>

        <div>
            <h3 class="font-semibold mb-2">Existing Items:</h3>
            {#if items.length === 0}
                <p class="text-gray-500">No items found.</p>
            {:else}
                <table class="w-full border-collapse border">
                    <thead class="bg-gray-50">
                    <tr>
                        <th class="border p-2">ID</th>
                        <th class="border p-2">Name</th>
                        <th class="border p-2">Type</th>
                        <th class="border p-2">Owner ID</th>
                    </tr>
                    </thead>
                    <tbody>
                    {#each items as item}
                        <tr>
                            <td class="border p-2">{item.id}</td>
                            <td class="border p-2">{item.name}</td>
                            <td class="border p-2">{item.type}</td>
                            <td class="border p-2">{item.owner_id}</td>
                        </tr>
                    {/each}
                    </tbody>
                </table>
            {/if}
        </div>
    </section>
</div>
