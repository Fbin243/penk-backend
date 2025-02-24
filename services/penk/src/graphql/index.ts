import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { MessageType, Resolvers } from "../__generated__/types";
import { MessageModel } from "../db/mongo";
import { getMessages, getUserContext } from "../db/utils";
import { chat } from "../openai";

const typeDefs = gql(
  readFileSync(resolve(__dirname, "schema.graphql"), "utf8"),
);

const getTempProfileByTokenId = (token: string) => {
  const tempProfileId = "6735a19cc0e37098e0286d6b";
  return tempProfileId;
};

const resolvers: Resolvers = {
  Query: {
    messages: async (_, __, context) => {
      const messages = await getMessages(
        getTempProfileByTokenId(context.token),
      );
      return messages.map((m) => ({
        content: m.content,
        timestamp: m.timestamp.toISOString(),
        type: m.type,
      }));
    },
    userContext: async (_, __, context) => {
      return await getUserContext(getTempProfileByTokenId(context.token));
    },
  },
  Mutation: {
    chat: async (_, args, context) => {
      const tempProfileId = getTempProfileByTokenId(context.token);

      const [userContext, messages] = await Promise.all([
        getUserContext(tempProfileId),
        getMessages(tempProfileId),
      ]);

      const botMessage = await chat({
        userContext,
        userData: {
          id: tempProfileId,
          name: "Test User",
          currentCharacterId: "test-char-1",
          characters: [
            {
              id: "test-char-1",
              name: "Test Char 1",
              categories: [
                {
                  id: "ct1",
                  name: "books",
                },
                {
                  id: "ct2",
                  name: "music",
                },
              ],
            },
            {
              id: "test-char-2",
              name: "Test Char 2",
              categories: [
                {
                  id: "ct3",
                  name: "coding",
                },
              ],
            },
          ],
        },
        history: messages.map((m) => ({
          role: m.type === MessageType.UserMessage ? "user" : "assistant",
          content: m.content,
        })),
        content: args.content,
      });

      await MessageModel.insertMany([
        {
          profile_id: tempProfileId,
          type: MessageType.UserMessage,
          content: args.content,
          timestamp: new Date(),
        },
        {
          profile_id: tempProfileId,
          type: MessageType.AiMessage,
          content: botMessage.content,
          timestamp: new Date(botMessage.timestamp),
        },
      ]);

      return botMessage;
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
