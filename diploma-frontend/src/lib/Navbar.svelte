<script>
    import { isAuthenticated } from '$lib/stores/authStore';
    import { page } from '$app/stores';
    import { derived } from 'svelte/store';
    import { onMount } from "svelte";
    import { authorization } from "./api.ts";
    import { goto } from "$app/navigation";

    const currentRoute = derived(page, ($page) => $page.url.pathname);

    async function logout() {
        await authorization.logout();
        await goto('/');
    }

    onMount(() => {
        import('$lib/stores/authStore').then((module) => {
            module.checkAuthStatus();
        });
    })
</script>

<header class="bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-md">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <a href="/" class="text-2xl font-bold hover:opacity-90">MyApp</a>

        <div class="flex space-x-6">
            {#if $isAuthenticated}
                <a
                        href="/query"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/query' && 'underline font-bold'}
                >
                    Query
                </a>
                <a
                        href="/items"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/items' && 'underline font-bold'}
                >
                    Items
                </a>
                <a
                        on:click={logout}
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/auth' && 'underline font-bold'}
                >
                    Logout
                </a>
            {:else }
                <a
                        href="/auth"
                        class="hover:text-purple-200 font-medium text-lg transition-colors"
                        class:selected={$currentRoute === '/auth' && 'underline font-bold'}
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
