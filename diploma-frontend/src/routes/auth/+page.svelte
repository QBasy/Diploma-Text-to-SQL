<script>
    let isRegister = false;

    import {slide} from "svelte/transition";

    let isAuth = true;

    let username = "";
    let password = "";

    let regUsername = "";
    let regEmail = "";
    let regPassword = "";
    let confirmPassword = "";

    function validateEmail() {
        /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(regEmail)
    }

    function toggleForm() {
        isRegister = !isRegister;
    }
</script>

<div class="flex items-center justify-center h-screen">
    <div class="w-full max-w-md">
        <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
            <div class="flex justify-between mb-4">
                <h2 class="text-xl font-bold">
                    {#if isRegister}
                        Register
                    {:else}
                        Sign Up
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

            <form
                    class="transition-all duration-500 ease-in-out"
                    style:transform={isRegister ? 'translateX(100%)' : 'translateX(0)'}>
                {#if isRegister}
                    <!-- Registration form fields -->
                    <div class="mb-4">
                        <label class="block text-gray-700 font-bold mb-2" for="name">
                            Name
                        </label>
                        <input
                                class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                id="name"
                                type="text"
                                placeholder="Enter your name" />
                    </div>
                {:else}
                    <!-- Sign-up form fields -->
                    <div class="mb-4">
                        <label class="block text-gray-700 font-bold mb-2" for="email">
                            Email
                        </label>
                        <input
                                class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                id="email"
                                type="email"
                                placeholder="Enter your email" />
                    </div>
                {/if}

                <div class="mb-6">
                    <label class="block text-gray-700 font-bold mb-2" for="password">
                        Password
                    </label>
                    <input
                            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                            id="password"
                            type="password"
                            placeholder="Enter your password" />
                </div>

                <div class="flex items-center justify-between">
                    <button
                            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                            type="button">
                        {#if isRegister}
                            Register
                        {:else}
                            Sign Up
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>
