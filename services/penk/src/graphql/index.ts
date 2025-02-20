import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { Resolvers } from "./__generated__/resolvers-types";

const typeDefs = gql(
  readFileSync(resolve(__dirname, "schema.graphql"), "utf8"),
);

const resolvers: Resolvers = {
  Query: {
    helloPenK: (_, __, context) => {
      console.log("--> context", context);
      return "PenK";
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
