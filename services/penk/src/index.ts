import { ApolloServer } from "apollo-server";
import dotenv from "dotenv";

import { schema } from "./graphql";
import { ResolverContext } from "./graphql/types";

dotenv.config({
  path: `.env.${process.env.NODE_ENV || "development"}`,
});

const server = new ApolloServer({
  cors: {
    allowedHeaders: "Authorization",
  },
  context: ({ req }) => {
    const token = req.headers.authorization
      ? `${req.headers.authorization}`.split(" ")[1]
      : "";

    return {
      token,
    } satisfies ResolverContext;
  },
  schema,
});

server.listen(8099).then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});
