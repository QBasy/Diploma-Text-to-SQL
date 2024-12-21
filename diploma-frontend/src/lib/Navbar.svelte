<script lang="ts">
    // @ts-ignore
    import { isAuthenticated, checkAuthStatus } from '$lib/stores/authStore';
    import { page } from '$app/stores';
    import { derived } from 'svelte/store';
    // @ts-ignore
    import { authorization } from '$lib/api';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    const currentRoute = derived(page, ($page) => $page.url.pathname);

    onMount(() => {
        checkAuthStatus();
    });

    async function logout() {
        await authorization.logout();
        checkAuthStatus();
        goto('/');
    }
</script>

<header class="bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-md">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <a href="/" class="text-2xl font-bold hover:opacity-90">MyApp</a>

        <div class="flex space-x-6">
            {#if $isAuthenticated}
                <a
                        href="/query"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/query'}
                >
                    Query
                </a>
                <a
                        href="/database"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/database'}
                >
                    Database
                </a>
                <a
                        on:click={logout}
                        class="cursor-pointer hover:text-purple-200 font-medium text-lg transition-colors"
                >
                    Logout
                </a>
            {:else}
                <a
                        href="/auth"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/auth'}
                >
                    Auth
                </a>
            {/if}
        </div>
    </div>
</header>

<style>
    a.selected {
        text-decoration: underline;
        font-weight: bold;
    }
</style>
