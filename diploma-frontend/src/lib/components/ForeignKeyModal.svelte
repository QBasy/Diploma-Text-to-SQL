<script lang="ts">
    import type { Column, Table } from "$lib/types/table.js";

    export let table: ExtendedTable[]  = []

    type PrimaryKeyType = 'INTEGER_AUTOINCREMENT' | 'UUID' | 'INTEGER' | 'TEXT';

    interface ExtendedColumn extends Column {
        isPrimaryKey?: boolean;
        primaryKeyType?: PrimaryKeyType;
    }

    interface ExtendedTable extends Omit<Table, 'columns'> {
        columns: ExtendedColumn[];
    }
</script>