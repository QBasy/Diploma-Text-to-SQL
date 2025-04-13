<script lang="ts">
    import { onMount } from "svelte";
    import { userStore } from "$lib/stores/";
    import { goto } from "$app/navigation";
    import Navbar from "$lib/components/Navbar.svelte";

    let { children } = $props();

    onMount(() => {
        const unsubscribe = userStore.subscribe(user => {
            if (!user) {
                goto('/auth');
            }
        });

        return () => unsubscribe();
    });
</script>

<Navbar />

{@render children()}
