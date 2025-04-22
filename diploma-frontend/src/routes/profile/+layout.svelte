<script lang="ts">
    import { onMount } from 'svelte';
    import { userStore, loadingUser, initializeUser } from '$lib/stores';
    import { get } from 'svelte/store';
    import { goto } from '$app/navigation';

    onMount(async () => {
        await initializeUser();

        const user = get(userStore);
        if (!user) {
            goto('/auth');
        }
    });
    import Navbar from "$lib/components/Navbar.svelte";
</script>

<Navbar />


{#if $loadingUser}
    <p>Loading...</p>
{:else}
    <slot />
{/if}

