import { z } from "zod";

export const mailSearch = z.object({
  eventIds: z.array(z.number().positive()).optional(),
})

