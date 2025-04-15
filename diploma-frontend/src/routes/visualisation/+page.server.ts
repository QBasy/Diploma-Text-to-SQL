import { getComplexSchema } from '$lib/api';
import { convertSchemaToGptFormat } from '$lib/utils/converter';
import { generateComplexSQLbyGPT } from '$lib/api';
import type { Actions } from './$types'

export const actions: Actions = {
    convertQuery: async ({ request }) => {
        const data = await request.formData();
        const query = data.get('query')?.toString() || '';

        try {
            const rawSchema = await getComplexSchema();
            const convertedSchema = convertSchemaToGptFormat(rawSchema);

            const result = await generateComplexSQLbyGPT(query, convertedSchema);

            return {
                success: true,
                sql: result.sql
            };
        } catch (error) {
            console.error('Text-to-SQL conversion failed:', error);
            return {
                success: false,
                error: 'Failed to process query'
            };
        }
    }
};

