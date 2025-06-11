<script lang="ts">
    import { onMount } from 'svelte';
    import { userStore } from '$lib/stores/userStore';
    import Notification from '$lib/components/Notification.svelte';
    import { goto } from '$app/navigation';
    import type { User } from '$lib/stores/userStore';
    import { UserCircle, Lock, Mail, Shield, AlertTriangle, Save, CheckCircle } from 'lucide-svelte';
    import { getCurrentUser, logout } from '$lib/api/auth';
    import api from '$lib/api/index';

    let currentUser: User | null;
    let notification = '';
    let notificationType = '';
    let isLoading = false;
    let userProfile = {
        name: '',
        email: ''
    };

    let passwordForm = {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
    };

    userStore.subscribe(value => {
        currentUser = value;
        if (currentUser) {
            userProfile.name = currentUser.name || '';
            userProfile.email = currentUser.email || '';
        }
    });

    onMount(async () => {
        isLoading = true;
        try {
            // Use the API service to fetch user data
            const userData = await getCurrentUser();
            userProfile.name = userData.name;
            userProfile.email = userData.email;
        } catch (error) {
            showNotification('Failed to load profile data', 'error');
        } finally {
            isLoading = false;
        }
    });

    function showNotification(message: string, type: 'success' | 'error') {
        notification = message;
        notificationType = type;
        setTimeout(() => {
            notification = '';
        }, 5000);
    }

    async function updatePassword() {
        if (passwordForm.newPassword !== passwordForm.confirmPassword) {
            showNotification('Passwords do not match', 'error');
            return;
        }

        if (passwordForm.newPassword.length < 8) {
            showNotification('Password must be at least 8 characters', 'error');
            return;
        }

        isLoading = true;

        try {
            // Using direct API call since you don't have a dedicated function for password change
            await api.post('/auth/change-password', {
                old_password: passwordForm.oldPassword,
                new_password: passwordForm.newPassword
            });

            showNotification('Password updated successfully', 'success');
            passwordForm = { oldPassword: '', newPassword: '', confirmPassword: '' };
        } catch (error) {
            showNotification(error.response?.data?.error || 'Failed to update password', 'error');
        } finally {
            isLoading = false;
        }
    }

    async function deleteAccount() {
        if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
            isLoading = true;

            try {
                // Direct API call for account deletion
                await api.delete('/auth/delete-account');

                showNotification('Account successfully deleted', 'success');
                await logout();
                userStore.set(null);
                setTimeout(() => goto('/'), 2000);
            } catch (error) {
                showNotification(error.response?.data?.error || 'Failed to delete account', 'error');
            } finally {
                isLoading = false;
            }
        }
    }
</script>

<svelte:head>
    <title>User Profile</title>
</svelte:head>

<div class="container mx-auto p-4 md:p-6 max-w-3xl">
    <div class="flex items-center gap-3 mb-8">
        <UserCircle class="w-8 h-8 text-blue-600" />
        <h1 class="text-2xl font-bold text-gray-800">User Profile</h1>
    </div>

    {#if isLoading}
        <div class="flex justify-center my-12">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
    {:else}
        <div class="bg-white shadow-md rounded-lg overflow-hidden mb-6">
            <div class="border-b border-gray-100 bg-blue-50 p-4">
                <div class="flex items-center gap-2">
                    <UserCircle class="w-5 h-5 text-blue-600" />
                    <h2 class="text-lg font-semibold text-gray-800">Personal Information</h2>
                </div>
                <p class="text-sm text-gray-500 mt-1">Your account details and information</p>
            </div>

            <div class="p-6">
                <div class="grid md:grid-cols-2 gap-6">
                    <div>
                        <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <UserCircle class="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                    id="username"
                                    type="text"
                                    value={userProfile.name}
                                    disabled
                                    class="pl-10 w-full p-2 border rounded-md bg-gray-50 text-gray-700"
                            />
                        </div>
                    </div>
                    <div>
                        <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Mail class="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                    id="email"
                                    type="email"
                                    value={userProfile.email}
                                    disabled
                                    class="pl-10 w-full p-2 border rounded-md bg-gray-50 text-gray-700"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="bg-white shadow-md rounded-lg overflow-hidden mb-6">
            <div class="border-b border-gray-100 bg-blue-50 p-4">
                <div class="flex items-center gap-2">
                    <Lock class="w-5 h-5 text-blue-600" />
                    <h2 class="text-lg font-semibold text-gray-800">Security</h2>
                </div>
                <p class="text-sm text-gray-500 mt-1">Update your password and account security</p>
            </div>

            <div class="p-6">
                <form on:submit|preventDefault={updatePassword} class="space-y-4">
                    <div>
                        <label for="oldPassword" class="block text-sm font-medium text-gray-700 mb-1">Current Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Shield class="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                    id="oldPassword"
                                    type="password"
                                    bind:value={passwordForm.oldPassword}
                                    class="pl-10 w-full p-2 border rounded-md"
                                    placeholder="Enter your current password"
                                    required
                            />
                        </div>
                    </div>

                    <div class="grid md:grid-cols-2 gap-4">
                        <div>
                            <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
                            <div class="relative">
                                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <Lock class="h-5 w-5 text-gray-400" />
                                </div>
                                <input
                                        id="newPassword"
                                        type="password"
                                        bind:value={passwordForm.newPassword}
                                        class="pl-10 w-full p-2 border rounded-md"
                                        placeholder="Enter new password"
                                        required
                                />
                            </div>
                        </div>
                        <div>
                            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm New Password</label>
                            <div class="relative">
                                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <CheckCircle class="h-5 w-5 text-gray-400" />
                                </div>
                                <input
                                        id="confirmPassword"
                                        type="password"
                                        bind:value={passwordForm.confirmPassword}
                                        class="pl-10 w-full p-2 border rounded-md"
                                        placeholder="Confirm new password"
                                        required
                                />
                            </div>
                        </div>
                    </div>

                    <div class="pt-2">
                        <button
                                type="submit"
                                class="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
                                disabled={isLoading}
                        >
                            <Save class="w-4 h-4" />
                            Update Password
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <div class="bg-white shadow-md rounded-lg overflow-hidden">
            <div class="border-b border-gray-100 bg-red-50 p-4">
                <div class="flex items-center gap-2">
                    <AlertTriangle class="w-5 h-5 text-red-600" />
                    <h2 class="text-lg font-semibold text-red-600">Danger Zone</h2>
                </div>
                <p class="text-sm text-gray-600 mt-1">Actions in this section cannot be undone</p>
            </div>

            <div class="p-6">
                <div class="flex items-center justify-between p-4 border border-red-100 rounded-lg bg-red-50">
                    <div>
                        <h3 class="font-medium text-gray-800">Delete Account</h3>
                        <p class="text-sm text-gray-600">Permanently delete your account and all data</p>
                    </div>
                    <button
                            on:click={deleteAccount}
                            class="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 transition-colors flex items-center gap-2"
                            disabled={isLoading}
                    >
                        <AlertTriangle class="w-4 h-4" />
                        Delete Account
                    </button>
                </div>
            </div>
        </div>
    {/if}

    {#if notification}
        <Notification
                message={notification}
                type={notificationType}
        />
    {/if}
</div>