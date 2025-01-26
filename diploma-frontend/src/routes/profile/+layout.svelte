<script>
    import { onMount } from "svelte";
    import { userStore } from "$lib/stores/index.ts";
    import { goto } from "$app/navigation";
    import Navbar from "$lib/components/Navbar.svelte";

    onMount(() => {
        userStore.subscribe((user) => {
            if (!user) {
                goto('/auth');
            }
        });
    });
</script>

<Navbar />