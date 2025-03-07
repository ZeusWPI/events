export interface JSON {
  [key: string]: string | number | boolean | JSON | JSON[];
}

export type JSONBody = JSON | JSON[];
