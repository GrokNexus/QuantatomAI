declare module 'apache-arrow' {
    export class Table {
        numRows: number;
        schema: {
            fields: Array<{ name: string }>;
        };
        getChild(name: string): any;
    }
    export function tableFromArrays(data: Record<string, any[]>): Table;
    export class RecordBatchReader {
        static from(data: any): any;
    }
}
