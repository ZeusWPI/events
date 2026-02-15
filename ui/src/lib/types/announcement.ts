import { API } from "./api";
import { Base, JSONBody } from "./general";
import { z } from "zod";

export interface Announcement extends Base {
  yearId: number;
  eventIds: number[];
  author_id: number;
  content: string;
  sendTime: Date;
  send: boolean;
  draft: boolean;
  error?: string;
}

export function convertAnnouncementToModel(announcement: API.Announcement): Announcement {
  return {
    id: announcement.id,
    yearId: announcement.year_id,
    eventIds: announcement.event_ids,
    author_id: announcement.author_id,
    content: announcement.content,
    sendTime: new Date(announcement.send_time),
    draft: announcement.draft,
    send: announcement.send,
    error: announcement.error
  }
}

export function convertAnnouncementsToModel(announcements: API.Announcement[]): Announcement[] {
  return announcements.map(convertAnnouncementToModel)
}

export const announcementSchema = z.object({
  id: z.number().optional(),
  yearId: z.number().positive(),
  eventIds: z.array(z.number().positive()),
  content: z.string().nonempty(),
  sendTime: z.date().min(new Date()).optional(),
  draft: z.boolean(),
})
export type AnnouncementSchema = z.infer<typeof announcementSchema> & JSONBody
