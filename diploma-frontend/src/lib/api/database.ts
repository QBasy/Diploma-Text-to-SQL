import api from './index';

interface QueryRequest {
    query: string;
}

interface QueryResponse {
    result: any;
}

interface VisualisationResponse {
    svg: string,
    result: any,
    columns: any,
    row_count: number
}

export const createTableAPI = async (schema: TableSchema): Promise<void> => {
    return api.post('/api/database/tables', schema);
};

export const executeQueryAPI = async (query: string): Promise<QueryResponse> => {
    return api.post('/api/database/execute-sql', { query });
};

export const executeQueryVisualisation = async (query: string): Promise<VisualisationResponse> => {
    return api.post('/api/database/visualise', { query });
}

export const getSchemaVisualisationSvg = async (query: string): Promise<{ svg: string }> => {
    const response = await api.post('/api/database/visualise', { query });
    return response.data;
};

export interface ColumnInfo {
    name: string;
    type: string;
    isForeignKey: boolean;
    referencedTable?: string;
    referencedColumn?: string;
}

export interface TableSchema {
    columns: ColumnInfo[];
    primaryKey: string;
    foreignKeys: { table: string; from: string; to: string }[];
}

export type Schema = Record<string, TableSchema>;

export interface SchemaResponse {
    status: "success" | "error";
    schema?: Record<string, TableSchema>;
    message?: string;
}

export const getSchema = async (): Promise<SchemaResponse> => {
    return await api.get('/api/database/schema');
}

export const getComplexSchema = async (): Promise<Schema> => {
    return api.get('/api/database/schema-complex').then(res => res.data);
};
