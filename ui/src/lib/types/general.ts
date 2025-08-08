export interface JSON {
  [key: string]: string | number | boolean | Date | JSON | JSON[] | unknown;
}

export type JSONBody = JSON | JSON[];

export interface Base extends JSON {
  id: number;
}

export const weightCategory = 10;
export const weightSubcategory = 20;
export const weightItem = 30;
