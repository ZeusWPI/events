import { API } from "./api";
import { Base } from "./general";

export interface Mail extends Base {
  content: string;
  sendTime: Date;
  send: boolean;
  error?: string;
}

export function convertMailToModel(mail: API.Mail): Mail {
  return {
    id: mail.id,
    content: mail.content,
    sendTime: new Date(mail.send_time),
    send: mail.send,
    error: mail.error,
  }
}

export function convertMailsToModel(mails: API.Mail[]): Mail[] {
  return mails.map(convertMailToModel)
}

