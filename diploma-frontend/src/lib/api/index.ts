import axios from 'axios';
import { createHistoryEntry } from './history';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:5001/';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});


api.interceptors.response.use(
    async (response) => {
        const url = response.config.url;

        if (url?.includes('/database') || url?.includes('/text-to-sql')) {
            try {
                await createHistoryEntry({
                    endpoint: url,
                    timestamp: new Date().toISOString(),
                    query: response.config.data,
                    result: response.data,
                    success: response.status >= 200 && response.status < 300
                });
            } catch (err) {
                console.error('Ошибка при создании истории:', err);
            }
        }

        return response.data;
    },
    (error) => {
        const url = error.config?.url;

        if (url?.includes('/database/execute-sql') || url?.includes('/text-to-sql')) {
            createHistoryEntry({
                endpoint: url,
                timestamp: new Date().toISOString(),
                query: error.config?.data,
                result: error.response?.data || { message: 'Unknown error' },
                success: false
            }).catch((err) => {
                console.error('Ошибка при создании истории (ошибка запроса):', err);
            });
        }

        if (error.response) {
            throw new Error(error.response.data.message || 'An error occurred');
        } else {
            throw new Error('Network error');
        }
    }
);


export default api;

export * from './auth';
export * from './database';
export * from './text-to-sql';
export * from './history';
