<script lang="ts">
    import { page } from '$app/stores';

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
            <div class="flex items-center space-x-6">
                <a href="/" class="text-2xl font-bold text-blue-600 hover:text-blue-800">VizSQL</a>
                {#each links as link}
                    <a
                            href={link.href}
                            class="pb-1 border-b-2 transition duration-200
							{($page.url.pathname === link.href
								|| ($page.url.pathname.startsWith(link.href + '/') && link.href !== '/'))
								? 'border-blue-600 font-semibold text-blue-600'
								: 'border-transparent text-gray-600 hover:text-blue-600'}"
                    >
                        {link.name}
                    </a>
                {/each}
            </div>
        </div>
    </div>
</nav>
