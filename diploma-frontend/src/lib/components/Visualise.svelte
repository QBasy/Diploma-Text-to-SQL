<script lang="ts">
    import { onMount } from 'svelte';
    import { getSchema } from "$lib/api";
    import { goto } from "$app/navigation";

    let schema: Record<string, any> = {}; // Обновление типа данных для более точной работы
    let loading = true;
    let error: string | null = null;
    let message: string | null = null;
    let isEmpty = false;

    const loadSchema = async () => {
        try {
            const response: any = await getSchema();
            console.log("API Response:", response);

            if (response.status === "success") {
                if (response.schema && Object.keys(response.schema).length === 0) {
                    isEmpty = true;
                    message = response.message || "The database exists but does not contain any tables or connections.";
                } else {
                    schema = response.schema;
                }
            } else {
                throw new Error(response.message || "Failed to fetch schema");
            }
        } catch (err: any) {
            console.error("Error:", err);
            if (err.message === "Unauthorized") {
                goto('/login');
            } else {
                error = err.message || "Failed to fetch schema";
            }
        } finally {
            loading = false;
        }
    };

    const refreshSchema = async () => {
        loading = true;
        error = null;
        isEmpty = false;
        await loadSchema();
    };

    onMount(loadSchema);
</script>

<main class="p-6 bg-gray-100 min-h-screen">
    <h1 class="text-3xl font-bold mb-6 text-gray-800 text-center">Database Schema</h1>

    <div class="flex flex-col items-center justify-center">
        {#if loading}
            <div class="text-center text-gray-600">Loading schema...</div>
        {:else if error}
            <div class="text-center text-red-600">Error: {error}</div>
        {:else if isEmpty}
            <div class="text-center text-gray-600">
                The database exists but does not contain any tables or connections.
            </div>
        {:else}
            <div class="relative grid grid-cols-2 gap-6">
                {#each Object.entries(schema) as [table, data], index}
                    <div class="bg-white shadow-lg rounded-lg p-4 w-96 border-2 border-gray-200 relative">
                        <h2 class="text-xl font-semibold text-gray-700 mb-3 text-center">
                            {table}
                        </h2>
                        <ul class="space-y-2">
                            {#each data.columns as column}
                                <li class="bg-gray-50 p-2 rounded-md text-gray-600 text-center">
                                    {column}
                                </li>
                            {/each}
                        </ul>
                        <div class="absolute top-2 right-2 text-green-500 text-sm">PK: {data.primaryKey}</div>
                    </div>

                    {#if index < Object.entries(schema).length - 1}
                        {#each data.foreignKeys as foreignKey}
                            <div
                                    class="absolute transform translate-x-1/2"
                                    style="top: 50%; left: calc(50% + {index * 150}px);"
                            >
                                <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-6 w-6 text-blue-500"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                        stroke-width="2"
                                >
                                    <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            d="M17 9l4 4m0 0l-4 4m4-4H3"
                                    />
                                </svg>
                            </div>
                        {/each}
                    {/if}
                {/each}
            </div>
        {/if}

        <button
                class="mt-6 px-4 py-2 bg-blue-500 text-white rounded-lg shadow-md hover:bg-blue-600 transition"
                on:click={refreshSchema}
        >
            Refresh Schema
        </button>
    </div>
</main>
