import { google } from "googleapis";

export const oauth2Client = new google.auth.OAuth2(
  `${process.env.GOOGLE_OAUTH_APP_GUID}.apps.googleusercontent.com`,
  process.env.WEB_CLIENT_SECRET,
  process.env.GOOGLE_OAUTH_REDIRECT_URI,
);
