<script lang="ts">
    import { page } from '$app/stores';
    import { logoutUser, userStore } from '$lib/stores/';
    import { goto } from '$app/navigation';
    import { derived } from 'svelte/store';
    import { UserCircle, LogOut } from 'lucide-svelte';

    const userLogIn = derived(userStore, $userStore => $userStore !== null);

    let showDropdown = false;

    function handleLogout(): void {
        logoutUser();
        setTimeout(() => goto('/'), 500);
    }

    function toggleDropdown(): void {
        showDropdown = !showDropdown;
    }

    function closeDropdown(): void {
        showDropdown = false;
    }
</script>

<nav class="bg-white border-b border-gray-200 shadow-sm">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <a href="/" class="text-2xl font-bold text-blue-600 tracking-tight">SQL Generator / Query Visualisation</a>

        <div class="flex items-center space-x-6 text-gray-700 font-medium">
            <a href="/" class="hover:text-blue-600 transition" class:font-bold={$page.url.pathname === '/'}>Home</a>
            <a href="/generate/simple" class="hover:text-blue-600 transition" class:font-bold={$page.url.pathname === '/generate/simple'}>Simple</a>
            <a href="/generate/complex" class="hover:text-blue-600 transition" class:font-bold={$page.url.pathname === '/generate/complex'}>Complex</a>

            {#if $userLogIn}
                <div class="relative">
                    <button on:click={toggleDropdown} class="flex items-center gap-1 hover:text-blue-600 transition">
                        <UserCircle class="w-5 h-5" /> Profile
                    </button>
                    {#if showDropdown}
                        <div class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-xl border z-50">
                            <a href="/profile" class="block px-4 py-2 hover:bg-gray-100">User Profile</a>
                            <a href="/profile/visualisation" class="block px-4 py-2 hover:bg-gray-100">Visualisation</a>
                            <a href="/profile/database" class="block px-4 py-2 hover:bg-gray-100">Database</a>
                            <a href="/profile/settings" class="block px-4 py-2 hover:bg-gray-100">Settings</a>
                            <a href="/profile/history" class="block px-4 py-2 hover:bg-gray-100">History</a>
                            <hr class="my-1" />
                            <button on:click={handleLogout} class="flex items-center gap-2 px-4 py-2 text-red-600 hover:bg-red-50 w-full text-left">
                                <LogOut class="w-4 h-4" /> Logout
                            </button>
                        </div>
                    {/if}
                </div>
            {:else}
                <a href="/auth" class="hover:text-blue-600 transition">Login/Register</a>
            {/if}
        </div>
    </div>
</nav>
