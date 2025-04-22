import type { PageLoad } from '../../../../.svelte-kit/types/src/routes';
import { getComplexSchema } from '$lib/api';

export const load: PageLoad = async ({ fetch }) => {
    try {
        const schema = await getComplexSchema();
        if (!schema) {
            throw new Error("Empty schema");
        }

        console.log(schema)
        return {
            schema: schema
        };
    } catch (error) {
        console.error("Error loading visualization page data:", error);
        return {
            schema: {},
            error: "Failed to load schema information"
        };
    }
};
