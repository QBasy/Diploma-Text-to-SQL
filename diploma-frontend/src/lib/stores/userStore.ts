import { writable } from 'svelte/store';
import { login, register, logout, getCurrentUser } from '$lib/api/auth';

export interface User {
    id: string;
    name: string;
    email: string;
}

interface LoginRequest {
    email: string;
    password: string;
}

interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

export const loadingUser = writable(true);

export const userStore = writable<User | null>(null);


export const loginUser = async (data: LoginRequest) => {
    try {
        const response = await login(data);
        userStore.set(response.user);
        localStorage.setItem('token', response.token);
    } catch (error) {
        throw error;
    }
};

export const registerUser = async (data: RegisterRequest) => {
    try {
        const response = await register(data);

        userStore.set(response.user);
        localStorage.setItem('token', response.token);
    } catch (error) {
        throw error;
    }
};

export const logoutUser = async () => {
    try {
        localStorage.removeItem('token');
        userStore.set(null);
    } catch (error) {
        throw error;
    }
};

export const initializeUser = async () => {
    if (typeof localStorage === 'undefined') {
        loadingUser.set(false);
        return;
    }

    const token = localStorage.getItem('token');
    loadingUser.set(true);

    if (token) {
        try {
            const user = await getCurrentUser();
            userStore.set(user);
        } catch (error) {
            console.log(error);
            localStorage.removeItem('token');
            userStore.set(null);
        }
    } else {
        userStore.set(null);
    }

    loadingUser.set(false); // Завершаем загрузку
};

