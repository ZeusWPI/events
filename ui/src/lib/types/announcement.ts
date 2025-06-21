import { API } from "./api";
import { Base } from "./general";

export interface Announcement extends Base {
  eventId: number;
  content: string;
  sendTime: Date;
  send: boolean;
  error?: string;
}

export function convertAnnouncementToModel(announcement: API.Announcement): Announcement {
  return {
    id: announcement.id,
    eventId: announcement.event_id,
    content: announcement.content,
    sendTime: new Date(announcement.send_time),
    send: announcement.send,
    error: announcement.error
  }
}

export function convertAnnouncementsToModel(announcements: API.Announcement[]): Announcement[] {
  return announcements.map(convertAnnouncementToModel)
}

