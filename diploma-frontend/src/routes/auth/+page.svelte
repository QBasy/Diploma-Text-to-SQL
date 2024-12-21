<script lang="ts">
    import { isAuthenticated, loginSuccess } from '$lib/stores/authStore';
    import { authorization } from '$lib/api';
    import { slide } from 'svelte/transition';
    import {onMount} from "svelte";
    import {getAuthToken} from "../../lib/api";
    import {checkAuthStatus} from "../../lib/stores/authStore";

    let messageBox: string = '';
    let isRegister: boolean = false;

    let email: string = "";
    let password: string = '';
    let rememberMe: boolean = false;

    let name_register: string = '';
    let email_register: string = '';
    let password_register: string = '';
    let password_repeat: string = '';

    async function auth() {
        if (isRegister) {
            await register();
        } else {
            await login();
        }
    }

    async function login() {
        try {
            await authorization.login(email, password, rememberMe);
            location.reload()
        } catch (e) {
            messageBox = "User not Found";
        }
    }

    async function register() {
        if (password_repeat === password_register) {
            try {
                console.log(name_register, email_register, password_register)
                await authorization.register(name_register, email_register, password_register);
                messageBox = "User registered. Please login.";
                isRegister = false;
            } catch (e) {
                messageBox = "Error on creating user " + e;
            }
        } else {
            messageBox = "Passwords do not match";
        }
    }

    function toggleForm() {
        isRegister = !isRegister;
    }

    onMount(() => {
        if (checkAuthStatus()) {
            location.assign("/");
        }
    })
</script>

<div class="relative h-screen flex items-center justify-center overflow-hidden">
    <!-- Background Image -->
    <div class="absolute inset-0 bg-cover bg-center transition-all duration-500"
         style="background-image: url('/background.jpg');"
         style:transform={isRegister ? 'translateX(50%)' : 'translateX(-50%)'}
    ></div>
    <div class="absolute inset-0 bg-cover bg-center transition-all duration-500"
         style="background-image: url('/background.jpg');"
         style:transform={isRegister ? 'translateX(-50%)' : 'translateX(50%)'}
    ></div>

    <!-- Form Container -->
    <div class="relative w-full md:w-1/4 bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4 transition-transform duration-500 ease-in-out"
         style:transform={isRegister ? 'translateX(-90%)' : 'translateX(90%)'}>

        <div class="flex justify-between mb-4">
            <h2 class="text-xl font-bold">
                {#if isRegister}
                    Register
                {:else}
                    Sign In
                {/if}
            </h2>
            <button
                    on:click={toggleForm}
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                {#if isRegister}
                    Sign In
                {:else}
                    Register
                {/if}
            </button>
        </div>

        <!-- Form with sliding transition -->
        <form class="transition-transform duration-500 ease-in-out"
              in:slide={{ x: 200 }} out:slide={{ x: -200 }}>
            {#if isRegister}
                <!-- Registration form fields -->
                <div class="mb-4">
                    <label class="block text-gray-700 font-bold mb-2" for="name_register">Name</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="name_register" type="text" placeholder="Enter your name" bind:value={name_register} required/>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-700 font-bold mb-2" for="email_register">Email</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="email_register" type="email" placeholder="Enter your Email" bind:value={email_register} required/>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-700 font-bold mb-2" for="password_register">Password</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="password_register" type="password" placeholder="Enter your password" bind:value={password_register} required/>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-700 font-bold mb-2" for="password_repeat">Repeat Password</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="password_repeat" type="password" placeholder="Repeat your password" bind:value={password_repeat} required/>
                </div>
            {:else}
                <!-- Sign-in form fields -->
                <div class="mb-4">
                    <label class="block text-gray-700 font-bold mb-2" for="email">Email</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="email" type="email" placeholder="Enter your email" bind:value={email}/>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-700 font-bold mb-2" for="password">Password</label>
                    <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                           id="password" type="password" placeholder="Enter your password" bind:value={password}/>
                </div>
                <div class="mb-6">
                    <input id="rememberMe" type="checkbox" bind:checked={rememberMe} />
                    <label class="text-gray-700 font-bold mb-2" for="rememberMe">Remember Me</label>
                </div>
            {/if}
            <div class="flex items-center justify-between text-red-500">
                <p>{messageBox}</p>
            </div>

            <div class="flex items-center justify-between">
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        type="button"
                on:click={auth}>
                    {#if isRegister}
                        Register
                    {:else}
                        Sign In
                    {/if}
                </button>
            </div>
        </form>
    </div>
</div>
