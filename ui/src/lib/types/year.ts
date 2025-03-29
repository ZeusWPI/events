import type { API } from "./api";

export interface Year {
  id: number;
  startYear: number;
  endYear: number;
  formatted: string;
}

export function convertYearToModel(year: API.Year): Year {
  const startFormatted = (year.start_year % 100).toString().padStart(2, "0");
  const endFormatted = (year.end_year % 100).toString().padStart(2, "0");

  return {
    id: year.id,
    startYear: year.start_year,
    endYear: year.end_year,
    formatted: `${startFormatted}-${endFormatted}`,
  };
}

export function convertYearsToModel(years: API.Year[]): Year[] {
  return years.map(convertYearToModel);
}

export function convertYearToJSON(year: Year): API.Year {
  return {
    id: year.id,
    start_year: year.startYear,
    end_year: year.endYear,
  };
}
