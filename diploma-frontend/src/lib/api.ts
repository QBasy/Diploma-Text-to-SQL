import axios from 'axios';

const API_BASE_URL = 'http://localhost:5001';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json'
    }
});

export const textToSqlService = {
    convertQuery: async (naturalLanguageQuery: string) => {
        const response = await apiClient.post('/text-to-sql', { query: naturalLanguageQuery });
        return response.data;
    }
};

export const customQueryService = {
    executeQuery: async (sql: string) => {
        const response = await apiClient.post('/custom-query', { sql });
        return response.data;
    }
};

export const itemService = {
    getItems: async () => {
        const response = await apiClient.get('/items');
        return response.data;
    },
    createItem: async (item: any) => {
        const response = await apiClient.post('/items', item);
        return response.data;
    }
};

export const authorization = {
    login: async (email: string, password: string, rememberMe: boolean) => {
        try {
            const response = await apiClient.post('/auth/login', { email, password, rememberMe });
            const { token } = response.data;

            if (rememberMe) {
                document.cookie = `token=${token};path=/;max-age=${60 * 60 * 24 * 30}`;
            } else {
                sessionStorage.setItem('token', token);
            }
            return response.data;
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    },
    register: async (name: string, email: string, password: string) => {
        try {
            const response = await apiClient.post('/auth/register', { name, email, password });
            return response.data;
        } catch (error) {
            console.error('Registration error:', error);
            throw error;
        }
    }
};