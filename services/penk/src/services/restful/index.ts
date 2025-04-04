import express from "express";

import { OAuthTokenModel } from "../../utils/database/mongo";
import { getAuthResult } from "../../utils/googleapis";
import { encrypt } from "./../../utils/encrypt";

const app = express();
const PORT = 8097;

// Route to handle OAuth redirect from Google
app.get("/oauth_redirect", async (req, res) => {
  const { state, code, scope } = req.query;

  if (!code) {
    return res.status(400).send("Authorization code not found");
  }

  if (!state) {
    return res.status(400).send("Linking account not found");
  }

  try {
    const authResult = await getAuthResult(code as string);

    if (!authResult.refreshToken) {
      return res.status(400).send("Failed to link account");
    }

    await OAuthTokenModel.findOneAndUpdate(
      { profile_id: state, email: authResult.email, scope },
      { refresh_token: encrypt(authResult.refreshToken) },
      { upsert: true },
    );

    res.send(`
      <h1>PenK Assistant: Authentication successful!</h1>
      <p>Account: ${authResult.email}</p>
      <p>You can now close this window.</p>
    `);
  } catch (error) {
    console.error("Error handling OAuth redirect:", error);
    res.status(500).send("Authentication failed");
  }
});

export const startOAuthServer = () => {
  app.listen(PORT, () => {
    console.log(`OAuth server running at http://localhost:${PORT}`);
  });
};
