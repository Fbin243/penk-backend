import express from "express";
import { readFileSync } from "fs";

import { OAuthTokenModel } from "../../utils/database/mongo";
import { getRedisClient } from "../../utils/database/redis";
import { getAuthResult } from "../../utils/googleapis";
import { LinkedAccount, LinkedAccountType } from "../../utils/types";
import { encrypt } from "./../../utils/encrypt";

const app = express();
const PORT = 8097;

const authSuccessHtml = readFileSync("resources/html/auth-success.html", "utf8");
const authFailureHtml = readFileSync("resources/html/auth-failure.html", "utf8");

const getTypeFromScope = (scope: string) => {
  if (scope.includes("gmail")) {
    return LinkedAccountType.Gmail;
  }

  if (scope.includes("calendar")) {
    return LinkedAccountType.GoogleCalendar;
  }

  throw new Error(`Invalid scope: ${scope}`);
};

// Route to handle OAuth redirect from Google
app.get("/oauth_redirect", async (req, res) => {
  const { state, code, scope } = req.query;

  if (!code) {
    return res.status(400).send("Authorization code not found");
  }

  if (!state) {
    return res.status(400).send("Linking account not found");
  }

  if (!scope) {
    return res.status(400).send("Scope not found");
  }

  try {
    const authResult = await getAuthResult(code as string);

    if (!authResult.refreshToken) {
      return res.status(400).send("Failed to link account");
    }

    const newAccount = await OAuthTokenModel.findOneAndUpdate(
      { profile_id: state, email: authResult.email, type: getTypeFromScope(scope.toString()) },
      { refresh_token: encrypt(authResult.refreshToken) },
      { upsert: true, new: true },
    );

    // Update the Redis cache
    const redisClient = await getRedisClient();
    const currentCache = await redisClient.get(`linked_accounts:${state}`);
    if (currentCache) {
      const linkedAccounts = JSON.parse(currentCache) as LinkedAccount[];
      linkedAccounts.push({
        id: newAccount._id.toString(),
        email: newAccount.email,
        type: newAccount.type,
        accessToken: authResult.accessToken || "",
      });
      await redisClient.set(`linked_accounts:${state}`, JSON.stringify(linkedAccounts));
    }

    res.send(authSuccessHtml);
  } catch (error) {
    console.error("Error handling OAuth redirect:", error);
    res.send(authFailureHtml);
  }
});

export const startOAuthServer = () => {
  app.listen(PORT, () => {
    console.log(`OAuth server running at http://localhost:${PORT}`);
  });
};
