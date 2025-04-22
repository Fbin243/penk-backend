import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";

import { getLinkedAccounts } from "../../database/utils";
import { getMails } from "../../googleapis/gmail";
import { LinkedAccountType } from "../../types";
import { SharedDescription } from "./shared";
import { FunctionName } from "./types";

const getMailsDefinition: FunctionDefinition = {
  name: FunctionName.GetMails,
  description:
    "Retrieves emails from the user's Gmail inbox based on specific search criteria. Use this tool when the user asks to check their email, find specific emails, or get updates on their inbox status. This function requires the user's profile ID and a Gmail-compatible search query. The tool returns relevant emails matching the search parameters and should be used to help users find, summarize, or organize their email correspondence. The search query parameter supports all standard Gmail search operators, allowing for precise filtering of emails.",
  parameters: {
    type: "object",
    properties: {
      profileId: {
        type: "string",
        description: SharedDescription.profileId,
      },
      q: {
        type: "string",
        description: `Gmail search query. When fetching emails, consider these scenarios:
        - If the user asks for the latest emails, use the \`newer_than:7d is:unread\` query.
        - All Unread: For all unread emails, use the \`is:unread\` query.
        - Specific Sender: To find emails from a particular sender, use \`from:sender\`.
        - Specific Date: To find emails from a specific date, use \`after:date\`.
        Infer the user's intent and use the most appropriate query to get the job done.`,
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
