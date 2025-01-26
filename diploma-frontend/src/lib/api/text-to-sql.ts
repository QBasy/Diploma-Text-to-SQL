import api from './index';

interface TextToSQLRequest {
    query: string;
    schema?: any;
}

interface TextToSQLResponse {
    sql: string;
}

export const generateSimpleSQL = async (query: string): Promise<TextToSQLResponse> => {
    return api.post('/text-to-sql/simple', { query });
};

export const generateComplexSQL = async (query: string, schema: any): Promise<TextToSQLResponse> => {
    return api.post('/text-to-sql/complex', { query, schema });
};