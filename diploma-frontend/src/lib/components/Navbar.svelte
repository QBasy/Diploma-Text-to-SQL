<script lang="ts">
    import { page } from '$app/stores';
    import { logoutUser, userStore } from "$lib/stores/";
    import { goto } from "$app/navigation";
    import { derived } from 'svelte/store';

    // Derived store to check login status
    const userLogIn = derived(userStore, $userStore => $userStore !== null);

    let showDropdown: boolean = false;

    function handleLogout(): void {
        logoutUser();
        setTimeout(() => {
            goto("/");
        }, 500);
    }

    function toggleDropdown(): void {
        showDropdown = !showDropdown;
    }

    function closeDropdown(): void {
        showDropdown = false;
    }
</script>

<nav class="bg-blue-500 text-white p-4 flex justify-between items-center">
    <div class="text-xl font-bold">SQL Generator</div>

    <div class="flex space-x-4 mx-12">
        <a href="/" class:active={$page.url.pathname === '/'}>Home</a>
        <a href="/generate/simple" class:active={$page.url.pathname === '/generate/simple'}>Generate SQL</a>
        <a href="/generate/complex" class:active={$page.url.pathname === '/generate/complex'}>Generate Complex SQL</a>

        {#if $userLogIn}
            <div class="relative">
                <button
                        on:click={toggleDropdown}
                        class="hover:underline focus:outline-none"
                >
                    Profile
                </button>
                {#if showDropdown}
                    <div id="dropdown" class="absolute top-full right-0 bg-white text-black shadow-lg rounded mt-1">
                        <a href="/profile" class="block px-4 py-2 hover:bg-gray-200">User Profile</a><hr>
                        <a href="/profile/database" class="block px-4 py-2 hover:bg-gray-200">Profile Database</a><hr>
                        <a href="/profile/settings" class="block px-4 py-2 hover:bg-gray-200">Profile Settings</a><hr>
                        <a href="/profile/history" class="block px-4 py-2 hover:bg-gray-200">Profile History</a><hr>
                        <button on:click={handleLogout} class="block w-full text-left px-4 py-2 hover:bg-gray-200">Logout</button>
                    </div>
                {/if}
            </div>
        {:else}
            <a href="/auth">Login/Register</a>
        {/if}
    </div>
</nav>

<style>
    .active {
        font-weight: bold;
        text-decoration: underline;
    }

    button:focus {
        outline: none;
    }
</style>