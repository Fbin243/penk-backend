import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";

import { getLinkedAccounts } from "../../database/utils";
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

  const linkedAccounts = (await getLinkedAccounts(props.profileId)).filter(
    (linkedAccount) => linkedAccount.type === LinkedAccountType.Gmail,
  );

  if (linkedAccounts.length === 0) {
    throw new Error(
      'No linked Gmail accounts found. Please ask user to go to App Settings, click on the "Linked Accounts" button, and link their Gmail account.',
    );
  }

  const messageFetchPromises = linkedAccounts.map(async (linkedAccount) => {
    try {
      const messages = await getMails({
        accessToken: linkedAccount.accessToken,
        q: props.q,
      });

      return messages || [];
    } catch (err) {
      console.error(`Error fetching messages for ${linkedAccount.email}:`, err);
      return [];
    }
  });

  const allMessages = (await Promise.all(messageFetchPromises)).flat();
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

  const linkedAccounts = await getLinkedAccounts(props.profileId);

  const accessToken = linkedAccounts.find(
    (linkedAccount) => linkedAccount.email === props.email,
  )?.accessToken;

  if (!accessToken) {
    throw new Error(`No access token found for ${props.email}`);
  }

  const mailBody = await getMailBody({
    accessToken,
    messageId: props.messageId,
  });

  return mailBody;
};

export const toolGetMailBody: ChatCompletionTool = {
  type: "function",
  function: getMailBodyDefinition,
};
