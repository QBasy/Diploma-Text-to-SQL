<!-- Spinner.svelte -->
<script lang="ts">
    // Props with defaults
    export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
    export let color: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'neutral' = 'primary';
    export let variant: 'border' | 'dots' = 'border';
    export let className: string = '';

    const sizeMap = {
        sm: {
            border: 'w-4 h-4 border-2',
            dots: 'gap-1'
        },
        md: {
            border: 'w-6 h-6 border-2',
            dots: 'gap-1.5'
        },
        lg: {
            border: 'w-8 h-8 border-3',
            dots: 'gap-2'
        },
        xl: {
            border: 'w-12 h-12 border-4',
            dots: 'gap-2.5'
        }
    };

    const colorMap = {
        primary: 'border-blue-600 border-t-blue-200',
        secondary: 'border-purple-600 border-t-purple-200',
        success: 'border-green-600 border-t-green-200',
        warning: 'border-yellow-600 border-t-yellow-200',
        danger: 'border-red-600 border-t-red-200',
        neutral: 'border-gray-600 border-t-gray-200'
    };

    const dotColorMap = {
        primary: 'bg-blue-600',
        secondary: 'bg-purple-600',
        success: 'bg-green-600',
        warning: 'bg-yellow-600',
        danger: 'bg-red-600',
        neutral: 'bg-gray-600'
    };

    const dotSizeMap = {
        sm: 'w-1.5 h-1.5',
        md: 'w-2 h-2',
        lg: 'w-2.5 h-2.5',
        xl: 'w-3 h-3'
    };
</script>

{#if variant === 'border'}
    <div class="inline-block {sizeMap[size].border} rounded-full {colorMap[color]} animate-spin {className}" role="status" aria-label="Loading">
        <span class="sr-only">Loading...</span>
    </div>
{:else if variant === 'dots'}
    <div class="flex items-center justify-center {sizeMap[size].dots} {className}" role="status" aria-label="Loading">
        <div class="{dotSizeMap[size]} rounded-full {dotColorMap[color]} animate-bounce" style="animation-delay: 0ms;"></div>
        <div class="{dotSizeMap[size]} rounded-full {dotColorMap[color]} animate-bounce" style="animation-delay: 150ms;"></div>
        <div class="{dotSizeMap[size]} rounded-full {dotColorMap[color]} animate-bounce" style="animation-delay: 300ms;"></div>
        <span class="sr-only">Loading...</span>
    </div>
{/if}