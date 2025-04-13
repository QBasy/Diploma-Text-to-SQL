<script lang="ts">
    import { goto } from '$app/navigation';
    import Notification from '$lib/components/Notification.svelte';
    import { forgotPassword } from '$lib/api';
    import { loginUser, registerUser, userStore } from '$lib/stores';
    import { onDestroy, onMount } from 'svelte';
    import { writable, type Writable } from 'svelte/store';

    let activeTab: 'login' | 'register' | 'forgot-password' = 'login';
    let notification: Writable<{ message: string; type: 'success' | 'error' } | null> = writable(null);

    let loginForm = {
        email: '',
        password: ''
    };

    let registerForm = {
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
    };

    let forgotPasswordForm = {
        email: ''
    };

    let currentBackground = 0;
    const backgrounds = ['/wallpaper1.jpg', '/wallpaper2.jpg'];

    let interval: ReturnType<typeof setInterval>;

    function startBackgroundChange() {
        interval = setInterval(() => {
            currentBackground = (currentBackground + 1) % backgrounds.length;
        }, 10000);
    }

    function showNotification(message: string, type: 'success' | 'error') {
        notification.set({ message, type });
        setTimeout(() => notification.set(null), 3000);
    }

    async function handleLogin() {
        if (!loginForm.email || !loginForm.password) {
            showNotification('Please fill in all fields', 'error');
            return;
        }

        try {
            await loginUser(loginForm);
            showNotification('Login successful', 'success');
            goto('/');
        } catch (error: any) {
            showNotification('Login failed: ' + (error.message || 'An error occurred'), 'error');
        }
    }

    async function handleRegister() {
        const { name, email, password, confirmPassword } = registerForm;

        if (!name || !email || !password || !confirmPassword) {
            showNotification('Please fill in all fields', 'error');
            return;
        }

        if (password !== confirmPassword) {
            showNotification('Passwords do not match', 'error');
            return;
        }

        try {
            await registerUser(registerForm);
            showNotification('Registration successful', 'success');
            goto('/');
        } catch (error: any) {
            showNotification('Registration failed: ' + error.message, 'error');
        }
    }

    async function handleForgotPassword() {
        if (!forgotPasswordForm.email) {
            showNotification('Please enter your email', 'error');
            return;
        }

        try {
            await forgotPassword(forgotPasswordForm);
            showNotification('Password reset link sent to your email', 'success');
        } catch (error: any) {
            showNotification('Failed to send reset link: ' + error.message, 'error');
        }
    }

    onMount(() => {
        const unsubscribe = userStore.subscribe((user) => {
            if (user) goto('/');
        });

        startBackgroundChange();
        return () => {
            unsubscribe();
        };
    });

    onDestroy(() => {
        clearInterval(interval);
    });
</script>

<div class="min-h-screen flex items-center justify-center relative">
    {#each backgrounds as bg, index}
        <div
                class="absolute inset-0 bg-cover bg-center transition-opacity duration-1000"
                style="background-image: url({bg}); opacity: {currentBackground === index ? 1 : 0};"
        />
    {/each}

    <div class="relative z-10 bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h1 class="text-2xl font-bold mb-6 text-center">
            {#if activeTab === 'login'}
                Login
            {:else if activeTab === 'register'}
                Register
            {:else}
                Forgot Password
            {/if}
        </h1>

        <div class="flex justify-center mb-6 flex-wrap gap-2">
            <button
                    class="px-4 py-2 rounded-md transition-all {activeTab === 'login' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => (activeTab = 'login')}
            >
                Login
            </button>
            <button
                    class="px-4 py-2 rounded-md transition-all {activeTab === 'register' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => (activeTab = 'register')}
            >
                Register
            </button>
            <button
                    class="px-4 py-2 rounded-md transition-all {activeTab === 'forgot-password' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => (activeTab = 'forgot-password')}
            >
                Forgot Password
            </button>
            <button
                    class="px-4 py-2 ml-2 hover:bg-blue-500 hover:text-white bg-gray-200 button-transition rounded-md"
                    on:click={() => goto('/')}
            >
                Back
            </button>
        </div>

        {#if activeTab === 'login'}
            <form on:submit|preventDefault={handleLogin} class="space-y-4">
                <input type="email" bind:value={loginForm.email} class="input" placeholder="Email" required />
                <input type="password" bind:value={loginForm.password} class="input" placeholder="Password" required />
                <button type="submit" class="btn-primary">Login</button>
            </form>
        {:else if activeTab === 'register'}
            <form on:submit|preventDefault={handleRegister} class="space-y-4">
                <input type="text" bind:value={registerForm.name} class="input" placeholder="Name" required />
                <input type="email" bind:value={registerForm.email} class="input" placeholder="Email" required />
                <input type="password" bind:value={registerForm.password} class="input" placeholder="Password" required />
                <input type="password" bind:value={registerForm.confirmPassword} class="input" placeholder="Confirm Password" required />
                <button type="submit" class="btn-success">Register</button>
            </form>
        {:else}
            <form on:submit|preventDefault={handleForgotPassword} class="space-y-4">
                <input type="email" bind:value={forgotPasswordForm.email} class="input" placeholder="Email" required />
                <button type="submit" class="btn-primary">Send Reset Link</button>
            </form>
        {/if}

        {#if $notification}
            <Notification message={$notification.message} type={$notification.type} />
        {/if}
    </div>
</div>

<style>
    .transition-opacity {
        transition: opacity 3s ease-in-out;
    }
    .input {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 0.375rem;
    }
    .btn-primary {
        width: 100%;
        background-color: #3b82f6;
        color: white;
        padding: 0.5rem;
        border-radius: 0.375rem;
        transition: background-color 0.3s;
    }
    .btn-primary:hover {
        background-color: #2563eb;
    }
    .btn-success {
        width: 100%;
        background-color: #22c55e;
        color: white;
        padding: 0.5rem;
        border-radius: 0.375rem;
        transition: background-color 0.3s;
    }
    .btn-success:hover {
        background-color: #16a34a;
    }
</style>
