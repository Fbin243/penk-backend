import "./bootstrap";

import { ApolloServer } from "apollo-server";

import { ResolverContext, schema } from "./services/graphql";
import { getProfileByEmail } from "./utils/database/utils";
import { decodeFirebaseJwt } from "./utils/firebase";

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
  console.log(`🚀 Server ready at ${url}`);
});
