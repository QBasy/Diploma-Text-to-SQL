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
