import type { API } from "./api";
import { Base } from "./general";

export interface Year extends Base {
  start: number;
  end: number;
  formatted: string;
}

export function convertYearToModel(year: API.Year): Year {
  const startFormatted = (year.start % 100).toString().padStart(2, "0");
  const endFormatted = (year.end % 100).toString().padStart(2, "0");

  return {
    id: year.id,
    start: year.start,
    end: year.end,
    formatted: `${startFormatted}-${endFormatted}`,
  };
}

export function convertYearsToModel(years: API.Year[]): Year[] {
  return years.map(convertYearToModel);
}
