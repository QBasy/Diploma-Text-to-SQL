<script lang="ts">
    import { Code2, Database, PlayCircle, BarChart3, Settings } from 'lucide-svelte';
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">📘 System Documentation</h1>
<p class="text-gray-600 mb-8">
    Welcome to the documentation for the automatic data visualization system. This section provides guidance on setup, architecture, and development.
</p>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <a href="/documentation/start">
    <div class="p-6 bg-white rounded-2xl shadow hover:shadow-md transition">
        <div class="flex items-center gap-3 mb-3 text-blue-600">
            <PlayCircle class="w-5 h-5" />
            <h2 class="text-xl font-semibold">How to Start</h2>
        </div>
        <p class="text-gray-600 text-sm">
            Step-by-step guide to clone, run, and test the entire system locally using Docker Compose and environment variables.
        </p>
    </div></a>
    <a href="/documentation/stack">
    <div class="p-6 bg-white rounded-2xl shadow hover:shadow-md transition">
        <div class="flex items-center gap-3 mb-3 text-green-600">
            <Code2 class="w-5 h-5" />
            <h2 class="text-xl font-semibold">Tech Stack</h2>
        </div>
        <p class="text-gray-600 text-sm">
            Go, FastAPI, SvelteKit, PostgreSQL, SQLite, Docker, gRPC, and PyTorch — detailed breakdown of each component.
        </p>
    </div></a>
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
    <a href="/documentation/database">
    <div class="p-6 bg-white rounded-2xl shadow hover:shadow-md transition">
        <div class="flex items-center gap-3 mb-3 text-yellow-600">
            <Database class="w-5 h-5" />
            <h2 class="text-xl font-semibold">Database Connection</h2>
        </div>
        <p class="text-gray-600 text-sm">
            Instructions on configuring PostgreSQL and SQLite for local and production environments using `.env`.
        </p>
    </div></a>
    <a href="/documentation/configuration">
    <div class="p-6 bg-white rounded-2xl shadow hover:shadow-md transition">
        <div class="flex items-center gap-3 mb-3 text-gray-800">
            <Settings class="w-5 h-5" />
            <h2 class="text-xl font-semibold">Configuration</h2>
        </div>
        <p class="text-gray-600 text-sm">
            Environment variables, ports, service dependencies, and Docker Compose setup across all microservices.
        </p>
    </div></a>
</div>
