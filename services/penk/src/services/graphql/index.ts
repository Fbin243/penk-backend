import { buildSubgraphSchema } from "@apollo/subgraph";
import { ApolloServer, gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { OAuthTokenModel, PenKContextModel } from "../../utils/database/mongo";
import { getLinkedAccounts, getPenKMessages, getProfileByEmail } from "../../utils/database/utils";
import { decodeFirebaseJwt } from "../../utils/firebase";
import { getGoogleAuthUrl } from "../../utils/googleapis";
import { Resolvers } from "../../utils/types";

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
    googleAuthUrl: async (_, args, context) => {
      requireAuth(context);
      const url = await getGoogleAuthUrl(context.profileId, args.type);
      return url;
    },
    linkedAccounts: async (_, __, context) => {
      requireAuth(context);
      const linkedAccounts = await getLinkedAccounts(context.profileId);
      return linkedAccounts;
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
    revokeLinkedAccount: async (_, args, context) => {
      requireAuth(context);
      await OAuthTokenModel.deleteOne({ _id: args.id });
      return true;
    },
  },
};

const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});

export const startGraphQLServer = async () => {
  const server = new ApolloServer({
    cors: {
      allowedHeaders: "Authorization",
    },
    context: async ({ req }) => {
      const resolverContext: ResolverContext = {
        token: "",
        email: "",
        userId: "",
        profileId: "",
      };

      if (req.headers.authorization) {
        const token = `${req.headers.authorization}`.split(" ")[1];
        if (token) {
          const decodedToken = await decodeFirebaseJwt(token);
          if (!decodedToken?.email) throw new Error("invalid jwt");
          const mongoProfile = await getProfileByEmail(decodedToken.email);
          if (!mongoProfile) throw new Error("profile not found");

          resolverContext.token = token;
          resolverContext.email = decodedToken.email;
          resolverContext.userId = mongoProfile._id.toString();
          resolverContext.profileId = mongoProfile.current_character_id.toString();
        }
      }

      return resolverContext;
    },
    schema,
  });

  server.listen(8099).then(({ url }) => {
    console.log(`🚀 GraphQL Server ready at ${url}`);
  });
};
