import mongoose from "mongoose";

import { MessageType } from "../__generated__/types";

const conn = mongoose.createConnection(
  process.env.MONGO_CONNECTION_STRING || "",
);

conn.useDb("");

const Schema = mongoose.Schema;

const UserContextSchema = new Schema({
  profileId: {
    type: Schema.Types.ObjectId,
    ref: "profiles",
    required: true,
    unique: true,
  },
  timezone: { type: String, default: "UTC" },
  locale: { type: String, default: "en-US" },
  preferences: {
    tone: String,
  },
  context: String,
});

export const UserContextModel = conn.model("user_contexts", UserContextSchema);

const MessageSchema = new Schema(
  {
    profileId: {
      type: Schema.Types.ObjectId,
      ref: "profiles",
      required: true,
      unique: true,
    },
    timestamp: { type: Date, required: true },
    content: { type: String, required: true },
    type: {
      type: String,
      required: true,
      enum: [
        MessageType.UserMessage,
        MessageType.AiMessage,
        MessageType.AiError,
      ],
    },
  },
  {
    timeseries: {
      metaField: "profileId",
      timeField: "timestamp",
      granularity: "seconds",
    },
    versionKey: false,
  },
);

export const MessageModel = conn.model("penk_messages", MessageSchema);
