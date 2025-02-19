import dotenv from "dotenv";
import express from "express";
import { createHandler } from "graphql-http/lib/use/express";

import { schema } from "./schema";

dotenv.config({
  path: ".env.penk",
});

const app = express();
const port = 3011;

app.all("/graphql", createHandler({ schema }));

app.get("/", (_, res) => {
  res.send("Hello PenK!");
});

app.listen(port, () => {
  console.log(`PenK Service listening on port ${port}`);
});
