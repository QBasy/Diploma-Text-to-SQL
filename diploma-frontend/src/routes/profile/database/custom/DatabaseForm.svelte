<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { addCustomDatabase, type AddDatabasePayload } from '$lib/api/customDatabase';
    import { Database } from 'lucide-svelte';

    const dispatch = createEventDispatcher();

    let formData: AddDatabasePayload = {
        name: '',
        db_type: 'postgres',
        host: '',
        port: 5432,
        username: '',
        password: '',
        database: '',
        ssl_mode: 'disable',
        description: ''
    };

    let isSubmitting = false;
    let error: string | null = null;

    function updateDefaultPort() {
        switch (formData.db_type) {
            case 'postgres':
                formData.port = 5432;
                break;
            case 'mysql':
                formData.port = 3306;
                break;
            case 'sqlite':
                formData.port = 0;
                break;
        }
    }

    async function handleSubmit() {
        isSubmitting = true;
        error = null;

        try {
            const response = await addCustomDatabase(formData);
            dispatch('databaseAdded', response);
            resetForm();
        } catch (err: any) {
            error = err?.response?.data?.error || 'Failed to add database';
            console.error(err);
        } finally {
            isSubmitting = false;
        }
    }

    function resetForm() {
        formData = {
            name: '',
            db_type: 'postgres',
            host: '',
            port: 5432,
            username: '',
            password: '',
            database: '',
            ssl_mode: 'disable',
            description: ''
        };
    }
</script>


<form on:submit|preventDefault={handleSubmit} class="space-y-4">
    {#if error}
        <div class="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
            {error}
        </div>
    {/if}

    <div class="grid md:grid-cols-2 gap-4">
        <div class="space-y-2">
            <label for="name" class="block text-sm font-medium text-gray-700">Connection Name*</label>
            <input
                    type="text"
                    id="name"
                    bind:value={formData.name}
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    placeholder="My Database"
            />
        </div>

        <div class="space-y-2">
            <label for="db_type" class="block text-sm font-medium text-gray-700">Database Type*</label>
            <select
                    id="db_type"
                    bind:value={formData.db_type}
                    on:change={updateDefaultPort}
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
                <option value="postgres">PostgreSQL</option>
                <option value="mysql">MySQL</option>
                <option value="sqlite">SQLite</option>
            </select>
        </div>
    </div>

    {#if formData.db_type !== 'sqlite'}
        <div class="grid md:grid-cols-2 gap-4">
            <div class="space-y-2">
                <label for="host" class="block text-sm font-medium text-gray-700">Host*</label>
                <input
                        type="text"
                        id="host"
                        bind:value={formData.host}
                        required={formData.db_type !== 'sqlite'}
                        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                        placeholder="localhost"
                />
            </div>

            <div class="space-y-2">
                <label for="port" class="block text-sm font-medium text-gray-700">Port*</label>
                <input
                        type="number"
                        id="port"
                        bind:value={formData.port}
                        required={formData.db_type !== 'sqlite'}
                        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
            </div>
        </div>

        <div class="grid md:grid-cols-2 gap-4">
            <div class="space-y-2">
                <label for="username" class="block text-sm font-medium text-gray-700">Username*</label>
                <input
                        type="text"
                        id="username"
                        bind:value={formData.username}
                        required={formData.db_type !== 'sqlite'}
                        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
            </div>

            <div class="space-y-2">
                <label for="password" class="block text-sm font-medium text-gray-700">Password*</label>
                <input
                        type="password"
                        id="password"
                        bind:value={formData.password}
                        required={formData.db_type !== 'sqlite'}
                        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
            </div>
        </div>
    {/if}

    <div class="space-y-2">
        <label for="database" class="block text-sm font-medium text-gray-700">
            {formData.db_type === 'sqlite' ? 'Database File Path*' : 'Database Name*'}
        </label>
        <input
                type="text"
                id="database"
                bind:value={formData.database}
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={formData.db_type === 'sqlite' ? '/path/to/database.db' : 'db_name'}
        />
    </div>

    {#if formData.db_type === 'postgres'}
        <div class="space-y-2">
            <label for="ssl_mode" class="block text-sm font-medium text-gray-700">SSL Mode</label>
            <select
                    id="ssl_mode"
                    bind:value={formData.ssl_mode}
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
                <option value="disable">Disable</option>
                <option value="require">Require</option>
                <option value="verify-full">Verify Full</option>
            </select>
        </div>
    {/if}

    <div class="space-y-2">
        <label for="description" class="block text-sm font-medium text-gray-700">Description</label>
        <textarea
                id="description"
                bind:value={formData.description}
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Optional description about this database connection"
        ></textarea>
    </div>

    <div class="flex justify-end pt-4">
        <button
                type="submit"
                disabled={isSubmitting}
                class="inline-flex items-center gap-2 px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
            {#if isSubmitting}
                <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                Submitting...
            {:else}
                <Database class="w-4 h-4" />
                Add Database
            {/if}
        </button>
    </div>
</form>