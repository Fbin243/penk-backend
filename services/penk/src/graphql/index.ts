import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { chat } from "../openai";
import { MessageModel } from "../utils/db/mongo";
import { getMessages, getUserContext, getUserData } from "../utils/db/utils";
import { MessageType, Profile, Resolvers } from "../utils/types";

export interface ResolverContext {
  token: string;
  profile?: Profile;
}

const typeDefs = gql(readFileSync(resolve(__dirname, "schema.graphql"), "utf8"));

const resolvers: Resolvers = {
  Query: {
    messages: async (_, __, context) => {
      const messages = await getMessages(context.profile.id);
      return messages.map((m) => ({
        content: m.content,
        timestamp: m.timestamp.toISOString(),
        type: m.type,
      }));
    },
    userContext: async (_, __, context) => {
      return await getUserContext(context.profile.id);
    },
  },
  Mutation: {
    chat: async (_, args, context) => {
      const [userData, messages] = await Promise.all([
        getUserData(context.profile.id),
        getMessages(context.profile.id),
      ]);

      const userMessage = {
        profile_id: context.profile.id,
        type: MessageType.UserMessage,
        content: args.content,
        timestamp: new Date(),
      };

      const botMessage = await chat({
        userData: JSON.stringify(userData),
        history: messages.map((m) => ({
          role: m.type === MessageType.UserMessage ? "user" : "assistant",
          content: m.content,
        })),
        content: userMessage.content,
        jwt: context.token,
      });

      await MessageModel.insertMany([
        userMessage,
        {
          profile_id: context.profile.id,
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
