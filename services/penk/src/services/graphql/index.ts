import { buildSubgraphSchema } from "@apollo/subgraph";
import { ApolloServer, gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { MembershipModel, OAuthTokenModel, PenKContextModel } from "../../utils/database/mongo";
import { getRedisClient } from "../../utils/database/redis";
import { getLinkedAccounts, getPenKMessages, getProfileByEmail } from "../../utils/database/utils";
import { decodeFirebaseJwt } from "../../utils/firebase";
import { getGoogleAuthUrl } from "../../utils/googleapis";
import { LinkedAccount, Resolvers } from "../../utils/types";

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
      if (penkContext) {
        return {
          locale: penkContext.locale,
          timezone: penkContext.timezone,
          context: penkContext.context,
        };
      }
      return null;
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
    membership: async (_, __, context) => {
      requireAuth(context);
      const membership = await MembershipModel.findOne({ email: context.email });
      if (membership) {
        return {
          monthlyCredit: membership.monthly_credit,
          persistentCredit: membership.persistent_credit,
          periodEnd: membership.period_end ? membership.period_end.toISOString() : null,
        };
      } else {
        // TODO: Get the free persistent credit from the config/database
        const freePersistentCredit = 20;
        const newMembership = new MembershipModel({
          email: context.email,
          monthly_credit: 0,
          persistent_credit: freePersistentCredit,
          period_end: null,
        });
        await newMembership.save();
        return {
          monthlyCredit: newMembership.monthly_credit,
          persistentCredit: newMembership.persistent_credit,
          periodEnd: newMembership.period_end ? newMembership.period_end.toISOString() : null,
        };
      }
    },
  },
  Mutation: {
    upsertContext: async (_, args, context) => {
      requireAuth(context);
      const penkContext = await PenKContextModel.findOneAndUpdate(
        { user_id: context.userId },
        args.input,
        { upsert: true, new: true },
      );
      return penkContext;
    },
    revokeLinkedAccount: async (_, args, context) => {
      requireAuth(context);

      await OAuthTokenModel.deleteOne({ _id: args.id });

      const redisClient = await getRedisClient();
      const currentCache = await redisClient.get(`linked_accounts:${context.profileId}`);
      if (currentCache) {
        const linkedAccounts = JSON.parse(currentCache) as LinkedAccount[];
        const updatedAccounts = linkedAccounts.filter((account) => account.id !== args.id);
        await redisClient.set(
          `linked_accounts:${context.profileId}`,
          JSON.stringify(updatedAccounts),
          { KEEPTTL: true },
        );
      }

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
