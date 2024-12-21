import axios from 'axios';
import {isAuthenticated} from "./stores/authStore";

const API_BASE_URL = 'http://localhost:5001';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json'
    }
});

export const textToSqlService = {
    convertQuery: async (naturalLanguageQuery: string) => {
        const response = await apiClient.post('/text-to-sql/convert', { query: naturalLanguageQuery });
        return response.data;
    }
};

export const customQueryService = {
    executeQuery: async (sql: string) => {
        const response = await apiClient.post('/database/execute-custom-query', { sql });
        return response.data;
    }
};

export const authorization = {
    login: async (email: string, password: string, rememberMe: boolean) => {
        try {
            const response = await apiClient.post('/auth/login', { email, password });
            const { token } = response.data;

            if (!token) throw new Error('Token not received');

            if (rememberMe) {
                document.cookie = `token=${token};path=/;max-age=${60 * 60 * 24 * 30}`;
            } else {
                document.cookie = `token=${token};path=/`;
            }
            return response.data;
        } catch (error: any) {
            console.error('Login error:', error);
            alert('Login failed: ' + error.message);
            throw error;
        }
    },

    register: async (name: string, email: string, password: string) => {
        try {
            const response = await apiClient.post('/auth/register', { username: name, email, password });
            return response.data;
        } catch (error: any) {
            if (error.status === 409) {
                alert('User with this email already exists');
                return;
            }
            console.error('Registration error:', error);
            alert('Registration failed: ' + error.message);
            throw Error('User with this email already exists');
        }
    },
    logout: () => {
        document.cookie = 'token=;path=/;max-age=0';
        isAuthenticated.set(false);
        alert('You have been logged out.');
    }
};

export function getAuthToken() {
    const cookies = document.cookie.split('; ');
    const tokenCookie = cookies.find(row => row.startsWith('token='));
    return tokenCookie ? tokenCookie.split('=')[1] : null;
}

apiClient.interceptors.request.use((config) => {
    const token = getAuthToken();
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
