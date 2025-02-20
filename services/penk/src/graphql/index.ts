import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

const typeDefs = gql(
  readFileSync(resolve(__dirname, "schema.graphql"), "utf8"),
);

const resolvers = {
  Query: {
    helloPenK: () => "PenK",
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
