import api from './index';

interface TableSchema {
    name: string;
    columns: { name: string; type: string }[];
}

interface QueryRequest {
    query: string;
}

interface QueryResponse {
    result: any;
}

export const createTable = async (schema: TableSchema): Promise<void> => {
    return api.post('/database/tables', schema);
};

export const executeQuery = async (query: string): Promise<QueryResponse> => {
    return api.post('/database/query', { query });
};

export const getSchema = async (): Promise<TableSchema[]> => {
    return api.get('/database/schema');
};