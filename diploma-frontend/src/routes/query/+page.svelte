<script lang="ts">
    import { writable } from 'svelte/store';
    import { textToSqlService, customQueryService } from '$lib/api';

    interface SQLResult {
        columns: string[];
        data: Record<string, any>[];
    }

    const naturalLanguageQuery = writable('');
    const sqlQuery = writable('');
    const queryResult = writable<SQLResult | null>(null);
    const error = writable<string | null>(null);

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
</script>

<div class="container mx-auto p-6 space-y-8">
    <h1 class="text-3xl font-bold">Integrated Management System</h1>

    <!-- Text-to-SQL Section -->
    <section class="bg-gray-100 p-4 rounded shadow">
        <h2 class="text-2xl font-semibold mb-4">Text-to-SQL</h2>
        <div>
            <label for="natural-query" class="block mb-2">Enter natural language query:</label>
            <textarea
                    id="natural-query"
                    bind:value={$naturalLanguageQuery}
                    class="w-full p-2 border rounded"
                    rows="3"
            ></textarea>
        </div>
        <button
                on:click={handleNaturalLanguageSubmit}
                class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 mt-2"
        >
            Convert & Execute Query
        </button>

        {#if $sqlQuery}
            <div class="mt-4">
                <h3 class="font-semibold">Generated SQL Query:</h3>
                <pre class="bg-gray-200 p-2 rounded">{$sqlQuery}</pre>
            </div>
        {/if}

        {#if $queryResult}
            <div class="mt-4">
                <h3 class="font-semibold">Query Result:</h3>
                <table class="w-full border-collapse border">
                    <thead>
                    <tr>
                        {#each $queryResult.columns as column}
                            <th class="border p-2">{column}</th>
                        {/each}
                    </tr>
                    </thead>
                    <tbody>
                    {#each $queryResult.data as row}
                        <tr>
                            {#each $queryResult.columns as column}
                                <td class="border p-2">{row[column]}</td>
                            {/each}
                        </tr>
                    {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </section>
</div>
