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

export const userStore = writable<User | null>(null);


export const loginUser = async (data: LoginRequest) => {
    try {
        const response = await login(data);
        userStore.set(response.user); // Сохраняем данные пользователя
        localStorage.setItem('token', response.token); // Сохраняем токен
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
    const token = localStorage.getItem('token');
    if (token) {
        try {
            const user = await getCurrentUser();
            userStore.set(user);
        } catch (error) {
            console.log(error)
            localStorage.removeItem('token');
            userStore.set(null);
        }
    }
};