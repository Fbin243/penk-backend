import { ApolloServer } from "apollo-server";
import dotenv from "dotenv";

import { schema } from "./graphql";

dotenv.config({
  path: ".env.penk",
});

const server = new ApolloServer({ schema });

server.listen(8099).then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});
