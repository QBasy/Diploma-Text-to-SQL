// src/lib/stores/authStore.js
import { writable } from 'svelte/store';
// @ts-ignore
import { getAuthToken } from '$lib/api';

export const isAuthenticated = writable(false);

export function checkAuthStatus() {
    const token = getAuthToken();
    isAuthenticated.set(!!token);
}
