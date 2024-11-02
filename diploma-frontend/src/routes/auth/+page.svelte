<script>
    import {slide} from "svelte/transition";

    let isAuth = true;

    let username = "";
    let password = "";

    let regUsername = "";
    let regEmail = "";
    let regPassword = "";
    let confirmPassword = "";

    function toggleAuth() {
        isAuth = !isAuth;
    }
</script>

<section class="container p-6 mx-auto w-full">
    <div class="mx-auto p-2 bg-gray-100 rounded-3xl flex flex-col">
        {#if isAuth}
            <div class="bg-emerald-800">
                <form class="bg-amber-300">
                    <label for="username">
                        <input class="bg-white" type="text" required>
                    </label>
                    <label for="password">
                        <input class="bg-white" type="text" required>
                    </label>
                </form>
            </div>
        {:else if !isAuth}
            <div class="m-2 bg-blue-100 p-6 border-2 border-gray-400 rounded-3xl">
                <form class="">
                    <div>
                        <strong>Username</strong>
                        <label for="regUsername">
                            <input class="bg-white" type="text" required>
                        </label>
                    </div>
                    <div>
                        <strong>Email</strong>
                        <label for="regEmail">
                            <input class="bg-white" type="email" required>
                        </label>
                    </div>
                    <div>
                        <strong>Password</strong>
                        <label for="regPassword">
                            <input class="bg-white" type="password" required>
                        </label>
                    </div>
                    <div>
                        <strong>Confirm Password</strong>
                        <label for="confirmPassword">
                            <input class="bg-white" type="password" required>
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