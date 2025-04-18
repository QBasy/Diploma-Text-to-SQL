<script lang="ts">
	import '../app.css';
    import {initializeUser} from "$lib/stores";
    import {onMount} from "svelte";
	let { children } = $props();
    onMount(async () => {
        await initializeUser();
    });
</script>

{@render children()}
