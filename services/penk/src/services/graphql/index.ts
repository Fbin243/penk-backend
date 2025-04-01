import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { PenKContextModel, PenKMessageModel } from "../../utils/database/mongo";
import { getPenKData, getPenKMessages } from "../../utils/database/utils";
import { Message, MessageType, Resolvers } from "../../utils/types";
import { base64ToUploadable, chat, transcribeAudio } from "../ai/utils";

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

      const { content, voiceMode } = args.input;

      const [transcribedInput, penkData, penkMessages] = await Promise.all([
        voiceMode
          ? transcribeAudio(base64ToUploadable(content, "audio.wav"))
          : { text: content, cost: 0 },
        getPenKData(context.userId),
        getPenKMessages(context.profileId),
      ]);

      const userMessage: Message = {
        type: MessageType.UserMessage,
        content: transcribedInput.text,
        timestamp: new Date().toISOString(),
      };

      const chatResult = await chat({
        userData: JSON.stringify(penkData),
        history: penkMessages.map((message) => ({
          type: message.type,
          content: message.content,
          timestamp: message.timestamp.toISOString(),
        })),
        newMessage: transcribedInput.text,
        voiceMode: voiceMode ?? false,
      });

      // Why is this not async?
      // -> Because I want to return the response ASAP
      // and it is not a big deal even if it fails to insert, I guess xD
      PenKMessageModel.insertMany([
        { ...userMessage, profile_id: context.profileId },
        { ...chatResult.aiMessage, profile_id: context.profileId },
      ]);

      return {
        message: chatResult.aiMessage,
        audio: chatResult.audio?.data,
      };
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
