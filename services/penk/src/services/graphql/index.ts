import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { Message, MessageType, Resolvers } from "../../utils/types";
import { chat } from "../ai/utils";
import { PenKContextModel, PenKMessageModel } from "../database/mongo";
import { getPenKData, getPenKMessages } from "../database/utils";

export interface ResolverContext {
  token: string;
  email: string;
  userId: string;
  profileId: string;
}

const typeDefs = gql(readFileSync(resolve(__dirname, "schema.graphql"), "utf8"));

const requireAuth = (context: ResolverContext) => {
  if (!context.token) throw new Error("unauthorized");
};

const resolvers: Resolvers = {
  Query: {
    context: async (_, __, context) => {
      requireAuth(context);
      const penkContext = await PenKContextModel.findOne({ user_id: context.userId });
      return penkContext;
    },
    messages: async (_, __, context) => {
      requireAuth(context);
      const messages = await getPenKMessages(context.profileId);
      return messages.map((message) => ({
        type: message.type,
        content: message.content,
        timestamp: message.timestamp.toISOString(),
      }));
    },
  },
  Mutation: {
    upsertContext: async (_, args, context) => {
      requireAuth(context);
      const penkContext = await PenKContextModel.findOneAndUpdate(
        { profile_id: context.profileId },
        args.input,
        { upsert: true, new: true },
      );
      return penkContext;
    },
    chat: async (_, args, context) => {
      requireAuth(context);

      const { content } = args;

      const penkData = await getPenKData(context.userId);

      const userMessage: Message = {
        type: MessageType.UserMessage,
        content,
        timestamp: new Date().toISOString(),
      };

      const previousMessages = (await getPenKMessages(context.profileId)).map((message) => ({
        type: message.type,
        content: message.content,
        timestamp: message.timestamp.toISOString(),
      }));

      const aiMessage = await chat({
        userData: JSON.stringify(penkData),
        messages: [...previousMessages, userMessage],
      });

      await PenKMessageModel.insertMany([
        { ...userMessage, profile_id: context.profileId },
        { ...aiMessage, profile_id: context.profileId },
      ]);

      return aiMessage;
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
