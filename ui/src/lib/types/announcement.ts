import { API } from "./api";

export interface Announcement {
  id: number;
  eventId: number;
  content: string;
  sendTime: Date;
  send: boolean;
}

export function convertAnnouncementToModel(announcement: API.Announcement): Announcement {
  return {
    id: announcement.id,
    eventId: announcement.event_id,
    content: announcement.content,
    sendTime: new Date(announcement.send_time),
    send: announcement.send,
  }
}

export function convertAnnouncementsToModel(announcements: API.Announcement[]): Announcement[] {
  return announcements.map(convertAnnouncementToModel)
}

export function convertAnnouncementToJSON(announcement: Announcement): API.Announcement {
  return {
    id: announcement.id,
    event_id: announcement.eventId,
    content: announcement.content,
    send_time: announcement.sendTime.toISOString(),
    send: announcement.send,
  }
}
