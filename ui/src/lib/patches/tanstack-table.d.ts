import "@tanstack/react-table";

declare module "@tanstack/react-table" {
  interface ColumnMeta<TData extends RowData, TValue> {
    small?: boolean;
    hideHeader?: boolean;
    colSpanHeader?: number;
    sticky?: boolean;
  }
}
