import api from './index';

export async function getHistory(queryString: string = ''): Promise<PaginatedResponse> {
    const url = queryString ? `/history?${queryString}` : '/history';
    return await api.get(url);
}

export const clearHistory = async (): Promise<boolean> => {
    const response = await api.delete('/history');
    return response.status === 200;
};

export interface QueryHistory {
    id: string;
    user_id: string;
    database_uuid: string;
    query_type: 'SQL' | 'Natural Language';
    query: string;
    result: any;
    timestamp: string;
    success: boolean;
}

export interface PaginatedResponse {
    data: QueryHistory[];
    page: number;
    per_page: number;
    total: number;
    last_page: number;
}

export async function createHistoryEntry(entry: {
    endpoint: string;
    timestamp: string;
    query: any;
    result: any;
    success: boolean;
}) {

    const payload = {
        query_type: entry.endpoint.includes('text-to-sql')
            ? 'text-to-sql'
            : entry.endpoint.includes('execute-sql')
                ? 'execute-sql'
                : 'database',
        query: typeof entry.query === 'string' ? entry.query : JSON.stringify(entry.query),
        result: typeof entry.result === 'string' ? entry.result : JSON.stringify(entry.result),
        timestamp: entry.timestamp,
        success: entry.success
    };

    await api.post('/history', payload);
    console.log("Added to history")
}
