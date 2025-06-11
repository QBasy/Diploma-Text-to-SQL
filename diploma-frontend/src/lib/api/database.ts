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
export interface ColumnInfo {
    name: string;
    type: string;
    isForeignKey: boolean;
    referencedTable?: string;
    referencedColumn?: string;
}

export interface TableSchema {
    name: string;
    columns: ColumnInfo[];
    primaryKey?: string;
}

export const createTableAPI = async (schema: TableSchema): Promise<void> => {
    return api.post('/api/database/tables', schema);
};

export const executeQueryAPI = async (query: string): Promise<QueryResponse> => {
    return api.post('/api/database/execute-sql', { query });
};

export const executeQueryVisualisation = async (query: string, databaseUUID: string = ""): Promise<VisualisationResponse> => {
    return api.post('/api/database/visualise', { query, databaseUUID });
}

export const getSchemaVisualisationSvg = async (query: string, databaseUUID: string = ""): Promise<{ svg: string }> => {
    const response = await api.post('/api/database/visualise', { query, databaseUUID });
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


export const getComplexSchema = async (): Promise<Record<string, TableSchema>> => {
    const res = await api.get('/api/database/schema-complex');

    return res.schema;
};

