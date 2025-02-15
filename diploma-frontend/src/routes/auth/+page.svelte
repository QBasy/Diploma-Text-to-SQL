<script lang="ts">
    import { goto } from '$app/navigation';
    import Notification from '$lib/components/Notification.svelte';
    import {forgotPassword} from "$lib/api";
    import {loginUser, registerUser, userStore} from "$lib/stores";
    import {onDestroy, onMount} from "svelte";

    let activeTab = 'login';
    let notification = '';

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

    let interval;

    function startBackgroundChange() {
        interval = setInterval(() => {
            currentBackground = (currentBackground + 1) % backgrounds.length;
        }, 10000);
    }

    async function handleLogin() {
        if (!loginForm.email || !loginForm.password) {
            notification = 'Please fill in all fields';
            return;
        }

        try {
            await loginUser(loginForm);
            notification = 'Login successful';

            goto('/');
        } catch (error: any) {
            notification = 'Login failed: ' + (error.message || 'An error occurred');
        }
    }

    async function handleRegister() {
        if (!registerForm.name || !registerForm.email || !registerForm.password || !registerForm.confirmPassword) {
            notification = 'Please fill in all fields';
            return;
        }

        if (registerForm.password !== registerForm.confirmPassword) {
            notification = 'Passwords do not match';
            return;
        }

        try {
            await registerUser(registerForm);
            notification = 'Registration successful';

            goto('/');
        } catch (error: any) {
            notification = 'Registration failed: ' + error.message;
        }
    }

    async function handleForgotPassword() {
        if (!forgotPasswordForm.email) {
            notification = 'Please enter your email';
            return;
        }

        try {
            await forgotPassword(forgotPasswordForm);
            notification = 'Password reset link sent to your email';
        } catch (error: any) {
            notification = 'Failed to send reset link: ' + error.message;
        }
    }

    onMount(() => {
        userStore.subscribe((user) => { if (user) goto('/') });
        startBackgroundChange();
    });

    onDestroy(() => {
        clearInterval(startBackgroundChange);
    });
</script>

<div class="min-h-screen flex items-center justify-center relative">
    {#each backgrounds as bg, index}
        <div
                class="absolute inset-0 bg-cover bg-center transition-opacity duration-1000"
                style="background-image: url({bg}); opacity: {currentBackground === index ? 1 : 0};">
        </div>
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

        <div class="flex justify-center mb-6">
            <button
                    class="px-4 py-2 mr-2 {activeTab === 'login' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => activeTab = 'login'}
            >
                Login
            </button>
            <button
                    class="px-4 py-2 mr-2 {activeTab === 'register' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => activeTab = 'register'}
            >
                Register
            </button>
            <button
                    class="px-4 py-2 {activeTab === 'forgot-password' ? 'bg-blue-500 text-white' : 'bg-gray-200'}"
                    on:click={() => activeTab = 'forgot-password'}
            >
                Forgot Password
            </button>

            <button
                    class="px-4 py-2 ml-2 hover:bg-blue-500 hover:text-white bg-gray-200 button-transition"
                    on:click={() => goto("/")}
            >
                Back
            </button>
        </div>

        {#if activeTab === 'login'}
            <form on:submit|preventDefault={handleLogin} class="space-y-4">
                <div>
                    <label class="block text-gray-700">Email</label>
                    <input
                            type="email"
                            bind:value={loginForm.email}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <div>
                    <label class="block text-gray-700">Password</label>
                    <input
                            type="password"
                            bind:value={loginForm.password}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <button
                        type="submit"
                        class="w-full bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600"
                >
                    Login
                </button>

            </form>
        {:else if activeTab === 'register'}
            <form on:submit|preventDefault={handleRegister} class="space-y-4">
                <div>
                    <label class="block text-gray-700">Name</label>
                    <input
                            type="text"
                            bind:value={registerForm.name}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <div>
                    <label class="block text-gray-700">Email</label>
                    <input
                            type="email"
                            bind:value={registerForm.email}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <div>
                    <label class="block text-gray-700">Password</label>
                    <input
                            type="password"
                            bind:value={registerForm.password}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <div>
                    <label class="block text-gray-700">Confirm Password</label>
                    <input
                            type="password"
                            bind:value={registerForm.confirmPassword}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <button
                        type="submit"
                        class="w-full bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600"
                >
                    Register
                </button>
            </form>
        {:else}
            <form on:submit|preventDefault={handleForgotPassword} class="space-y-4">
                <div>
                    <label class="block text-gray-700">Email</label>
                    <input
                            type="email"
                            bind:value={forgotPasswordForm.email}
                            class="w-full p-2 border rounded-md"
                            required
                    />
                </div>
                <button
                        type="submit"
                        class="w-full bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600"
                >
                    Send Reset Link
                </button>
            </form>
        {/if}

        {#if notification}
            <Notification message={notification} type={notification.includes('success') ? 'success' : 'error'} />
        {/if}

    </div>

</div>

<style>
    .transition-opacity {
        transition: opacity 3s ease-in-out;
    }

    .bg-cover {
        background-size: cover;
        background-position: center;
    }

    .button-transition {
        transition: background-color 0.3s ease, color 0.3s ease;
    }
</style>
