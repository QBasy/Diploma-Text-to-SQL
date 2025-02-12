export interface Column {
    name: string;
    type: string;
    isForeignKey: boolean;
    referencedTable: string;
    referencedColumn: string;
}

export interface Table {
    name: string;
    columns: Column[];
    primaryKey: string;
}

export interface Schema {
    tables: Table[]
}