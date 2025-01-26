<script>
    import { onMount } from 'svelte';
    import Notification from '$lib/components/Notification.svelte';

    let theme = 'light';
    let language = 'en';
    let notification = '';

    onMount(() => {
        // TODO: Load actual user settings
        theme = localStorage.getItem('theme') || 'light';
        language = localStorage.getItem('language') || 'en';
    });

    function updateTheme(newTheme) {
        theme = newTheme;
        localStorage.setItem('theme', newTheme);
        document.documentElement.classList.toggle('dark', newTheme === 'dark');
        notification = `Theme set to ${newTheme}`;
    }

    function updateLanguage(newLanguage) {
        language = newLanguage;
        localStorage.setItem('language', newLanguage);
        notification = `Language set to ${newLanguage}`;
    }
</script>

<div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">User Settings</h1>

    <div class="bg-white shadow-md rounded-lg p-6 mb-6">
        <h2 class="text-xl font-semibold mb-4">Theme</h2>
        <div class="flex space-x-4">
            <button
                    class="px-4 py-2 rounded-md {theme === 'light' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => updateTheme('light')}
            >
                Light
            </button>
            <button
                    class="px-4 py-2 rounded-md {theme === 'dark' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => updateTheme('dark')}
            >
                Dark
            </button>
        </div>
    </div>

    <div class="bg-white shadow-md rounded-lg p-6">
        <h2 class="text-xl font-semibold mb-4">Language</h2>
        <div class="flex space-x-4">
            <button
                    class="px-4 py-2 rounded-md {language === 'en' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => updateLanguage('en')}
            >
                English
            </button>
            <button
                    class="px-4 py-2 rounded-md {language === 'ru' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => updateLanguage('ru')}
            >
                Русский
            </button>
        </div>
    </div>

    {#if notification}
        <Notification
                message={notification}
                type="success"
        />
    {/if}
</div>