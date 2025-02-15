    import api from './index';

    export interface TableSchema {
        name: string;
        columns: { name: string; type: string }[];
    }

    interface QueryRequest {
        query: string;
    }

    interface QueryResponse {
        result: any;
    }

    export const createTableAPI = async (schema: TableSchema): Promise<void> => {
        return api.post('/database/tables', schema);
    };

    export const executeQueryAPI = async (query: string): Promise<QueryResponse> => {
        return api.post('/database/execute-sql', { query });
    };

    export const getSchema = async (): Promise<TableSchema[]> => {
        return api.get('/database/schema');
    };
