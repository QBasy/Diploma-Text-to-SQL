<script lang="ts">
    import { page } from '$app/stores';
    import { logoutUser, userStore } from '$lib/stores/';
    import { goto } from '$app/navigation';
    import { derived } from 'svelte/store';
    import { UserCircle, LogOut, Home, Database, BarChart3, LineChart, Menu, X, Settings, History, ChevronDown } from 'lucide-svelte';
    import { onMount, onDestroy } from 'svelte';
    import { clickOutside } from '$lib/actions/clickOutside';

    const userLogIn = derived(userStore, $userStore => $userStore !== null);

    let showDropdown = false;
    let showMobileMenu = false;
    let lastScrollY = 0;
    let navHidden = false;
    let navElement: HTMLElement;

    function handleLogout(): void {
        logoutUser();
        showDropdown = false;
        showMobileMenu = false;
        setTimeout(() => goto('/'), 500);
    }

    function toggleDropdown(): void {
        showDropdown = !showDropdown;
    }

    function closeDropdown(): void {
        showDropdown = false;
    }

    function toggleMobileMenu(): void {
        showMobileMenu = !showMobileMenu;
        if (showMobileMenu) {
            document.body.classList.add('overflow-hidden');
        } else {
            document.body.classList.remove('overflow-hidden');
        }
    }

    function closeMobileMenu(): void {
        showMobileMenu = false;
        document.body.classList.remove('overflow-hidden');
    }

    // Handle scroll behavior for hiding/showing navbar
    function handleScroll() {
        const currentScrollY = window.scrollY;

        if (currentScrollY > 100) {
            if (currentScrollY > lastScrollY) {
                navHidden = true;
            } else {
                navHidden = false;
            }
        } else {
            navHidden = false;
        }

        lastScrollY = currentScrollY;
    }

    onMount(() => {
        if (typeof window !== 'undefined') {
            window.addEventListener('scroll', handleScroll);
        }
    });

    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('scroll', handleScroll);
        }
    });


    const navLinks = [
        { href: '/', label: 'Home', icon: Home },
        { href: '/generate/simple', label: 'Simple', icon: BarChart3 },
        { href: '/generate/complex', label: 'Complex', icon: LineChart }
    ];

    const profileLinks = [
        { href: '/profile', label: 'User Profile', icon: UserCircle },
        { href: '/profile/visualisation', label: 'Visualisation', icon: BarChart3 },
        { href: '/profile/database', label: 'Database', icon: Database },
        { href: '/profile/settings', label: 'Settings', icon: Settings },
        { href: '/profile/history', label: 'History', icon: History }
    ];

    $: isActive = (path) => {
        if (path === '/') {
            return $page.url.pathname === '/';
        }
        return $page.url.pathname.startsWith(path);
    };
</script>

<nav
        bind:this={navElement}
        class="fixed top-0 left-0 right-0 z-50 bg-white border-b border-gray-200 shadow-sm transition-transform duration-300 {navHidden ? '-translate-y-full' : 'translate-y-0'}"
        use:clickOutside={closeDropdown}
>
    <div class="container mx-auto px-4 py-3 flex justify-between items-center">
        <!-- Logo -->
        <a href="/" class="text-xl md:text-2xl font-bold text-blue-600 tracking-tight flex items-center gap-2">
            <LineChart class="w-6 h-6" />
            <span class="hidden sm:inline">SQL Translation</span>
            <span class="inline sm:hidden">Visualisation of SQL</span>
        </a>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-6 text-gray-700 font-medium">
            {#each navLinks as link}
                <a
                        href={link.href}
                        class="flex items-center gap-2 py-2 px-3 rounded-md transition-colors {isActive(link.href) ? 'bg-blue-50 text-blue-600 font-medium' : 'hover:text-blue-600 hover:bg-blue-50'}"
                >
                    <svelte:component this={link.icon} class="w-4 h-4" />
                    {link.label}
                </a>
            {/each}

            {#if $userLogIn}
                <div class="relative">
                    <button
                            on:click={toggleDropdown}
                            class="flex items-center gap-2 py-2 px-3 rounded-md transition-colors hover:text-blue-600 hover:bg-blue-50"
                    >
                        <UserCircle class="w-5 h-5" />
                        Profile
                        <ChevronDown class="w-4 h-4 {showDropdown ? 'rotate-180' : ''} transition-transform" />
                    </button>

                    {#if showDropdown}
                        <div class="absolute right-0 mt-2 w-64 bg-white rounded-lg shadow-xl border z-50 overflow-hidden">
                            <div class="pt-3 pb-2 px-4 border-b border-gray-100">
                                <p class="text-sm text-gray-500">Signed in as</p>
                                <p class="font-semibold text-gray-800">{$userStore?.email || 'User'}</p>
                            </div>

                            <div class="py-1">
                                {#each profileLinks as link}
                                    <a
                                            href={link.href}
                                            class="flex items-center gap-3 px-4 py-2 hover:bg-gray-50 text-gray-700 {$page.url.pathname === link.href ? 'text-blue-600 bg-blue-50' : ''}"
                                            on:click={closeDropdown}
                                    >
                                        <svelte:component this={link.icon} class="w-4 h-4" />
                                        {link.label}
                                    </a>
                                {/each}
                            </div>

                            <div class="border-t border-gray-100 mt-1 py-1">
                                <button
                                        on:click={handleLogout}
                                        class="flex items-center gap-3 px-4 py-2 text-red-600 hover:bg-red-50 w-full text-left"
                                >
                                    <LogOut class="w-4 h-4" />
                                    Logout
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>
            {:else}
                <a
                        href="/auth"
                        class="flex items-center gap-2 py-2 px-4 rounded-md bg-blue-100 text-blue-700 hover:bg-blue-200 transition-colors font-medium"
                >
                    <UserCircle class="w-4 h-4" />
                    Login
                </a>
            {/if}
        </div>

        <!-- Mobile Menu Toggle -->
        <div class="md:hidden flex items-center">
            <button
                    on:click={toggleMobileMenu}
                    class="p-2 text-gray-600 hover:text-blue-600 focus:outline-none"
                    aria-label={showMobileMenu ? "Close menu" : "Open menu"}
            >
                {#if showMobileMenu}
                    <X class="w-6 h-6" />
                {:else}
                    <Menu class="w-6 h-6" />
                {/if}
            </button>
        </div>
    </div>

    <!-- Mobile Navigation -->
    {#if showMobileMenu}
        <div class="fixed inset-0 z-40 bg-white md:hidden pt-16">
            <div class="p-4 space-y-3 overflow-y-auto h-full pb-24">
                <div class="space-y-1">
                    {#each navLinks as link}
                        <a
                                href={link.href}
                                class="flex items-center gap-3 p-3 rounded-lg {isActive(link.href) ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-50'}"
                                on:click={closeMobileMenu}
                        >
                            <svelte:component this={link.icon} class="w-5 h-5" />
                            {link.label}
                        </a>
                    {/each}
                </div>

                {#if $userLogIn}
                    <div class="mt-6 pt-6 border-t border-gray-100">
                        <div class="px-3 py-2">
                            <p class="text-sm text-gray-500">Signed in as</p>
                            <p class="font-semibold text-gray-800">{$userStore?.email || 'User'}</p>
                        </div>

                        <div class="mt-4 space-y-1">
                            {#each profileLinks as link}
                                <a
                                        href={link.href}
                                        class="flex items-center gap-3 p-3 rounded-lg {$page.url.pathname === link.href ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-50'}"
                                        on:click={closeMobileMenu}
                                >
                                    <svelte:component this={link.icon} class="w-5 h-5" />
                                    {link.label}
                                </a>
                            {/each}
                        </div>

                        <div class="mt-6">
                            <button
                                    on:click={handleLogout}
                                    class="flex items-center gap-3 p-3 rounded-lg text-red-600 hover:bg-red-50 w-full"
                            >
                                <LogOut class="w-5 h-5" />
                                Logout
                            </button>
                        </div>
                    </div>
                {:else}
                    <div class="mt-6 pt-6 border-t border-gray-100">
                        <a
                                href="/auth"
                                class="flex items-center justify-center gap-2 p-3 rounded-lg bg-blue-600 text-white hover:bg-blue-700"
                                on:click={closeMobileMenu}
                        >
                            <UserCircle class="w-5 h-5" />
                            Login / Register
                        </a>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</nav>

<div class="h-16"></div>
