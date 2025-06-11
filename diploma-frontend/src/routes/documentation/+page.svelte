<script lang="ts">
    import { BarChart3 } from 'lucide-svelte';
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">📘 System Documentation</h1>
<p class="text-gray-600 mb-8">
    Welcome to the documentation for the automatic data visualization system. This section provides guidance on setup, architecture, and development.
</p>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <a href="/documentation/charts">
    <div class="p-6 bg-white rounded-2xl shadow hover:shadow-md transition">
        <div class="flex items-center gap-3 mb-3 text-purple-600">
            <BarChart3 class="w-5 h-5" />
            <h2 class="text-xl font-semibold">Chart Types</h2>
        </div>
        <p class="text-gray-600 text-sm">
            Overview of supported SVG chart types: bar, pie, line, and scatter — all rendered on backend without JS.
        </p>
    </div></a>
</div>
