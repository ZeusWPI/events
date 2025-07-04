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

  if (obj instanceof Date) {
    console.log("Date")
    return obj.toISOString()
  }

  if (typeof obj === "object") {
    console.log(obj)
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

export function formatDateDiff(first: Date, second: Date) {
  const diff = first.getTime() - second.getTime()
  const sign = diff >= 0 ? "-" : "+"
  const absDiff = Math.abs(diff)

  const totalHours = Math.floor(absDiff / (1000 * 60 * 60))
  const days = Math.floor(totalHours / 24)
  const hours = totalHours % 24

  return `${sign} ${days} days ${hours} hours`
}

export function randomNumber() {
  return Math.floor(Math.random() * (100000 - 10000)) + 10000;
}
