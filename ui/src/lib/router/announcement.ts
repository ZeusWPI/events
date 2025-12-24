import { z } from "zod";

export const announcementSearch = z.object({
  eventIds: z.array(z.number().positive()).optional(),
})

