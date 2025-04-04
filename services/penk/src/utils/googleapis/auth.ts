import { google } from "googleapis";

import { LinkedAccountType } from "../types";
import { oauth2Client } from "./config";

const getScopeFromType = (type: LinkedAccountType) => {
  switch (type) {
    case LinkedAccountType.Gmail:
      return "https://www.googleapis.com/auth/gmail.readonly";
    case LinkedAccountType.GoogleCalendar:
      return "https://www.googleapis.com/auth/calendar";
    default:
      throw new Error(`Invalid linked account type: ${type}`);
  }
};

export const getGoogleAuthUrl = async (profile_id: string, type: LinkedAccountType) => {
  return oauth2Client.generateAuthUrl({
    access_type: "offline",
    scope: ["https://www.googleapis.com/auth/userinfo.email", getScopeFromType(type)],
    prompt: "consent",
    state: profile_id,
  });
};

export const getAuthResult = async (code: string) => {
  const { tokens } = await oauth2Client.getToken(code);
  oauth2Client.setCredentials(tokens);
  const oauth2 = google.oauth2({
    version: "v2",
    auth: oauth2Client,
  });
  const { data } = await oauth2.userinfo.get();
  return {
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token,
    email: data.email,
  };
};

export const refreshToken = async (refreshToken: string) => {
  oauth2Client.setCredentials({ refresh_token: refreshToken });
  const res = await oauth2Client.getAccessToken();
  return res.token;
};
