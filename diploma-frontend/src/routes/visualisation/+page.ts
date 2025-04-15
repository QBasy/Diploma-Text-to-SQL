import type { PageLoad } from './$types';
import { getComplexSchema } from '$lib/api';

export const load: PageLoad = async ({ fetch }) => {
    try {
        const schema = await getComplexSchema();

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
