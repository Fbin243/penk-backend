import dotenv from "dotenv";
import express from "express";

dotenv.config({
  path: ".env.penk",
});

const app = express();
const port = 3011;

app.get("/", (req, res) => {
  res.send("Hello PenK!");
});

app.listen(port, () => {
  console.log(`PenK Service listening on port ${port}`);
});
