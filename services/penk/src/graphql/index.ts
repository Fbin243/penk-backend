import { buildSubgraphSchema } from "@apollo/subgraph";
import { gql } from "apollo-server";
import { readFileSync } from "fs";
import { resolve } from "path";

import { UserContext } from "../database/mongodb";
import { getRedisClient } from "../database/redis";
import { Resolvers } from "./__generated__/resolvers-types";

const typeDefs = gql(
  readFileSync(resolve(__dirname, "schema.graphql"), "utf8"),
);

const resolvers: Resolvers = {
  Query: {
    helloPenK: async () => {
      const redisClient = await getRedisClient();
      await redisClient.set("mykey", "myvalue", { EX: 300 });
      return "PenK";
    },
    userContext: async () => {
      const tempProfileId = "6735a19cc0e37098e0286d6b";

      let userContext = await UserContext.findOne({
        profileId: tempProfileId,
      });

      if (!userContext) {
        userContext = await UserContext.create({
          profileId: tempProfileId,
          locale: "vi",
          timezone: "Asia/Ho_Chi_Minh",
          context: "",
          preferences: {
            tone: "funny",
          },
        });
      }

      return {
        id: userContext._id.toString(),
        profileID: userContext.profileId.toString(),
        timezone: userContext.timezone,
        locale: userContext.locale,
        context: userContext.context || "",
        preferences: {
          tone: userContext.preferences?.tone || "",
        },
      };
    },
  },
};

export const schema = buildSubgraphSchema({
  typeDefs,
  resolvers,
});
