import { google } from "googleapis";

import { oauth2Client } from "./config";

export const getCalendarEvents = async (props: {
  accessToken: string;
  timeMin: Date;
  timeMax: Date;
}) => {
  oauth2Client.setCredentials({ access_token: props.accessToken });
  const calendar = google.calendar({ version: "v3", auth: oauth2Client });
  const res = await calendar.events.list({
    calendarId: "primary",
    timeMin: props.timeMin.toISOString(),
    timeMax: props.timeMax.toISOString(),
  });
  return res.data.items;
};
