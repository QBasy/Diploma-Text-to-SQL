<script>
    import Navbar from '$lib/components/Navbar.svelte';
    import Footer from '$lib/components/Footer.svelte';
</script>

<div class="min-h-screen flex flex-col">
    <Navbar />

    <main class="flex-grow container mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold mb-4 text-blue-500">SQL Generator</h1>

        <p class="mb-6">
            Transform natural language queries into precise SQL statements with ease.
        </p>

        <div class="grid md:grid-cols-3 gap-4">
            <div class="bg-white p-6 rounded-lg shadow-md">
                <h2 class="text-xl font-semibold mb-2">Simple Generation</h2>
                <p>Quick SQL generation for straightforward queries.</p>
                <a href="/generate/simple" class="btn bg-blue-500 text-white mt-4">Try Now</a>
            </div>

            <div class="bg-white p-6 rounded-lg shadow-md">
                <h2 class="text-xl font-semibold mb-2">Complex Generation</h2>
                <p>Advanced SQL generation with database schema context.</p>
                <a href="/generate/complex" class="btn bg-green-500 text-white mt-4">Explore</a>
            </div>

            <div class="bg-white p-6 rounded-lg shadow-md">
                <h2 class="text-xl font-semibold mb-2">Documentation</h2>
                <p>See example queries and generated SQL.</p>
                <a href="/documentation" class="btn bg-gray-500 text-white mt-4">View Documentation</a>
            </div>
        </div>
    </main>

    <Footer />
</div>

<style>
    .btn {
        @apply inline-block px-4 py-2 rounded hover:opacity-80 transition-all;
    }
</style>