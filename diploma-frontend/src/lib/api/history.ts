import api from './index';

export interface QueryHistory {
    id: string;
    database_uuid: string;
    query_type: 'SQL' | 'Natural Language';
    query: string;
    result: any;
    timestamp: string;
}

export const getHistory = async (): Promise<QueryHistory[]> => {
    return api.get('/history');
};

export const clearHistory = async (): Promise<boolean> => {
    const response = await api.delete('/history');
    return response.status === 200;

};

export async function createHistoryEntry(entry: {
    endpoint: string;
    timestamp: string;
    query: any;
    result: any;
    success: boolean;
}) {

    const payload = {
        query_type: entry.endpoint.includes('text-to-sql') ? 'text-to-sql' : 'database',
        query: typeof entry.query === 'string' ? entry.query : JSON.stringify(entry.query),
        result: typeof entry.result === 'string' ? entry.result : JSON.stringify(entry.result),
        timestamp: entry.timestamp,
        success: entry.success
    };

    await api.post('/history', payload);
    console.log("Added to history")
}
