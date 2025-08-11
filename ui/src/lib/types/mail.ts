import { API } from "./api";

import { z } from "zod";
import { Base, JSONBody } from "./general";

export interface Mail extends Base {
  yearId: number;
  eventIds: number[];
  author_id: number;
  title: string;
  content: string;
  sendTime: Date;
  send: boolean;
  error?: string;
}

export function convertMailToModel(mail: API.Mail): Mail {
  return {
    id: mail.id,
    yearId: mail.year_id,
    eventIds: mail.event_ids,
    author_id: mail.author_id,
    title: mail.title,
    content: mail.content,
    sendTime: new Date(mail.send_time),
    send: mail.send,
    error: mail.error,
  }
}

export function convertMailsToModel(mails: API.Mail[]): Mail[] {
  return mails.map(convertMailToModel)
}

export const mailSchema = z.object({
  id: z.number().optional(),
  yearId: z.number().positive(),
  eventIds: z.array(z.number().positive()),
  title: z.string().nonempty(),
  content: z.string().nonempty(),
  sendTime: z.date().min(new Date()),
})
export type MailSchema = z.infer<typeof mailSchema> & JSONBody;
