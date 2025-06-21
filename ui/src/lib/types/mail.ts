import { API } from "./api";
import { Event, convertEventsToModel } from "./event";
import { Base } from "./general";

export interface Mail extends Base {
  title: string;
  content: string;
  sendTime: Date;
  send: boolean;
  events: Event[];
  error?: string;
}

export function convertMailToModel(mail: API.Mail): Mail {
  return {
    id: mail.id,
    title: mail.title,
    content: mail.content,
    sendTime: new Date(mail.send_time),
    send: mail.send,
    events: convertEventsToModel(mail.events),
    error: mail.error,
  }
}

export function convertMailsToModel(mails: API.Mail[]): Mail[] {
  return mails.map(convertMailToModel)
}

