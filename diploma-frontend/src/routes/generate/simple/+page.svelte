<script lang="ts">
    import { generateSimpleSQL } from '$lib/api';
    import Notification from '$lib/components/Notification.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import { CheckCircle, AlertTriangle } from 'lucide-svelte';

    let isLoading = false;
    let notificationMessage = '';
    let notificationType: 'success' | 'error' = 'success';
    let query = '';
    let sqlResult = '';

    const submitQuery = async () => {
        isLoading = true;
        try {
            const response = await generateSimpleSQL(query);
            sqlResult = response.sql;
            notificationMessage = 'SQL generated successfully!';
            notificationType = 'success';
        } catch (err: any) {
            notificationMessage = `Error: ${err.message}`;
            notificationType = 'error';
        } finally {
            isLoading = false;
        }
    };
</script>

<svelte:head>
    <title>Simple SQL Generation</title>
</svelte:head>

<section class="bg-gray-50 min-h-screen py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto">
        <div class="text-center mb-10">
            <h1 class="text-4xl font-bold text-blue-700 tracking-tight">Simple SQL Generator</h1>
            <p class="mt-2 text-lg text-gray-600">Just describe what you want, and we’ll generate the SQL for you.</p>
        </div>

        <form on:submit|preventDefault={submitQuery} class="bg-white p-6 rounded-2xl shadow-lg space-y-6">
            <div>
                <label for="query" class="block text-sm font-semibold text-gray-700 mb-1">Your Query in English</label>
                <textarea
                        id="query"
                        bind:value={query}
                        rows="4"
                        placeholder="e.g., List all users who joined after 2022"
                        class="w-full p-3 border border-gray-300 rounded-xl shadow-sm focus:ring-blue-500 focus:border-blue-500 resize-none text-sm placeholder-gray-400"
                ></textarea>
            </div>

            <button
                    type="submit"
                    class="w-full flex justify-center items-center gap-2 px-4 py-2 bg-blue-600 text-white font-medium rounded-xl hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-blue-500 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={isLoading}
            >
                {#if isLoading}
                    <LoadingSpinner size={20} color="white" />
                    Generating...
                {:else}
                    Generate SQL
                {/if}
            </button>
        </form>

        {#if sqlResult}
            <div class="mt-10 bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-lg font-semibold text-gray-800 mb-2 flex items-center gap-2">
                    <CheckCircle class="text-green-500 w-5 h-5" />
                    Generated SQL
                </h2>
                <pre class="p-4 bg-gray-100 rounded-xl text-sm font-mono overflow-x-auto">{sqlResult}</pre>
            </div>
        {/if}

        {#if notificationMessage}
            <Notification {notificationMessage} {notificationType}>
                {#if notificationType === 'success'}
                    <CheckCircle class="text-green-500 w-5 h-5" />
                {:else}
                    <AlertTriangle class="text-red-500 w-5 h-5" />
                {/if}
            </Notification>
        {/if}
    </div>
</section>
