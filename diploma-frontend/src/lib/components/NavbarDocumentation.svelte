<script lang="ts">
    import { page } from '$app/stores';
    import { BookOpen, Home, MessageSquare, BarChart3 } from 'lucide-svelte';

    const links = [
        { name: 'Home', href: '/', icon: Home },
        { name: 'Docs', href: '/documentation', icon: BookOpen },
        { name: 'Generate', href: '/generate', icon: BarChart3 },
        { name: 'Contact', href: '/contact', icon: MessageSquare }
    ];

    // Mobile dropdown menu state
    import { writable } from 'svelte/store';
    const isDropdownOpen = writable(false);

    const toggleDropdown = () => {
        isDropdownOpen.update(value => !value);
    };
</script>

<nav class="fixed top-0 left-0 right-0 z-50 w-full bg-white shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
            <div class="flex items-center">
                <a href="/" class="flex items-center gap-2 text-2xl font-bold text-blue-600 hover:text-blue-800 transition-colors">
                    <BarChart3 class="w-6 h-6" />
                    Visualization of NLQ
                </a>

                <!-- Desktop navigation links -->
                <div class="hidden md:flex ml-8 space-x-6">
                    {#each links as link}
                        {@const isActive = $page.url.pathname === link.href ||
                        ($page.url.pathname.startsWith(link.href + '/') && link.href !== '/')}
                        <a
                                href={link.href}
                                class="flex items-center gap-2 py-2 px-3 rounded-md transition-colors {
                                isActive
                                ? 'bg-blue-50 text-blue-600 font-medium'
                                : 'text-gray-600 hover:text-blue-600 hover:bg-blue-50'
                            }"
                        >
                            <svelte:component this={link.icon} class="w-4 h-4" />
                            {link.name}
                        </a>
                    {/each}
                </div>
            </div>

            <!-- Mobile menu button -->
            <div class="md:hidden">
                <button
                        on:click={toggleDropdown}
                        class="p-2 rounded-md text-gray-600 hover:text-blue-600 hover:bg-blue-50 focus:outline-none"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                </button>
            </div>
        </div>
    </div>

    <!-- Mobile dropdown menu -->
    {#if $isDropdownOpen}
        <div class="md:hidden">
            <div class="px-3 py-3 space-y-1 bg-white border-t border-gray-100 shadow-lg">
                {#each links as link}
                    {@const isActive = $page.url.pathname === link.href ||
                    ($page.url.pathname.startsWith(link.href + '/') && link.href !== '/')}
                    <a
                            href={link.href}
                            class="flex items-center gap-3 px-4 py-3 rounded-md {
                            isActive
                            ? 'bg-blue-50 text-blue-600 font-medium'
                            : 'text-gray-600 hover:text-blue-600 hover:bg-blue-50'
                        }"
                            on:click={toggleDropdown}
                    >
                        <svelte:component this={link.icon} class="w-5 h-5" />
                        {link.name}
                    </a>
                {/each}
            </div>
        </div>
    {/if}
</nav>