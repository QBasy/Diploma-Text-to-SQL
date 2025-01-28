<script>
    export let message = '';
    export let type = 'success';
    export let duration = 3000;

    let visible = false;

    $: if (message) {
        visible = true;
        setTimeout(() => {
            visible = false;
        }, duration);
    }
</script>

{#if visible}
    <div
            class="notification fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 transition-all"
            class:bg-green-500={type === 'success'}
            class:bg-red-500={type === 'error'}
            class:text-white={true}
    >
        {message}
    </div>
{/if}

<style>
    .notification {
        transition: opacity 0.3s, transform 0.3s;
    }
    .notification-enter {
        opacity: 0;
        transform: translateY(-20px);
    }
    .notification-exit {
        opacity: 0;
        transform: translateY(-20px);
    }
</style>