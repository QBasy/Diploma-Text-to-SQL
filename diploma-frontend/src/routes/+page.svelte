<script>
    import Navbar from '$lib/components/Navbar.svelte';
    import Footer from '$lib/components/Footer.svelte';
    import { Rocket, Sparkles, BookText } from 'lucide-svelte';
</script>

<svelte:head>
    <title>Text-To-SQL / Visualisation</title>
</svelte:head>

<div class="min-h-screen flex flex-col bg-gradient-to-b from-white to-slate-100">
    <Navbar />

    <main class="flex-grow container mx-auto px-4 py-16">
        <h1 class="text-4xl md:text-5xl font-bold mb-4 text-center text-blue-600 tracking-tight">
            Transform Natural Language into SQL
        </h1>

        <h1 class="text-3xl md:text-4xl font-bold mb-4 text-center text-blue-600 tracking-tight">
            And Visualise Result
        </h1>

        <p class="mb-12 text-center text-lg text-gray-600 max-w-2xl mx-auto">
            Effortlessly convert human language queries into structured SQL commands and generate instant visualizations.
        </p>

        <div class="grid md:grid-cols-3 gap-6 max-w-6xl mx-auto">
            <div class="bg-white p-6 rounded-2xl shadow-md border hover:shadow-lg transition-all group">
                <div class="flex items-center gap-3 mb-4 text-blue-600">
                    <Rocket class="w-6 h-6 group-hover:scale-110 transition" />
                    <h2 class="text-xl font-semibold">Simple Generation</h2>
                </div>
                <p class="text-gray-600">Quick SQL generation for straightforward queries.</p>
                <a href="/generate/simple" class="mt-6 inline-block bg-blue-600 text-white px-4 py-2 rounded-xl hover:bg-blue-700 transition">Try Now</a>
            </div>

            <div class="bg-white p-6 rounded-2xl shadow-md border hover:shadow-lg transition-all group">
                <div class="flex items-center gap-3 mb-4 text-green-600">
                    <Sparkles class="w-6 h-6 group-hover:scale-110 transition" />
                    <h2 class="text-xl font-semibold">Complex Generation</h2>
                </div>
                <p class="text-gray-600">Advanced SQL generation with schema context and filters.</p>
                <a href="/generate/complex" class="mt-6 inline-block bg-green-600 text-white px-4 py-2 rounded-xl hover:bg-green-700 transition">Explore</a>
            </div>

            <div class="bg-white p-6 rounded-2xl shadow-md border hover:shadow-lg transition-all group">
                <div class="flex items-center gap-3 mb-4 text-gray-700">
                    <BookText class="w-6 h-6 group-hover:scale-110 transition" />
                    <h2 class="text-xl font-semibold">Documentation</h2>
                </div>
                <p class="text-gray-600">Explore example queries and generated SQL statements.</p>
                <a href="/documentation" class="mt-6 inline-block bg-gray-700 text-white px-4 py-2 rounded-xl hover:bg-gray-800 transition">View Docs</a>
            </div>
        </div>
    </main>

    <Footer />
</div>
