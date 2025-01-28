<script>
import Navbar from "$lib/components/Navbar.svelte";
import Footer from "$lib/components/Footer.svelte";

let { children } = $props();

</script>

<div class="min-h-screen flex flex-col">
    <Navbar />
    {@render children()}
    <Footer />
</div>