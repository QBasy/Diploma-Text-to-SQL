import { writable } from 'svelte/store';
import { getSchema } from '$lib/api/database';

interface TableSchema {
    name: string;
    columns: { name: string; type: string }[];
}

export const schemaStore = writable<TableSchema[]>([]);

export const loadSchema = async () => {
    try {
        const schema = await getSchema();
        schemaStore.set(schema); // Сохраняем схему
    } catch (error) {
        throw error;
    }
};
