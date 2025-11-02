import { useMutation } from "@tanstack/react-query"
import { apiPut, NO_CONVERTER } from "./query"

const ENDPOINT = "image"

export function useImageCreate() {
  return useMutation({
    mutationFn: async (args: { name: string, file: File }) => (await apiPut<{ id: number }>(ENDPOINT, { name: args.name }, NO_CONVERTER, [{ file: args.file, field: "file" }])).data,
  })
}
