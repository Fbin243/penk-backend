import "./bootstrap";

import { ApolloServer } from "apollo-server";

import { schema } from "./graphql";
import { ResolverContext } from "./graphql/types";

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
