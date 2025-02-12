import api from './index';

interface QueryHistory {
    id: string;
    type: 'SQL' | 'Natural Language';
    query: string;
    result: any;
    timestamp: string;
}

export const getHistory = async (): Promise<QueryHistory[]> => {
    return api.get('/history');
};

export const clearHistory = async (): Promise<void> => {
    return api.delete('/history');
};
