import { zodFunction } from "openai/helpers/zod";
import z from "zod";

import { GoalModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { isoStringToUnixSeconds } from "../../time";
import { Tool } from "../../types";
import { CheckboxInput, SharedDescription } from "./shared";

const GoalMetricInput = z.object({
  id: z.string().describe("Metric ID"),
  condition: z
    .enum(["eq", "gt", "gte", "ir", "lt", "lte"])
    .describe(
      "Condition for the metric (eq: equal, gt: greater than, gte: greater than or equal, ir: in range, lt: less than, lte: less than or equal)",
    ),
  targetValue: z
    .union([z.number(), z.null()])
    .describe("Target value for the metric. Use null for ir condition"),
  rangeValue: z
    .union([
      z.object({
        min: z.number().describe("Minimum value for the range"),
        max: z.number().describe("Maximum value for the range"),
      }),
      z.null(),
    ])
    .describe("Range value for the metric. Use null when the condition is not ir"),
});

const GoalInput = z
  .object({
    id: z
      .union([z.string(), z.null()])
      .describe("Null when creating a new goal, otherwise the goal ID"),
    name: z.string().describe("Goal name"),
    description: z.string().describe("Goal description"),
    startTime: z.string().describe(SharedDescription.datetime),
    endTime: z.string().describe(SharedDescription.datetime),
    metrics: z.array(GoalMetricInput).describe("Metrics for the goal"),
    checkboxes: z.array(CheckboxInput).describe("Checkboxes for the goal"),
  })
  .describe(
    "Goal input parameters. Goal input must have at least one metric or checkbox. Default start time is now, end time is 30 days from now",
  );

const getGoalsParams = z.object({
  profileId: z.string().describe("Profile ID"),
  categoryId: z.union([z.string(), z.null()]).describe("Assigned category ID").optional(),
  name: z
    .string()
    .describe("Filter goals by name using case-insensitive pattern matching")
    .optional(),
});

export const functionGetGoals = async (props: {
  profileId: string;
  categoryId?: string | null;
  name?: string;
}) => {
  console.log(`[Tool: ${Tool.GetGoals}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  const query: Record<string, unknown> = { character_id: props.profileId };
  if (props.categoryId) {
    query.category_id = props.categoryId;
  }
  if (props.name) {
    query.name = { $regex: props.name, $options: "i" };
  }

  const goals = await GoalModel.find(query);

  return goals;
};

export const toolGetGoals = zodFunction({
  name: Tool.GetGoals,
  description:
    "Retrieve goals for a user with optional filtering by category and name. Always use this tool before creating new goals to check if similar goals already exist.",
  parameters: getGoalsParams,
});

const createGoalParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: GoalInput,
});

export const functionCreateGoal = async (props: {
  firebaseUID: string;
  input: z.infer<typeof GoalInput>;
}) => {
  console.log(`[Tool: ${Tool.CreateGoal}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertGoal(
      {
        name: props.input.name,
        description: props.input.description,
        startTime: isoStringToUnixSeconds(props.input.startTime),
        endTime: isoStringToUnixSeconds(props.input.endTime),
        metrics: props.input.metrics?.map((metric) => ({
          condition: metric.condition,
          targetValue: metric.targetValue || undefined,
          rangeValue: metric.rangeValue || undefined,
        })),
        checkboxes: props.input.checkboxes?.map((checkbox) => ({
          id: checkbox.id || undefined,
          name: checkbox.name,
          completed: checkbox.completed,
        })),
      },
      createMetadata(props.firebaseUID),
      (err, res) => {
        if (err) {
          reject(err);
          return;
        }
        resolve(res);
      },
    );
  });
};

export const toolCreateGoal = zodFunction({
  name: Tool.CreateGoal,
  description: `Create a new goal. Before using this tool, always check for existing goals with ${Tool.GetGoals} first to avoid creating duplicates. If a similar goal exists, use ${Tool.UpdateGoal} instead.`,
  parameters: createGoalParams,
});

const updateGoalParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: GoalInput,
});

export const functionUpdateGoal = async (props: {
  firebaseUID: string;
  input: z.infer<typeof GoalInput>;
}) => {
  console.log(`[Tool: ${Tool.UpdateGoal}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertGoal(
      {
        id: props.input.id || undefined,
        name: props.input.name,
        description: props.input.description,
        startTime: isoStringToUnixSeconds(props.input.startTime),
        endTime: isoStringToUnixSeconds(props.input.endTime),
        metrics: props.input.metrics?.map((metric) => ({
          condition: metric.condition,
          targetValue: metric.targetValue || undefined,
          rangeValue: metric.rangeValue || undefined,
        })),
        checkboxes: props.input.checkboxes?.map((checkbox) => ({
          id: checkbox.id || undefined,
          name: checkbox.name,
          completed: checkbox.completed,
        })),
      },
      createMetadata(props.firebaseUID),
      (err, res) => {
        if (err) {
          reject(err);
          return;
        }
        resolve(res);
      },
    );
  });
};

export const toolUpdateGoal = zodFunction({
  name: Tool.UpdateGoal,
  description: `Update an existing goal. Before using this tool, always check for existing goals with ${Tool.GetGoals} first to avoid creating duplicates. If a similar goal exists, use ${Tool.UpdateGoal} instead.`,
  parameters: updateGoalParams,
});

const deleteGoalParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Goal ID"),
});

export const functionDeleteGoal = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteGoal}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteGoal({ id: props.id }, createMetadata(props.firebaseUID), (err, res) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(res);
    });
  });
};

export const toolDeleteGoal = zodFunction({
  name: Tool.DeleteGoal,
  description: `Delete an existing goal. Before using this tool, always check for existing goals with ${Tool.GetGoals} first to avoid creating duplicates. If a similar goal exists, use ${Tool.UpdateGoal} instead.`,
  parameters: deleteGoalParams,
});
