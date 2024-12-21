<script lang="ts">
    import { onMount } from "svelte";
    import { getUserDatabase, createTable, executeQuery } from "$lib/database";
    import {goto} from "$app/navigation";
    import {checkAuthStatus} from "$lib/stores/authStore";

    interface Column {
        name: string;
        type: string;
    }

    interface TableData {
        name: string;
        columns: Column[];
    }

    interface QueryResult {
        success: boolean;
        data?: any;
        error?: string;
    }

    let database: any = null;
    let tableName: string = "";
    let columns: Column[] = [{ name: "", type: "" }];
    let rawTableQuery: string = "";
    let customQuery: string = "";
    let queryResult: QueryResult | null = null;

    let isRaw: boolean = false;

    const toggleRaw = (): void => {
        isRaw = !isRaw;
    }

    // Получить данные базы данных
    const getDatabase = async (): Promise<void> => {
        try {
            database = await getUserDatabase();
        } catch (error) {
            console.error("Failed to fetch database:", error);
        }
    };

    // Добавить колонку для таблицы
    const addColumn = (): void => {
        columns.push({ name: "", type: "" });
    };

    // Создать таблицу через options
    const createTableOptions = async (): Promise<void> => {
        try {
            const tableData: TableData = {
                name: tableName,
                columns: columns.filter((col) => col.name && col.type),
            };
            const response = await createTable(tableData);
            console.log("Table created:", response);
            getDatabase();
        } catch (error) {
            console.error("Failed to create table with options:", error);
        }
    };

    // Создать таблицу через raw text
    const createTableRaw = async (): Promise<void> => {
        try {
            const response = await executeQuery({ query: rawTableQuery });
            console.log("Raw table created:", response);
            getDatabase();
        } catch (error) {
            console.error("Failed to create table with raw query:", error);
        }
    };

    // Выполнить пользовательский запрос
    const executeCustomQuery = async (): Promise<void> => {
        try {
            queryResult = await executeQuery({ query: customQuery });
        } catch (error) {
            console.error("Failed to execute custom query:", error);
        }
    };

    onMount(() => {
        if (!checkAuthStatus()) {
            alert(
                "You must be logged in to access this page. Please log in or register."
            );
            goto('/auth');
        } else {
            getDatabase();
        }
    });
</script>

<section class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-4">Database Management</h1>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Секция: Текущая база данных -->
        <div>
            <h2 class="text-xl font-semibold mb-2">Current Database</h2>
            {#if database}
                <pre class="bg-gray-100 p-4 rounded-lg shadow-md">
                    {JSON.stringify(database, null, 2)}
                </pre>
            {:else}
                <p>Loading database...</p>
            {/if}
        </div>

        {#if !isRaw}
            <!-- Секция: Создание таблицы -->
            <div>
                <h2 class="text-xl font-semibold mb-2">Create Table</h2>
                <div class="mb-4">
                    <label class="block font-medium mb-2">Table Name</label>
                    <input
                            type="text"
                            class="w-full p-2 border rounded"
                            placeholder="Enter table name"
                            bind:value={tableName}
                    />
                </div>

                <div class="mb-4">
                    <h3 class="font-medium mb-2">Columns</h3>
                    {#each columns as column, index}
                        <div class="flex gap-4 mb-2">
                            <input
                                    type="text"
                                    class="flex-1 p-2 border rounded"
                                    placeholder="Column name"
                                    bind:value={columns[index].name}
                            />
                            <select
                                    class="flex-1 p-2 border rounded"
                                    bind:value={columns[index].type}
                            >
                                <option value="" disabled>Select type</option>
                                <option value="int">Integer</option>
                                <option value="text">Text</option>
                                <option value="boolean">Boolean</option>
                            </select>
                        </div>
                    {/each}
                    <button
                            class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                            on:click={addColumn}
                    >
                        Add Column
                    </button>
                </div>

                <button
                        class="w-full bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
                        on:click={createTableOptions}
                >
                    Create Table (Options)
                </button>
                <button
                        class="w-full bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
                        on:click={toggleRaw}
                >
                    Use Raw
                </button>
            </div>
        {:else}
            <!-- Секция: Создание таблицы через Raw Query -->
            <div class="mb-4">
                <div>
                    <h2 class="text-xl font-semibold mb-2">Create Table (Raw Query)</h2>
                    <textarea
                            class="w-full p-4 border rounded mb-4"
                            rows="5"
                            placeholder="Enter raw SQL query"
                            bind:value={rawTableQuery}
                    ></textarea>
                    <button
                            class="w-full bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                            on:click={createTableRaw}
                    >
                        Create Table (Raw Query)
                    </button>
                    <button
                            class="w-full bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
                            on:click={toggleRaw}
                    >
                        Use Structured
                    </button>
                </div>
            </div>
        {/if}
    </div>

    <!-- Секция: Выполнение пользовательского запроса -->
    <div class="mt-8 w-3/4 mx-auto">
        <h2 class="text-xl font-semibold mb-2">Custom Query</h2>
        <textarea
                class="w-full p-4 border rounded mb-4"
                rows="5"
                placeholder="Enter custom query"
                bind:value={customQuery}
        ></textarea>
        <button
                class="w-full bg-indigo-500 text-white px-4 py-2 rounded hover:bg-indigo-600"
                on:click={executeCustomQuery}
        >
            Execute Query
        </button>

        {#if queryResult}
            <div class="mt-4">
                <h3 class="text-lg font-semibold">Query Result</h3>
                <pre class="bg-gray-100 p-4 rounded-lg shadow-md">
                    {JSON.stringify(queryResult, null, 2)}
                </pre>
            </div>
        {/if}
    </div>
</section>
