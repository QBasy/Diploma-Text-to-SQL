<script lang="ts">
    import { goto } from '$app/navigation';

    const links = [
        { name: 'Home', href: '/' },
        { name: 'Docs', href: '/documentation' },
        { name: 'Generate', href: '/generate' },
        { name: 'Contact', href: '/contact' }
    ];
</script>

<nav class="fixed top-0 w-full bg-white shadow z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
            <div class="flex items-center space-x-4">
                <a class="text-xl font-bold text-blue-600 hover:text-blue-800" href="/"></a>
                {#each links as link}
                    <a
                            href={link.href}
                            class="text-gray-600 hover:text-blue-600 transition duration-200"
                    >
                        {link.name}
                    </a>
                {/each}
            </div>
        </div>
    </div>
</nav>
