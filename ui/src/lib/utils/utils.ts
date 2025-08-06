import type { ClassValue } from "clsx";
import { clsx } from "clsx";
import { format } from "date-fns";
import { twMerge } from "tailwind-merge";
import { v4 as uuid } from 'uuid'

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
    return obj.toISOString()
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

export function formatDate(date?: Date) {
  if (date === undefined) {
    return ""
  }

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

export function getBuildTime() {
  const buildTime = import.meta.env.VITE_BUILD_TIME as string | "";

  return buildTime ? new Date(buildTime) : undefined
}

export function randomNumber() {
  return Math.floor(Math.random() * (100000 - 10000)) + 10000;
}

export function getUuid() {
  return uuid()
}

export function arrayEqual<T extends ArrayLike<unknown>>(a: T, b: T): boolean {
  if (a.length !== b.length) {
    return false;
  }

  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) {
      return false;
    }
  }

  return true;
}

export function capitalize(text: string): string {
  if (!text) {
    return text
  }

  return text.charAt(0).toUpperCase() + text.substring(1).toLowerCase()
}


const A4_ASPECT_RATIO = 1.4142;
const A4_ASPECT_RATIO_TOLERANCE = 0.01;
export function isA4AspectRatio(file: File): Promise<boolean> {
  return new Promise((resolve, reject) => {
    if (!file.type.startsWith("image/")) {
      resolve(false);
      return;
    }

    const img = new Image();
    const url = URL.createObjectURL(file);

    img.onload = () => {
      const { width, height } = img;
      const aspectRatio = height / width;
      URL.revokeObjectURL(url);

      const isA4 = Math.abs(aspectRatio - A4_ASPECT_RATIO) <= A4_ASPECT_RATIO_TOLERANCE;
      resolve(isA4);
    };

    img.onerror = (err) => {
      URL.revokeObjectURL(url);
      reject(err);
    };

    img.src = url;
  });
}
