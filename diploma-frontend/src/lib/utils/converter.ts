import type { Schema as APISchema, TableSchema } from '$lib/api';
import type { Schema as GptSchema, Table, Column } from '$lib/types/table';

export function convertSchemaToGptFormat(apiSchema: APISchema): GptSchema {
    const tables: Table[] = Object.entries(apiSchema).map(([tableName, tableSchema]: [string, TableSchema]) => {
        const columns: Column[] = tableSchema.columns.map(col => ({
            name: col.name,
            type: col.type,
            isForeignKey: col.isForeignKey,
            referencedTable: col.referencedTable || '',
            referencedColumn: col.referencedColumn || ''
        }));

        return {
            name: tableName,
            columns,
            primaryKey: tableSchema.primaryKey
        };
    });

    return { tables };
}
