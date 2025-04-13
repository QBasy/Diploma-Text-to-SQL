<script lang="ts">
    export let message: string = '';
    export let type: 'success' | 'error' = 'success';
    export let duration: number = 3000;

    let visible = false;
    let timeout: ReturnType<typeof setTimeout>;

    $: if (message) {
        visible = true;
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            visible = false;
        }, duration);
    }
</script>

{#if visible}
    <div
            class="fixed top-4 right-4 z-50 px-4 py-3 rounded-md shadow-lg text-white transition-all duration-300"
            class:bg-green-500={type === 'success'}
            class:bg-red-500={type === 'error'}
    >
        {message}
    </div>
{/if}
