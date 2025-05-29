import api from './index';

export type CustomDatabase = {
    uuid: string;
    name: string;
    db_type: string;
    host: string;
    port: number;
    database: string;
    description: string;
};

export type AddDatabasePayload = {
    name: string;
    db_type: string;
    host: string;
    port: number;
    username: string;
    password: string;
    database: string;
    ssl_mode: string;
    description: string;
};


export async function listCustomDatabases(): Promise<CustomDatabase[]> {
    return await api.get('/api/database/custom/list');

}

export async function addCustomDatabase(payload: AddDatabasePayload) {
    return api.post('/api/database/custom/add', payload);
}

export async function deleteCustomDatabase(uuid: string) {
    return api.delete('/api/database/custom/delete', {
        data: { uuid },
    });
}

export async function getCustomSchema(uuid: string) {
    return api.get('/api/database/custom/schema', {
        params: { database_uuid: uuid },
    });
}

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

export async function getCustomSchemaComplex(uuid: string): Promise<Record<string, TableSchema>> {
    const request = await api.get('/api/database/custom/schema-complex', {
        params: { database_uuid: uuid },
    });

    const response = request.data;

    if (response && response.schema && Array.isArray(response.schema)) {
        const transformedSchema: Record<string, TableSchema> = {};
        response.schema.forEach(table => {
            transformedSchema[table.name] = {
                name: table.name,
                columns: table.columns,
                primaryKey: table.primaryKey
            };
        });
        return transformedSchema;
    }

    throw new Error("Invalid schema format returned from server");
}
