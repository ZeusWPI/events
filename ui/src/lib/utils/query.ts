import type { JSONBody } from "../types/general";
import { camelToSnake } from "./utils";

const URLS: Record<string, string> = {
  API: import.meta.env.VITE_BACKEND_URL as string,
};

export async function getApi<T, U = unknown>(endpoint: string, convertData?: (data: U) => T) {
  return _fetch<T, U>(
    `${URLS.API}/${endpoint}`,
    {},
    convertData,
  );
}

export async function postApi<T, U = unknown>(endpoint: string, body: JSONBody = {}, convertData?: (data: U) => T) {
  return _fetch<T, U>(
    `${URLS.API}/${endpoint}`,
    {
      method: "POST",
      body: JSON.stringify(camelToSnake(body)),
      headers: new Headers({ "content-type": "application/json" }),
    },
    convertData,
  );
}

export async function putApi<T, U = unknown>(endpoint: string, body: JSONBody = {}, convertData?: (data: U) => T) {
  return _fetch<T, U>(
    `${URLS.API}/${endpoint}`,
    {
      method: "PUT",
      body: JSON.stringify(camelToSnake(body)),
      headers: new Headers({ "content-type": "application/json" }),
    },
    convertData,
  );
}

export async function patchApi<T, U = unknown>(endpoint: string, body: JSONBody = {}, convertData?: (data: U) => T) {
  return _fetch<T, U>(
    `${URLS.API}/${endpoint}`,
    {
      method: "PATCH",
      body: JSON.stringify(camelToSnake(body)),
      headers: new Headers({ "content-type": "application/json" }),
    },
    convertData,
  );
}

export async function deleteApi<T, U = unknown>(endpoint: string, convertData?: (data: U) => T) {
  return _fetch<T, U>(
    `${URLS.API}/${endpoint}`,
    {
      method: "DELETE",
    },
    convertData,
  );
}

interface ResponseNot200Error extends Error {
  response: Response;
}

export function isResponseNot200Error(error: unknown): error is ResponseNot200Error {
  return (error as ResponseNot200Error).response !== undefined;
}

async function _fetch<T, U>(url: string, options: RequestInit = {}, convertData?: (data: U) => T): Promise<T> {
  return fetch(
    url,
    { credentials: "include", ...options },
  ).then(async (response) => {
    if (!response.ok) {
      const error = new Error(`Fetch failed with status: ${response.status}`) as ResponseNot200Error;
      error.response = response;
      throw error;
    }

    const contentType = response.headers.get("content-type");

    if (contentType?.includes("application/json"))
      return response.json() as Promise<unknown>;
    else if (contentType?.includes("image/png"))
      return response.blob();
    else
      return response.text();
  }).then(data => (convertData ? convertData(data as U) : (data as T)));
}
