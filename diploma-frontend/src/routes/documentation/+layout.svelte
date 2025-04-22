<script lang="ts">
    import { page } from '$app/stores';
    import { derived } from 'svelte/store';
    import NavbarDocumentation from '$lib/components/NavbarDocumentation.svelte';
    import Footer from '$lib/components/Footer.svelte';
    import { BookOpen, ChartBar, Database, Settings, Play } from 'lucide-svelte';

    // Mobile menu state
    import { writable } from 'svelte/store';
    const isMobileMenuOpen = writable(false);

    const toggleMobileMenu = () => {
        isMobileMenuOpen.update(value => !value);
    };

    const services = [
        { name: 'How to Start', path: '/documentation/start', icon: Play },
        { name: 'Chart Types', path: '/documentation/charts', icon: ChartBar },
        { name: 'Tech Stack', path: '/documentation/stack', icon: BookOpen },
        { name: 'Configuration', path: '/documentation/configuration', icon: Settings },
        { name: 'Database Connection', path: '/documentation/database', icon: Database }
    ];

    const currentPath = derived(page, ($page) => $page.url.pathname);
</script>

<NavbarDocumentation />

<div class="flex min-h-screen pt-16 bg-gray-50">
    <!-- Mobile sidebar toggle -->
    <div class="md:hidden fixed bottom-4 right-4 z-20">
        <button
                on:click={toggleMobileMenu}
                class="flex items-center justify-center w-12 h-12 rounded-full bg-blue-600 text-white shadow-lg hover:bg-blue-700 transition-colors"
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
        </button>
    </div>

    <!-- Mobile sidebar -->
    {#if $isMobileMenuOpen}
        <div class="md:hidden fixed inset-0 z-10 bg-black bg-opacity-50" on:click={toggleMobileMenu}></div>
        <aside class="md:hidden fixed bottom-0 inset-x-0 z-10 bg-white rounded-t-xl shadow-lg transform transition-transform duration-300" style="max-height: 80vh; overflow-y: auto;">
            <div class="p-4 border-b border-gray-100">
                <h2 class="text-lg font-semibold flex items-center justify-between">
                    <a href="/documentation"><span>Documentation</span></a>
                    <button on:click={toggleMobileMenu} class="text-gray-500 hover:text-gray-700">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </h2>
            </div>
            <ul class="p-4 space-y-1">
                {#each services as service}
                    {@const isActive = $currentPath === service.path}
                    <li>
                        <a
                                href={service.path}
                                class="flex items-center gap-3 p-3 rounded-lg transition-colors {isActive ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-100'}"
                        >
                            <svelte:component this={service.icon} class="w-5 h-5 {isActive ? 'text-blue-600' : 'text-gray-500'}" />
                            {service.name}
                        </a>
                    </li>
                {/each}
            </ul>
        </aside>
    {/if}

    <!-- Desktop sidebar -->
    <aside class="w-64 py-6 px-4 bg-white shadow-md border-r hidden md:block sticky top-16 h-[calc(100vh-4rem)]">
        <a href="/documentation"><h2 class="text-lg font-semibold mb-6 text-gray-800 px-2">Documentation</h2></a>
        <ul class="space-y-1">
            {#each services as service}
                {@const isActive = $currentPath === service.path}
                <li>
                    <a
                            href={service.path}
                            class="flex items-center gap-3 p-3 rounded-lg transition-colors {isActive ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-100'}"
                    >
                        <svelte:component this={service.icon} class="w-5 h-5 {isActive ? 'text-blue-600' : 'text-gray-500'}" />
                        {service.name}
                    </a>
                </li>
            {/each}
        </ul>
    </aside>

    <main class="flex-1 p-6 md:p-8 pb-20">
        <div class="max-w-4xl mx-auto">
            <slot />
        </div>
    </main>
</div>

<Footer />