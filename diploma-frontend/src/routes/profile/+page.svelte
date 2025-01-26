<script>
    import { userStore } from '$lib/stores/userStore';
    import Notification from '$lib/components/Notification.svelte';

    let currentUser;
    let notification = '';

    let passwordForm = {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
    };

    userStore.subscribe(value => {
        currentUser = value.user;
    });

    function updatePassword() {
        if (passwordForm.newPassword !== passwordForm.confirmPassword) {
            notification = 'Passwords do not match';
            return;
        }

        // TODO: Implement actual password update logic
        notification = 'Password updated successfully';
        passwordForm = { oldPassword: '', newPassword: '', confirmPassword: '' };
    }

    function deleteAccount() {
        if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
            // TODO: Implement account deletion logic
            notification = 'Account deleted';
        }
    }
</script>

<div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">User Profile</h1>

    <div class="bg-white shadow-md rounded-lg p-6 mb-6">
        <h2 class="text-xl font-semibold mb-4">User Information</h2>

        <div class="grid md:grid-cols-2 gap-4">
            <div>
                <label class="block text-gray-700">Name</label>
                <input
                        type="text"
                        value={currentUser?.name || ''}
                        disabled
                        class="w-full p-2 border rounded-md bg-gray-100"
                >
            </div>
            <div>
                <label class="block text-gray-700">Email</label>
                <input
                        type="email"
                        value={currentUser?.email || ''}
                        disabled
                        class="w-full p-2 border rounded-md bg-gray-100"
                >
            </div>
        </div>
    </div>

    <div class="bg-white shadow-md rounded-lg p-6 mb-6">
        <h2 class="text-xl font-semibold mb-4">Change Password</h2>

        <form on:submit|preventDefault={updatePassword} class="space-y-4">
            <div>
                <label class="block text-gray-700">Old Password</label>
                <input
                        type="password"
                        bind:value={passwordForm.oldPassword}
                        class="w-full p-2 border rounded-md"
                        required
                >
            </div>
            <div>
                <label class="block text-gray-700">New Password</label>
                <input
                        type="password"
                        bind:value={passwordForm.newPassword}
                        class="w-full p-2 border rounded-md"
                        required
                >
            </div>
            <div>
                <label class="block text-gray-700">Confirm New Password</label>
                <input
                        type="password"
                        bind:value={passwordForm.confirmPassword}
                        class="w-full p-2 border rounded-md"
                        required
                >
            </div>
            <button
                    type="submit"
                    class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600"
            >
                Change Password
            </button>
        </form>
    </div>

    <div class="bg-white shadow-md rounded-lg p-6">
        <h2 class="text-xl font-semibold mb-4 text-red-600">Danger Zone</h2>
        <button
                on:click={deleteAccount}
                class="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600"
        >
            Delete Account
        </button>
    </div>

    {#if notification}
        <Notification
                message={notification}
                type={notification.includes('successfully') ? 'success' : 'error'}
        />
    {/if}
</div>
