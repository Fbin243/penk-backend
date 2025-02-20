import { ApolloServer } from "apollo-server";
import dotenv from "dotenv";

import { schema } from "./graphql";
import { ResolverContext } from "./graphql/types";

dotenv.config({
  path: ".env.penk",
});

const server = new ApolloServer({
  cors: {
    allowedHeaders: "Authorization",
  },
  context: ({ req }) => {
    const token = req.headers.Authorization
      ? `${req.headers.Authorization}`.split(" ")[1]
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
