// src/lib/stores/authStore.js
import { writable } from 'svelte/store';
// @ts-ignore
import { getAuthToken } from '$lib/api';
import {goto} from "$app/navigation";

export const isAuthenticated = writable(false);

export const loginSuccess = async (): Promise<void> => {
    await goto('/');
}
export function checkAuthStatus(): boolean {
    const token = getAuthToken();
    const isLoggedIn: boolean = !!token;
    isAuthenticated.set(isLoggedIn);
    return isLoggedIn;
}

