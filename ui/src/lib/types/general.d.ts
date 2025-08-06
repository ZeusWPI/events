export interface JSON {
  [key: string]: string | number | boolean | Date | JSON | JSON[] | unknown;
}

export type JSONBody = JSON | JSON[];

export interface Base extends JSON {
  id: number;
}

