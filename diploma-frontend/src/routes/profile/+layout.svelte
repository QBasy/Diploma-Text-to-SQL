<script lang="ts">
    import { onMount } from "svelte";
    import { userStore } from "$lib/stores/";
    import { goto } from "$app/navigation";
    import Navbar from "$lib/components/Navbar.svelte";

    let { children } = $props();

    onMount(() => {
        userStore.subscribe(user => {
            if (!user) {
                goto('/auth');
            }
        });
    });
</script>

<Navbar />

{@render children()}