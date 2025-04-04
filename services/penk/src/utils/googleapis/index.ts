// https://github.com/googleapis/google-api-nodejs-client
import { google } from "googleapis";

const oauth2Client = new google.auth.OAuth2(
  `${process.env.GOOGLE_OAUTH_APP_GUID}.apps.googleusercontent.com`,
  process.env.WEB_CLIENT_SECRET,
  `http://localhost:8097/oauth_redirect`,
);

export const getGoogleAuthUrl = async (profile_id: string, scope: string) => {
  return oauth2Client.generateAuthUrl({
    access_type: "offline",
    scope: ["https://www.googleapis.com/auth/userinfo.email", scope],
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
