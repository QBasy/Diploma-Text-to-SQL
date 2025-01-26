<script lang="ts">
    import { page } from '$app/stores';
    import { logoutUser, userStore } from "$lib/stores/index.js";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    let userLogIn: boolean | null = null;
    let showDropdown: boolean = false;

    function handleLogout(): void {
        logoutUser();
        setTimeout(() => {
            goto("/");
        }, 500);
    }

    function toggleDropdown(visible: boolean): void {
        showDropdown = visible;
    }

    onMount(() => {
        userStore.subscribe((user) => {
            if (!user) {
                userLogIn = false;
            } else {
                userLogIn = true;
            }
        });
    });
</script>

<nav class="bg-blue-500 text-white p-4 flex justify-between items-center">
    <div class="text-xl font-bold">SQL Generator</div>

    <div class="flex space-x-4">
        <a href="/" class:active={$page.url.pathname === '/'}>Home</a>
        <a href="/generate/simple" class:active={$page.url.pathname === '/generate/simple'}>Generate SQL</a>
        <a href="/generate/complex" class:active={$page.url.pathname === '/generate/complex'}>Generate Complex SQL</a>

        {#if userLogIn}
            <div
                    class="relative"
                    on:mouseenter={() => toggleDropdown(true)}
                    on:mouseleave={() => toggleDropdown(false)}
                    on:focusin={() => toggleDropdown(true)}
                    on:focusout={() => toggleDropdown(false)}
                    role="button"
                    tabindex="0"
            >
                <a href="/profile" class="hover:underline">Profile</a>
                {#if showDropdown}
                    <div class="absolute top-full left-0 bg-white text-black shadow-lg rounded mt-1">
                        <a href="/profile/database" class="block px-4 py-2 hover:bg-gray-200">Profile Database</a>
                        <a href="/profile/settings" class="block px-4 py-2 hover:bg-gray-200">Profile Settings</a>
                        <a href="/profile/history" class="block px-4 py-2 hover:bg-gray-200">Profile History</a>
                        <button on:click={handleLogout} class="block w-full text-left px-4 py-2 hover:bg-gray-200">Logout</button>
                    </div>
                {/if}
            </div>
        {:else}
            <a href="/auth">Login</a>
        {/if}
    </div>
</nav>

<style>
    .active {
        font-weight: bold;
        text-decoration: underline;
    }
</style>