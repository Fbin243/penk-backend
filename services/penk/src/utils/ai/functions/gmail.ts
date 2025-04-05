import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";

import { OAuthTokenModel } from "../../database/mongo";
import { decrypt } from "../../encrypt";
import { refreshToken } from "../../googleapis";
import { getMailBody, getMails } from "../../googleapis/gmail";
import { LinkedAccountType } from "../../types";
import { FunctionName } from "./types";

const getMailsDefinition: FunctionDefinition = {
  name: FunctionName.GetMails,
  description: "Get emails from the user's inbox",
  parameters: {
    type: "object",
    properties: {
      profileId: {
        type: "string",
        description: "User's profile ID",
      },
      q: {
        type: "string",
        description: "Gmail query",
      },
    },
    required: ["profileId", "q"],
    additionalProperties: false,
  },
  strict: true,
};

export const functionGetMails = async (props: { profileId: string; q: string }) => {
  console.log(chalk.cyan(`[Tool: ${FunctionName.GetMails}]`));
  console.dir(props, { depth: null, colors: true });
  console.log();

  const tokens = await OAuthTokenModel.find({
    profile_id: props.profileId,
    type: LinkedAccountType.Gmail,
  });

  const messageFetchPromises = tokens.map(async (token) => {
    try {
      const accessToken = await refreshToken(decrypt(token.refresh_token));

      if (!accessToken) {
        console.warn(`Failed to refresh token for account: ${token.email}`);
        return [];
      }

      const messages = await getMails({
        accessToken,
        q: props.q,
      });

      return messages || [];
    } catch (err) {
      console.error(`Error fetching messages for ${token.email}:`, err);
      return [];
    }
  });

  const allMessages = (await Promise.all(messageFetchPromises)).flat();

  console.log("All messages:", allMessages);

  return allMessages;
};

export const toolGetMails: ChatCompletionTool = {
  type: "function",
  function: getMailsDefinition,
};

const getMailBodyDefinition: FunctionDefinition = {
  name: FunctionName.GetMailBody,
  description: "Get the body of an email",
  parameters: {
    type: "object",
    properties: {
      profileId: {
        type: "string",
        description: "User's profile ID",
      },
      email: {
        type: "string",
        description: "Email address of the user",
      },
      messageId: {
        type: "string",
      },
    },
    required: ["profileId", "email", "messageId"],
    additionalProperties: false,
  },
  strict: true,
};

export const functionGetMailBody = async (props: {
  profileId: string;
  email: string;
  messageId: string;
}) => {
  console.log(chalk.cyan(`[Tool: ${FunctionName.GetMailBody}]`));
  console.dir(props, { depth: null, colors: true });
  console.log();

  const tokens = await OAuthTokenModel.find({
    profile_id: props.profileId,
    email: props.email,
    type: LinkedAccountType.Gmail,
  });

  const mailBody = await getMailBody({
    accessToken: decrypt(tokens[0].refresh_token),
    messageId: props.messageId,
  });

  return mailBody;
};

export const toolGetMailBody: ChatCompletionTool = {
  type: "function",
  function: getMailBodyDefinition,
};
