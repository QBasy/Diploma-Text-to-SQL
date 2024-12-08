<script>
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

    function toggleAuth() {
        isAuth = !isAuth;
    }
</script>

<section class="container p-6 mx-auto w-full">
    <div class="mx-auto p-2 bg-gray-100 rounded-3xl flex flex-col">
        {#if isAuth}
            <div class="m-2 bg-gray-300 p-6 w-1/4 border-2 border-gray-400 rounded-3xl" transition:slide={{ duration: 1500 }}>
                <form class="mx-auto">
                    <div class="border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Username</strong>
                        <label for="username">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="text" required>
                        </label>
                    </div>
                    <div class="border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Password</strong>
                        <label for="password">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="text" required>
                        </label>
                    </div>
                </form>
            </div>
        {:else if !isAuth}
            <div class="m-2 bg-gray-300 p-6 w-1/4 border-2 border-gray-400 rounded-3xl items-center" transition:slide={{ duration: 1500 }}>
                <form class="mx-auto">
                    <div class="mx-auto border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Username</strong>
                        <label for="regUsername">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="text" required>
                        </label>
                    </div>
                    <div class="mx-auto border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Email</strong>
                        <label for="regEmail">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="email" required>
                        </label>
                    </div>
                    <div class="mx-auto border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Password</strong>
                        <label for="regPassword">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="password" required>
                        </label>
                    </div>
                    <div class="mx-auto border-2 border-gray-400 bg-gray-200 rounded-3xl m-2 p-2 w-full flex">
                        <strong>Confirm Password</strong>
                        <label for="confirmPassword">
                            <input class="bg-white border-2 border-gray-400 rounded-3xl mx-2 px-2 w-full flex" type="password" required>
                        </label>
                    </div>
                </form>
            </div>
        {/if}
        <div class="p-6 mx-auto">
            <button class="bg-blue-200 p-2 rounded-3xl" on:click={toggleAuth}>
                {#if isAuth}
                Register
                {:else}
                Login
                {/if}
            </button>
        </div>
    </div>
</section>