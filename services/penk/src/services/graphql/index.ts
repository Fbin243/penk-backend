import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

export interface ResolverContext {
  token: string;
  email: string;
  profileId: string;
}

const typeDefs = gql(readFileSync(resolve(__dirname, "schema.graphql"), "utf8"));

const resolvers = {
  Query: {
    messages: async (_, args, context) => {
      throw new Error("Not implemented");
    },
  },
  Mutation: {
    chat: async (_, args, context) => {
      throw new Error("Not implemented");
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
