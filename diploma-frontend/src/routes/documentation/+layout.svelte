<script lang="ts">
    import { page } from '$app/stores';
    import { derived } from 'svelte/store';
    import NavbarDocumentation from '$lib/components/NavbarDocumentation.svelte';
    import Footer from '$lib/components/Footer.svelte';

    const services = [
        { name: 'How to Start', path: '/documentation/start' },
        { name: 'Chart Types', path: '/documentation/charts' },
        { name: 'Tech Stack', path: '/documentation/stack' },
        { name: 'Configuration', path: '/documentation/configuration' },
        { name: 'Database Connection', path: '/documentation/database' },
        { name: 'Architecture', path: '/documentation/architecture' }
    ];

    const currentPath = derived(page, ($page) => $page.url.pathname);
</script>

<NavbarDocumentation />

<div class="flex min-h-screen pt-20 bg-gray-50">
    <aside class="w-64 p-6 bg-white shadow-md border-r hidden md:block">
        <h2 class="text-lg font-semibold mb-4">Sections</h2>
        <ul class="space-y-2">
            {#each services as service}
                {#await currentPath}
                    <li><a href={service.path}>{service.name}</a></li>
                {:then $currentPath}
                    <li>
                        <a
                                href={service.path}
                                class="transition block rounded px-2 py-1.5
								{ $currentPath === service.path
									? 'text-white bg-blue-600 font-semibold'
									: 'text-gray-700 hover:text-blue-600' }"
                        >
                            {service.name}
                        </a>
                    </li>
                {/await}
            {/each}
        </ul>
    </aside>

    <main class="flex-1 p-6">
        <slot />
    </main>
</div>

<Footer />
