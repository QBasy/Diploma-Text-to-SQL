<script lang="ts">
    import Footer from "$lib/components/Footer.svelte";
</script>

<div class="min-h-screen flex flex-col">

    <main class="flex-1">
        <slot />
    </main>

    <Footer />
</div>
