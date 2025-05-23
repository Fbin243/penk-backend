import mongoose from "mongoose";

import { LinkedAccountType, MessageType } from "../types";

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

const PenKContextSchema = new Schema({
  user_id: {
    type: Schema.Types.ObjectId,
    ref: "profiles",
    required: true,
    unique: true,
  },
  timezone: { type: String, default: "Asia/Saigon" },
  locale: { type: String, default: "vi-VN" },
  context: { type: String, default: "" },
});

export const PenKContextModel = conn.model("penk_contexts", PenKContextSchema);

const PenKMessageSchema = new Schema(
  {
    profile_id: {
      type: Schema.Types.ObjectId,
      ref: "characters",
      required: true,
    },
    timestamp: { type: Date, required: true },
    content: { type: String, required: true },
    type: {
      type: String,
      required: true,
      enum: [MessageType.UserMessage, MessageType.AiMessage, MessageType.ToolCallMessage],
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

PenKMessageSchema.index(
  { timestamp: 1 },
  {
    expireAfterSeconds: 60 * 60 * 24 * 3,
    partialFilterExpression: { profile_id: { $exists: true } },
  },
);

export const PenKMessageModel = conn.model("penk_messages", PenKMessageSchema);

const MembershipSchema = new Schema({
  email: { type: String, required: true },
  monthly_credit: { type: Number, default: 0 },
  persistent_credit: { type: Number, default: 0 },
  period_end: { type: Date, default: null },
});

export const MembershipModel = conn.model("memberships", MembershipSchema);

// 1 profile can link multiple oauth tokens
export const OAuthTokenSchema = new Schema({
  profile_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  email: { type: String, required: true },
  type: {
    type: String,
    required: true,
    enum: [LinkedAccountType.Gmail, LinkedAccountType.GoogleCalendar],
  },
  refresh_token: { type: String, required: true },
});

export const OAuthTokenModel = conn.model("oauth_tokens", OAuthTokenSchema);

// Other services' collections
export const TaskSchema = new Schema({
  character_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  category_id: {
    type: Schema.Types.ObjectId,
    ref: "categories",
    default: null,
  },
  name: { type: String, required: true },
  priority: { type: Number, required: true },
  deadline: { type: Date },
  completed_time: { type: Date },
  subtasks: [
    {
      id: { type: Schema.Types.ObjectId, required: true },
      name: { type: String, required: true, default: "" },
      value: { type: Boolean, required: true, default: false },
    },
  ],
});

export const TaskModel = conn.model("tasks", TaskSchema);

export const MetricSchema = new Schema({
  character_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  category_id: {
    type: Schema.Types.ObjectId,
    ref: "categories",
    default: null,
  },
  name: { type: String, required: true },
  value: { type: Number, required: true },
  unit: { type: String },
});

export const MetricModel = conn.model("metrics", MetricSchema);

export const HabitSchema = new Schema({
  character_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  category_id: {
    type: Schema.Types.ObjectId,
    ref: "categories",
    default: null,
  },
  name: { type: String, required: true },
  value: { type: Number, required: true },
  unit: { type: String },
  completion_type: {
    type: String,
    required: true,
    enum: ["Number", "Time"],
  },
  rrule: { type: String },
  reset: {
    type: String,
    required: true,
    enum: ["Daily", "Weekly", "Monthly"],
  },
});

export const HabitModel = conn.model("habits", HabitSchema);

export const GoalSchema = new Schema({
  character_id: {
    type: Schema.Types.ObjectId,
    ref: "characters",
    required: true,
  },
  name: { type: String, required: true },
  description: { type: String },
  start_time: { type: Date },
  end_time: { type: Date },
  completed_time: { type: Date },
  metrics: [
    {
      id: { type: Schema.Types.ObjectId, required: true },
      condition: {
        type: String,
        required: true,
        enum: ["eq", "gt", "gte", "ir", "lt", "lte"],
      },
      target_value: { type: Number },
      range_value: {
        min: { type: Number },
        max: { type: Number },
      },
    },
  ],
  checkboxes: [
    {
      id: { type: Schema.Types.ObjectId, required: true },
      name: { type: String, required: true },
      value: { type: Boolean, required: true },
    },
  ],
});

export const GoalModel = conn.model("goals", GoalSchema);
