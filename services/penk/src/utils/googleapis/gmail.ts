import { google } from "googleapis";

import { oauth2Client } from "./config";

export const getMails = async (props: { accessToken: string; q: string }) => {
  oauth2Client.setCredentials({ access_token: props.accessToken });
  const gmail = google.gmail({ version: "v1", auth: oauth2Client });
  const res = await gmail.users.messages.list({
    userId: "me",
    q: props.q,
    maxResults: 5,
  });

  const messages = res.data.messages || [];

  const messageDetails = await Promise.all(
    messages.map(async (message) => {
      if (!message.id) {
        return null;
      }

      const res = await gmail.users.messages.get({ userId: "me", id: message.id });
      const headers = res.data.payload?.headers || [];
      const getHeader = (name: string) => headers.find((h) => h.name === name)?.value || "";

      let to = getHeader("To")
        .split(",")
        .map((email) => email.trim());

      if (to.length > 3) {
        to = [...to.slice(0, 3), `+${to.length - 3} more`];
      }

      return {
        id: message.id,
        from: getHeader("From"),
        to,
        subject: getHeader("Subject"),
        date: getHeader("Date"),
        snippet: res.data.snippet,
        labelIds: res.data.labelIds,
        isUnread: res.data.labelIds?.includes("UNREAD") || false,
        hasAttachment: res.data.payload?.parts?.some((part) => part.filename) || false,
      };
    }),
  );

  return messageDetails.filter((m) => m !== null);
};

export const getMailBody = async (props: { accessToken: string; messageId: string }) => {
  oauth2Client.setCredentials({ access_token: props.accessToken });
  const gmail = google.gmail({ version: "v1", auth: oauth2Client });
  const res = await gmail.users.messages.get({ userId: "me", id: props.messageId });
  const body = res.data.payload?.parts?.map((part) => part.body?.data).join("");
  return body;
};
