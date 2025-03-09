import mongoose from "mongoose";

import { MessageType } from "../types";

const conn = mongoose.createConnection(
  `mongodb+srv://${process.env.MONGO_USER}:${process.env.MONGO_PASSWORD}@${process.env.MONGO_ADDRESS}/${process.env.MONGO_DATABASE_NAME}?retryWrites=true&w=majority`,
  {
    autoIndex: true,
  },
);

const Schema = mongoose.Schema;

const ProfileSchema = new Schema({
  name: { type: String, required: true },
  email: { type: String, required: true },
  current_character_id: { type: Schema.Types.ObjectId, required: true },
});

export const ProfileModel = conn.model("profiles", ProfileSchema);

const UserContextSchema = new Schema({
  profile_id: {
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
    profile_id: {
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
      enum: [MessageType.UserMessage, MessageType.AiMessage],
    },
  },
  {
    timeseries: {
      metaField: "profile_id",
      timeField: "timestamp",
      granularity: "seconds",
    },
    versionKey: false,
  },
);

MessageSchema.index(
  { timestamp: 1 },
  { expireAfterSeconds: 60 * 60 * 24, partialFilterExpression: { profile_id: { $exists: true } } },
);

export const MessageModel = conn.model("penk_messages", MessageSchema);
