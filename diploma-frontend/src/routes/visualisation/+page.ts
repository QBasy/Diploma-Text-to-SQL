import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    try {
                const response = await fetch('/database/schema-complex');

        if (response.ok) {
            const data = await response.json();
            return {
                schema: data.schema || {}
            };
        }

        return {
            schema: {}
        };
    } catch (error) {
        console.error("Error loading visualization page data:", error);
        return {
            schema: {},
            error: "Failed to load schema information"
        };
    }
};

import type { Actions } from './$types';

export const actions: Actions = {
    convertQuery: async ({ request, fetch }) => {
        const data = await request.formData();
        const query = data.get('query')?.toString() || '';

        try {
            const response = await fetch('/text-to-sql/gpt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ query })
            });

            if (!response.ok) {
                return { success: false, error: 'Failed to convert query' };
            }

            const result = await response.json();
            return {
                success: true,
                sql: result.query
            };
        } catch (error) {
            return {
                success: false,
                error: 'Failed to process query'
            };
        }
    }
};