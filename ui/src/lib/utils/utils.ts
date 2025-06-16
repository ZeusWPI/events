import type { ClassValue } from "clsx";
import { clsx } from "clsx";
import { format } from "date-fns";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function camelToSnake(obj: unknown): unknown {
  if (obj === null || obj === undefined) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(camelToSnake);
  }

  if (typeof obj === "object") {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [
        stringCamelToSnake(key),
        camelToSnake(value),
      ]),
    );
  }

  return obj;
}

function stringCamelToSnake(str: string) {
  return str.replace(/[A-Z]+/g, l => `_${l.toLowerCase()}`);
}

export function formatDate(date: Date) {
  return format(date, "eee dd MMM, HH:mm");
}

export function randomNumber() {
  return Math.floor(Math.random() * (100000 - 10000)) + 10000;
}
