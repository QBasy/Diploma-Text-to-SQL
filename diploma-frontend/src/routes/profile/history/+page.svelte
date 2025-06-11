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
        ChevronRight,
        CheckCircle,
        XCircle,
        Copy,
        Eye,
        EyeOff,
        Clock,
        Tag,
        ArrowUpDown,
        ArrowUp,
        ArrowDown,
        X,
        Plus,
        Settings,
        MoreHorizontal,
        Menu
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

    // Enhanced filtering
    let databaseUUID = "";
    let queryTypes: string[] = [];
    let successFilter = "";
    let startDate = "";
    let endDate = "";
    let searchQuery = "";

    // Multiple sorting
    interface SortOption {
        field: string;
        direction: 'ASC' | 'DESC';
    }

    let sortOptions: SortOption[] = [
        { field: 'timestamp', direction: 'DESC' }
    ];

    let databaseOptions: string[] = [];
    let queryTypeOptions: string[] = [];

    // State management
    let expandedCards: Set<number> = new Set();
    let showFilters = false;
    let showAdvancedSort = false;
    let isMobile = false;

    // Mobile detection
    onMount(async () => {
        checkMobile();
        window.addEventListener('resize', checkMobile);
        await fetchHistory();
        await fetchFilterOptions();

        return () => {
            window.removeEventListener('resize', checkMobile);
        };
    });

    function checkMobile() {
        isMobile = window.innerWidth < 768;
    }

    async function fetchFilterOptions() {
        try {
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

            // Multiple query types
            queryTypes.forEach(type => {
                queryParams.append('query_type[]', type);
            });

            if (successFilter) queryParams.append('success', successFilter);
            if (startDate) queryParams.append('start_date', startDate);
            if (endDate) queryParams.append('end_date', endDate);
            if (searchQuery) queryParams.append('search', searchQuery);

            // Multiple sorting - send as JSON or multiple params
            sortOptions.forEach((sort, index) => {
                queryParams.append(`sort_${index}_by`, sort.field);
                queryParams.append(`sort_${index}_dir`, sort.direction);
            });

            const response = await getHistory(queryParams.toString());
            histories = response.data;
            totalItems = response.total;
            totalPages = response.last_page;
            currentPage = response.page;
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    function resetFilters() {
        databaseUUID = "";
        queryTypes = [];
        successFilter = "";
        startDate = "";
        endDate = "";
        searchQuery = "";
        sortOptions = [{ field: 'timestamp', direction: 'DESC' }];
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

    // Enhanced sorting functions
    function addSortOption() {
        if (sortOptions.length < 3) {
            sortOptions = [...sortOptions, { field: 'timestamp', direction: 'DESC' }];
        }
    }

    function removeSortOption(index: number) {
        if (sortOptions.length > 1) {
            sortOptions = sortOptions.filter((_, i) => i !== index);
            applySort();
        }
    }

    function updateSortOption(index: number, field: string, direction: 'ASC' | 'DESC') {
        sortOptions[index] = { field, direction };
        applySort();
    }

    function quickSort(field: string) {
        const existingIndex = sortOptions.findIndex(s => s.field === field);

        if (existingIndex >= 0) {
            // Toggle direction
            sortOptions[existingIndex].direction =
                sortOptions[existingIndex].direction === 'ASC' ? 'DESC' : 'ASC';
        } else {
            // Add new sort (move to front)
            sortOptions = [{ field, direction: 'DESC' }, ...sortOptions.slice(0, 2)];
        }

        applySort();
    }

    function applySort() {
        currentPage = 1;
        fetchHistory();
    }

    function applyFilters() {
        currentPage = 1;
        fetchHistory();
    }

    function toggleQueryType(type: string) {
        if (queryTypes.includes(type)) {
            queryTypes = queryTypes.filter(t => t !== type);
        } else {
            queryTypes = [...queryTypes, type];
        }
    }

    function formatDate(dateString: Date | string): string {
        return new Date(dateString).toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function toggleExpanded(index: number) {
        if (expandedCards.has(index)) {
            expandedCards.delete(index);
        } else {
            expandedCards.add(index);
        }
        expandedCards = expandedCards;
    }

    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
    }

    function formatJSON(jsonString: string): string {
        try {
            return JSON.stringify(JSON.parse(jsonString), null, 2);
        } catch {
            return jsonString;
        }
    }

    function truncateText(text: string, maxLength: number = 100): string {
        return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
    }

    function parseQueryData(history: QueryHistory) {
        try {
            const data = JSON.parse(history.query);
            return {
                userQuery: data.query || '',
                schema: data.schema || null,
                isSchemaQuery: !data.query || data.query === '',
                tableCount: data.schema?.tables?.length || 0
            };
        } catch {
            return {
                userQuery: history.query,
                schema: null,
                isSchemaQuery: false,
                tableCount: 0
            };
        }
    }

    function parseResultData(resultString: string) {
        try {
            const result = JSON.parse(resultString);

            if (result.tables && Array.isArray(result.tables)) {
                return {
                    type: 'schema',
                    tableCount: result.tables.length,
                    tables: result.tables.slice(0, 3).map(t => t.name)
                };
            }

            if (result.rows && Array.isArray(result.rows)) {
                return {
                    type: 'sql_result',
                    rowCount: result.rows.length,
                    columnCount: result.columns?.length || 0,
                    columns: result.columns?.slice(0, 4) || [],
                    sampleRows: result.rows.slice(0, 2)
                };
            }

            if (result.error) {
                return {
                    type: 'error',
                    message: result.error
                };
            }

            return {
                type: 'unknown',
                raw: result
            };
        } catch {
            return {
                type: 'raw',
                content: resultString
            };
        }
    }

    function getSortIcon(field: string) {
        const sortOption = sortOptions.find(s => s.field === field);
        if (!sortOption) return ArrowUpDown;
        return sortOption.direction === 'ASC' ? ArrowUp : ArrowDown;
    }

    function getSortPriority(field: string): number {
        const index = sortOptions.findIndex(s => s.field === field);
        return index >= 0 ? index + 1 : 0;
    }
</script>

<svelte:head>
    <title>Query History | Text-To-SQL</title>
</svelte:head>

<main class="flex-grow container mx-auto px-4 py-4 md:py-8 max-w-7xl">
    <!-- Mobile Header -->
    {#if isMobile}
        <div class="flex items-center justify-between mb-6">
            <div class="flex items-center gap-3">
                <div class="p-2 bg-blue-100 rounded-lg">
                    <History class="w-5 h-5 text-blue-600" />
                </div>
                <div>
                    <h1 class="text-xl font-bold text-gray-900">History</h1>
                    <p class="text-xs text-gray-500">{totalItems} queries</p>
                </div>
            </div>

            <button
                    on:click={() => showFilters = !showFilters}
                    class="p-2 bg-white border border-gray-300 rounded-lg"
            >
                <Menu class="w-5 h-5" />
            </button>
        </div>
    {:else}
        <!-- Desktop Header -->
        <div class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4 mb-8">
            <div class="flex items-center gap-3">
                <div class="p-2 bg-blue-100 rounded-lg">
                    <History class="w-6 h-6 text-blue-600" />
                </div>
                <div>
                    <h1 class="text-2xl font-bold text-gray-900">Query History</h1>
                    <p class="text-sm text-gray-500">Track and reuse your SQL queries</p>
                </div>
            </div>

            <div class="flex gap-2">
                <button
                        on:click={() => showFilters = !showFilters}
                        class="flex items-center gap-2 bg-white border border-gray-300 hover:bg-gray-50 text-gray-700 px-4 py-2 rounded-lg transition-colors"
                >
                    <Filter class="w-4 h-4" />
                    <span>Filters</span>
                    {#if databaseUUID || queryTypes.length > 0 || successFilter || startDate || endDate || searchQuery}
                        <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
                    {/if}
                </button>

                <button
                        on:click={fetchHistory}
                        class="flex items-center gap-2 bg-white border border-gray-300 hover:bg-gray-50 text-gray-700 px-4 py-2 rounded-lg transition-colors"
                >
                    <RefreshCcw class="w-4 h-4" />
                    <span>Refresh</span>
                </button>

                <button
                        on:click={clear}
                        class="flex items-center gap-2 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors"
                >
                    <Trash2 class="w-4 h-4" />
                    <span>Clear All</span>
                </button>
            </div>
        </div>
    {/if}

    <!-- Enhanced Filters Panel -->
    {#if showFilters}
        <div class="bg-white rounded-lg border border-gray-200 shadow-sm mb-6 overflow-hidden">
            <div class="p-4 border-b border-gray-100">
                <div class="flex items-center justify-between">
                    <h3 class="font-medium text-gray-900">Search & Filter</h3>
                    {#if isMobile}
                        <button
                                on:click={() => showFilters = false}
                                class="p-1 text-gray-400 hover:text-gray-600"
                        >
                            <X class="w-4 h-4" />
                        </button>
                    {/if}
                </div>
            </div>

            <div class="p-4 space-y-4">
                <!-- Search -->
                <div class="space-y-2">
                    <label class="block text-sm font-medium text-gray-700">Search in queries</label>
                    <div class="relative">
                        <input
                                type="text"
                                bind:value={searchQuery}
                                placeholder="Enter keywords..."
                                class="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        />
                        <Search class="w-4 h-4 text-gray-400 absolute left-3 top-3" />
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <!-- Query Types (Multiple Selection) -->
                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-gray-700">Query Types</label>
                        <div class="space-y-2">
                            {#each queryTypeOptions as type}
                                <label class="flex items-center space-x-2">
                                    <input
                                            type="checkbox"
                                            checked={queryTypes.includes(type)}
                                            on:change={() => toggleQueryType(type)}
                                            class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                                    />
                                    <span class="text-sm text-gray-700">{type}</span>
                                </label>
                            {/each}
                        </div>
                    </div>

                    <!-- Success Filter -->
                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-gray-700">Status</label>
                        <select
                                bind:value={successFilter}
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">All Status</option>
                            <option value="true">Success Only</option>
                            <option value="false">Failed Only</option>
                        </select>
                    </div>

                    <!-- Database -->
                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-gray-700">Database</label>
                        <select
                                bind:value={databaseUUID}
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">All Databases</option>
                            {#each databaseOptions as db}
                                <option value={db}>{db}</option>
                            {/each}
                        </select>
                    </div>

                    <!-- Date Range -->
                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-gray-700">From Date</label>
                        <input
                                type="datetime-local"
                                bind:value={startDate}
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        />
                    </div>

                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-gray-700">To Date</label>
                        <input
                                type="datetime-local"
                                bind:value={endDate}
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        />
                    </div>
                </div>

                <!-- Advanced Sorting -->
                <div class="border-t border-gray-200 pt-4">
                    <div class="flex items-center justify-between mb-3">
                        <label class="block text-sm font-medium text-gray-700">Sorting Options</label>
                        <button
                                on:click={() => showAdvancedSort = !showAdvancedSort}
                                class="text-sm text-blue-600 hover:text-blue-700"
                        >
                            {showAdvancedSort ? 'Simple' : 'Advanced'}
                        </button>
                    </div>

                    {#if showAdvancedSort}
                        <div class="space-y-3">
                            {#each sortOptions as sortOption, index}
                                <div class="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
                                    <span class="text-xs font-medium text-gray-500 w-6">#{index + 1}</span>

                                    <select
                                            bind:value={sortOption.field}
                                            on:change={() => updateSortOption(index, sortOption.field, sortOption.direction)}
                                            class="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                    >
                                        <option value="timestamp">Date</option>
                                        <option value="query_type">Type</option>
                                        <option value="success">Status</option>
                                    </select>

                                    <select
                                            bind:value={sortOption.direction}
                                            on:change={() => updateSortOption(index, sortOption.field, sortOption.direction)}
                                            class="px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                    >
                                        <option value="DESC">↓ Desc</option>
                                        <option value="ASC">↑ Asc</option>
                                    </select>

                                    {#if sortOptions.length > 1}
                                        <button
                                                on:click={() => removeSortOption(index)}
                                                class="p-2 text-red-500 hover:bg-red-50 rounded-lg"
                                        >
                                            <X class="w-4 h-4" />
                                        </button>
                                    {/if}
                                </div>
                            {/each}

                            {#if sortOptions.length < 3}
                                <button
                                        on:click={addSortOption}
                                        class="flex items-center gap-2 w-full p-3 text-sm text-gray-600 border-2 border-dashed border-gray-300 rounded-lg hover:border-gray-400 hover:bg-gray-50 transition-colors"
                                >
                                    <Plus class="w-4 h-4" />
                                    Add sorting criteria
                                </button>
                            {/if}
                        </div>
                    {:else}
                        <div class="flex flex-wrap gap-2">
                            {#each ['timestamp', 'query_type', 'success'] as field}
                                {@const priority = getSortPriority(field)}
                                {@const Icon = getSortIcon(field)}
                                <button
                                        on:click={() => quickSort(field)}
                                        class="flex items-center gap-1 px-3 py-1.5 text-sm rounded-lg transition-colors relative {priority > 0 ? 'bg-blue-100 text-blue-700 border border-blue-200' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}"
                                >
                                    <Icon class="w-3 h-3" />
                                    {field === 'timestamp' ? 'Date' : field === 'query_type' ? 'Type' : 'Status'}
                                    {#if priority > 0}
                                        <span class="absolute -top-1 -right-1 w-4 h-4 bg-blue-500 text-white text-xs rounded-full flex items-center justify-center">
                                            {priority}
                                        </span>
                                    {/if}
                                </button>
                            {/each}
                        </div>
                    {/if}
                </div>

                <div class="flex flex-col sm:flex-row gap-3 pt-4 border-t border-gray-200">
                    <button
                            on:click={applyFilters}
                            class="flex-1 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
                    >
                        Apply Filters
                    </button>
                    <button
                            on:click={resetFilters}
                            class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg transition-colors"
                    >
                        Reset All
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- Results Info -->
    {#if !loading && histories.length > 0}
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-4">
            <div class="text-sm text-gray-600">
                Showing {(currentPage - 1) * perPage + 1}-{Math.min(currentPage * perPage, totalItems)} of {totalItems} results
                {#if sortOptions.length > 1}
                    <span class="text-xs text-gray-500">
                        • Sorted by {sortOptions.length} criteria
                    </span>
                {/if}
            </div>

            <!-- Active Filters Summary -->
            {#if databaseUUID || queryTypes.length > 0 || successFilter || startDate || endDate || searchQuery}
                <div class="flex flex-wrap gap-1">
                    {#if searchQuery}
                        <span class="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                            Search: "{searchQuery}"
                            <button on:click={() => { searchQuery = ''; applyFilters(); }}>
                                <X class="w-3 h-3" />
                            </button>
                        </span>
                    {/if}
                    {#each queryTypes as type}
                        <span class="inline-flex items-center gap-1 px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">
                            {type}
                            <button on:click={() => toggleQueryType(type)}>
                                <X class="w-3 h-3" />
                            </button>
                        </span>
                    {/each}
                    {#if successFilter}
                        <span class="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-full">
                            {successFilter === 'true' ? 'Success' : 'Failed'}
                            <button on:click={() => { successFilter = ''; applyFilters(); }}>
                                <X class="w-3 h-3" />
                            </button>
                        </span>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}

    <!-- Content -->
    {#if loading}
        <div class="flex justify-center py-16">
            <div class="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent"></div>
        </div>
    {:else if error}
        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
            <div class="flex items-center gap-2 text-red-800 mb-1">
                <XCircle class="w-5 h-5" />
                <span class="font-medium">Error occurred</span>
            </div>
            <p class="text-red-700 text-sm">{error}</p>
        </div>
    {:else if histories.length === 0}
        <div class="text-center py-16">
            <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <History class="w-8 h-8 text-gray-400" />
            </div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No Query History</h3>
            <p class="text-gray-500 max-w-md mx-auto">
                {searchQuery || databaseUUID || queryTypes.length > 0 || successFilter || startDate || endDate
                    ? "No results match your current filters. Try adjusting your search criteria."
                    : "You haven't executed any queries yet. Start generating SQL queries to see them here."}
            </p>
        </div>
    {:else}
        <!-- History Cards -->
        <div class="space-y-3 md:space-y-4">
            {#each histories as history, index}
                <div class="bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow">
                    <!-- Card Header -->
                    <div class="p-3 md:p-4 border-b border-gray-100">
                        <div class="flex items-center justify-between gap-3">
                            <div class="flex items-center gap-3 min-w-0 flex-1">
                                <div class="flex items-center gap-2 flex-shrink-0">
                                    {#if history.success}
                                        <CheckCircle class="w-4 h-4 md:w-5 md:h-5 text-green-500" />
                                        <span class="text-xs md:text-sm font-medium text-green-700">Success</span>
                                    {:else}
                                        <XCircle class="w-4 h-4 md:w-5 md:h-5 text-red-500" />
                                        <span class="text-xs md:text-sm font-medium text-red-700">Failed</span>
                                    {/if}
                                </div>

                                {#if !isMobile}
                                    <div class="h-4 w-px bg-gray-300 flex-shrink-0"></div>

                                    <div class="flex items-center gap-4 text-xs md:text-sm text-gray-600 min-w-0">
                                        <span class="flex items-center gap-1 flex-shrink-0">
                                            <Tag class="w-3 h-3 md:w-4 md:h-4" />
                                            {history.query_type}
                                        </span>
                                        <span class="flex items-center gap-1 flex-shrink-0">
                                            <Database class="w-3 h-3 md:w-4 md:h-4" />
                                            {history.database_uuid || 'Default DB'}
                                        </span>
                                        <span class="flex items-center gap-1 flex-shrink-0">
                                            <Calendar class="w-3 h-3 md:w-4 md:h-4" />
                                            {formatDate(history.timestamp)}
                                        </span>
                                    </div>
                                {/if}
                            </div>

                            <div class="flex items-center gap-1 md:gap-2 flex-shrink-0">
                                <button
                                        on:click={() => toggleExpanded(index)}
                                        class="flex items-center gap-1 px-2 md:px-3 py-1 md:py-1.5 text-xs md:text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                                >
                                    {#if isMobile}
                                        {#if expandedCards.has(index)}
                                            <EyeOff class="w-4 h-4" />
                                        {:else}
                                            <Eye class="w-4 h-4" />
                                        {/if}
                                    {:else}
                                        {#if expandedCards.has(index)}
                                            <EyeOff class="w-4 h-4" />
                                            Hide Details
                                        {:else}
                                            <Eye class="w-4 h-4" />
                                            Show Details
                                        {/if}
                                    {/if}
                                </button>

                                <button class="flex items-center gap-1 px-2 md:px-3 py-1 md:py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs md:text-sm rounded-lg transition-colors">
                                    <Code class="w-3 h-3 md:w-4 md:h-4" />
                                    {#if !isMobile}Reuse{/if}
                                </button>
                            </div>
                        </div>

                        <!-- Mobile-specific info on mobile -->
                        {#if isMobile}
                            <div class="mt-3 pt-3 border-t border-gray-100">
                                <div class="flex flex-wrap items-center gap-3 text-xs text-gray-600">
                                    <span class="flex items-center gap-1">
                                        <Tag class="w-3 h-3" />
                                        {history.query_type}
                                    </span>
                                    <span class="flex items-center gap-1">
                                        <Database class="w-3 h-3" />
                                        {history.database_uuid ? history.database_uuid.slice(0, 8) + '...' : 'Default DB'}
                                    </span>
                                    <span class="flex items-center gap-1">
                                        <Clock class="w-3 h-3" />
                                        {formatDate(history.timestamp)}
                                    </span>
                                </div>
                            </div>
                        {/if}
                    </div>

                    <!-- Query Content -->
                    {#if expandedCards.has(index)}
                        {@const queryData = parseQueryData(history)}
                        {@const resultData = parseResultData(history.result)}
                        <div class="p-3 md:p-4 border-t border-gray-100">
                            <div class="space-y-4">
                                <!-- User Query -->
                                <div>
                                    <div class="flex items-center justify-between mb-2">
                                        <h4 class="text-sm font-medium text-gray-900">
                                            {queryData.isSchemaQuery ? 'Database Schema Request' : 'SQL Query'}
                                        </h4>
                                        <button
                                                on:click={() => copyToClipboard(queryData.userQuery || 'SHOW TABLES')}
                                                class="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
                                                title="Copy query"
                                        >
                                            <Copy class="w-3.5 h-3.5" />
                                        </button>
                                    </div>

                                    <div class="bg-slate-50 border border-slate-200 rounded-lg p-3">
                                        {#if queryData.isSchemaQuery}
                                            <div class="flex items-center gap-2 text-sm text-slate-600">
                                                <Database class="w-4 h-4" />
                                                <span>Requested database schema information</span>
                                                {#if queryData.tableCount > 0}
                                                    <span class="text-slate-400">•</span>
                                                    <span>{queryData.tableCount} tables found</span>
                                                {/if}
                                            </div>
                                        {:else}
                                            <code class="text-sm text-slate-800 font-mono block break-all">
                                                {queryData.userQuery}
                                            </code>
                                        {/if}
                                    </div>
                                </div>

                                <!-- Result Summary -->
                                <div>
                                    <h4 class="text-sm font-medium text-gray-900 mb-2">Result</h4>

                                    {#if resultData.type === 'schema'}
                                        <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
                                            <div class="flex items-center gap-2 text-blue-800 mb-2">
                                                <Database class="w-4 h-4" />
                                                <span class="font-medium">Schema Retrieved</span>
                                            </div>
                                            <div class="text-sm text-blue-700">
                                                <p>Found {resultData.tableCount} tables</p>
                                                {#if resultData.tables.length > 0}
                                                    <div class="mt-2">
                                                        <p class="text-xs text-blue-600 mb-1">Tables:</p>
                                                        <div class="flex flex-wrap gap-1">
                                                            {#each resultData.tables as table}
                                                                <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">{table}</span>
                                                            {/each}
                                                            {#if resultData.tableCount > 3}
                                                                <span class="px-2 py-1 bg-blue-100 text-blue-600 text-xs rounded">+{resultData.tableCount - 3} more</span>
                                                            {/if}
                                                        </div>
                                                    </div>
                                                {/if}
                                            </div>
                                        </div>
                                    {:else if resultData.type === 'sql_result'}
                                        <div class="bg-green-50 border border-green-200 rounded-lg p-3">
                                            <div class="flex items-center gap-2 text-green-800 mb-2">
                                                <CheckCircle class="w-4 h-4" />
                                                <span class="font-medium">Query Executed Successfully</span>
                                            </div>
                                            <div class="text-sm text-green-700">
                                                <p>{resultData.rowCount} rows returned</p>
                                                {#if resultData.columnCount > 0}
                                                    <p class="text-xs text-green-600 mt-1">Columns: {resultData.columns.join(', ')}</p>
                                                {/if}
                                            </div>

                                            {#if resultData.sampleRows.length > 0}
                                                <div class="mt-3 overflow-x-auto">
                                                    <table class="min-w-full text-xs">
                                                        <thead>
                                                        <tr class="border-b border-green-200">
                                                            {#each resultData.columns as column}
                                                                <th class="text-left py-1 px-2 font-medium text-green-800">{column}</th>
                                                            {/each}
                                                        </tr>
                                                        </thead>
                                                        <tbody>
                                                        {#each resultData.sampleRows as row}
                                                            <tr class="border-b border-green-100">
                                                                {#each resultData.columns as column}
                                                                    <td class="py-1 px-2 text-green-700">{row[column] || '-'}</td>
                                                                {/each}
                                                            </tr>
                                                        {/each}
                                                        </tbody>
                                                    </table>
                                                    {#if resultData.rowCount > 2}
                                                        <p class="text-xs text-green-600 mt-2">... and {resultData.rowCount - 2} more rows</p>
                                                    {/if}
                                                </div>
                                            {/if}
                                        </div>
                                    {:else if resultData.type === 'error'}
                                        <div class="bg-red-50 border border-red-200 rounded-lg p-3">
                                            <div class="flex items-center gap-2 text-red-800 mb-2">
                                                <XCircle class="w-4 h-4" />
                                                <span class="font-medium">Query Failed</span>
                                            </div>
                                            <div class="text-sm text-red-700">
                                                <code class="bg-red-100 px-2 py-1 rounded text-xs break-all">{resultData.message}</code>
                                            </div>
                                        </div>
                                    {:else}
                                        <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
                                            <div class="text-sm text-gray-600">
                                                <pre class="text-xs font-mono whitespace-pre-wrap break-all">{formatJSON(history.result)}</pre>
                                            </div>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>

        <!-- Pagination -->
        {#if totalPages > 1}
            <div class="mt-6 md:mt-8">
                {#if isMobile}
                    <!-- Mobile Pagination -->
                    <div class="flex items-center justify-between">
                        <button
                                class="flex items-center gap-2 px-3 py-2 text-sm bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                disabled={currentPage === 1}
                                on:click={() => changePage(currentPage - 1)}
                        >
                            <ChevronLeft class="w-4 h-4" />
                            Previous
                        </button>

                        <div class="flex items-center gap-1 text-sm text-gray-600">
                            <span>{currentPage}</span>
                            <span class="text-gray-400">of</span>
                            <span>{totalPages}</span>
                        </div>

                        <button
                                class="flex items-center gap-2 px-3 py-2 text-sm bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                disabled={currentPage === totalPages}
                                on:click={() => changePage(currentPage + 1)}
                        >
                            Next
                            <ChevronRight class="w-4 h-4" />
                        </button>
                    </div>
                {:else}
                    <!-- Desktop Pagination -->
                    <div class="flex justify-center">
                        <div class="flex items-center gap-1">
                            <button
                                    class="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                    disabled={currentPage === 1}
                                    on:click={() => changePage(currentPage - 1)}
                            >
                                <ChevronLeft class="w-5 h-5" />
                            </button>

                            {#if totalPages <= 7}
                                {#each Array(totalPages) as _, i}
                                    <button
                                            class="w-10 h-10 rounded-lg text-sm font-medium transition-colors {currentPage === i + 1 ? 'bg-blue-600 text-white' : 'border border-gray-300 hover:bg-gray-50'}"
                                            on:click={() => changePage(i + 1)}
                                    >
                                        {i + 1}
                                    </button>
                                {/each}
                            {:else}
                                <button
                                        class="w-10 h-10 rounded-lg text-sm font-medium transition-colors {currentPage === 1 ? 'bg-blue-600 text-white' : 'border border-gray-300 hover:bg-gray-50'}"
                                        on:click={() => changePage(1)}
                                >
                                    1
                                </button>

                                {#if currentPage > 3}
                                    <span class="px-2 text-gray-500">...</span>
                                {/if}

                                {#each Array(3) as _, i}
                                    {@const pageNum = currentPage - 1 + i}
                                    {#if pageNum > 1 && pageNum < totalPages}
                                        <button
                                                class="w-10 h-10 rounded-lg text-sm font-medium transition-colors {currentPage === pageNum ? 'bg-blue-600 text-white' : 'border border-gray-300 hover:bg-gray-50'}"
                                                on:click={() => changePage(pageNum)}
                                        >
                                            {pageNum}
                                        </button>
                                    {/if}
                                {/each}

                                {#if currentPage < totalPages - 2}
                                    <span class="px-2 text-gray-500">...</span>
                                {/if}

                                <button
                                        class="w-10 h-10 rounded-lg text-sm font-medium transition-colors {currentPage === totalPages ? 'bg-blue-600 text-white' : 'border border-gray-300 hover:bg-gray-50'}"
                                        on:click={() => changePage(totalPages)}
                                >
                                    {totalPages}
                                </button>
                            {/if}

                            <button
                                    class="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                    disabled={currentPage === totalPages}
                                    on:click={() => changePage(currentPage + 1)}
                            >
                                <ChevronRight class="w-5 h-5" />
                            </button>
                        </div>
                    </div>
                {/if}
            </div>
        {/if}
    {/if}
</main>