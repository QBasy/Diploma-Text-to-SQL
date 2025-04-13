<script lang="ts">
    import NavbarDocumentation from '$lib/components/NavbarDocumentation.svelte';
    import Footer from '$lib/components/Footer.svelte';

    const services = [
        { name: 'Overview', path: '/documentation' },
        { name: 'Auth Service', path: '/documentation/auth-service' },
        { name: 'API Gateway', path: '/documentation/api-gateway' },
        { name: 'Text-to-SQL Service', path: '/documentation/text-to-sql-service' },
        { name: 'Visualisation Service', path: '/documentation/visualisation-service' },
        { name: 'Database Service', path: '/documentation/database-service' },
        { name: 'History Service', path: '/documentation/history-service' }
    ];
</script>

<NavbarDocumentation />

<div class="flex min-h-screen pt-20 bg-gray-50">
    <aside class="w-64 p-6 bg-white shadow-md border-r hidden md:block">
        <h2 class="text-lg font-semibold mb-4">Services</h2>
        <ul class="space-y-2">
            {#each services as service}
                <li>
                    <a href={service.path} class="text-gray-700 hover:text-blue-600 transition">{service.name}</a>
                </li>
            {/each}
        </ul>
    </aside>

    <main class="flex-1 p-6">
        <slot />
    </main>
</div>

<Footer />
