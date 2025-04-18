import api from './index';
import type {Schema, Table} from "$lib/types/table"

export interface TextToSQLRequest {
    sql: string;
    schema?: any;
}

interface TextToSQLResponse {
    sql: string;
}

export const generateSimpleSQL = async (query: string): Promise<TextToSQLResponse> => {
    return api.post('/text-to-sql/simple', { query });
};

export const generateComplexSQL = async (query: string, schema: Schema[]): Promise<TextToSQLResponse> => {
    return api.post('/text-to-sql/complex', { query, schema });
};

export const generateComplexSQLbyGPT = async (query: string, schema: Schema): Promise<TextToSQLResponse> => {
    console.log(schema)
    return api.post('/text-to-sql/groc', { query, schema });
};