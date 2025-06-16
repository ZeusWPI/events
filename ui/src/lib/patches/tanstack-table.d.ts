import "@tanstack/react-table";

declare module "@tanstack/react-table" {
  interface ColumnMeta {
    small?: boolean;
    hideHeader?: boolean;
    colSpanHeader?: number;
    sticky?: boolean;
  }
}
