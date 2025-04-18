<script>
    import NavbarDocumentation from "$lib/components/NavbarDocumentation.svelte";
</script>

<NavbarDocumentation />
<section class="flex-grow container mx-auto px-4 py-16">
    <h1 class="text-4xl font-bold text-center text-gray-800 mb-12">SQL Query Generation</h1>

    <div class="grid md:grid-cols-2 gap-8">
        <div class="bg-white p-8 rounded-2xl shadow-lg hover:shadow-xl transition">
            <h2 class="text-2xl font-semibold text-blue-600 mb-4">Simple Generation</h2>
            <p class="text-gray-700 mb-6">
                Quick SQL generation for straightforward, natural language queries.
            </p>
            <a href="/generate/simple" class="inline-block bg-blue-600 text-white px-5 py-2 rounded-lg hover:bg-blue-700 transition">
                Try Now
            </a>
        </div>

        <div class="bg-white p-8 rounded-2xl shadow-lg hover:shadow-xl transition">
            <h2 class="text-2xl font-semibold text-green-600 mb-4">Complex Generation</h2>
            <p class="text-gray-700 mb-6">
                Advanced SQL generation with full schema understanding and context.
            </p>
            <a href="/generate/complex" class="inline-block bg-green-600 text-white px-5 py-2 rounded-lg hover:bg-green-700 transition">
                Explore
            </a>
        </div>
    </div>
</section>
