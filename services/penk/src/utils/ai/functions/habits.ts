import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { HabitModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

const HabitInput = z.object({
  id: z
    .union([z.string(), z.null()])
    .describe("Null when creating a new habit, otherwise the habit ID"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z.string().describe("Habit name"),
  value: z
    .number()
    .describe(
      "Habit value. IMPORTANT: For Time habits, value MUST be in seconds (not minutes). Example: 5 minutes = 300 seconds, 1 hour = 3600 seconds. For Number habits, value is the numeric quantity.",
    ),
  unit: z
    .string()
    .describe(
      "Habit unit. IMPORTANT: For Time habits, unit MUST be empty string. For Number habits, specify appropriate unit (e.g., 'pages', 'glasses').",
    ),
  completionType: z.enum(["Number", "Time"]).describe("Completion type"),
  rrule: z.string().describe("RRule for habit. Default is RRULE:FREQ=DAILY;INTERVAL=1"),
  resetDuration: z
    .enum(["Daily", "Weekly", "Monthly"])
    .describe(
      "Reset duration. Default is Daily. Some examples: Drink 2L of water every day, Read 100 pages every week, etc.",
    ),
});

const getHabitsParams = z.object({
  profileId: z.string().describe(SharedDescription.profileId),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z
    .union([z.string(), z.null()])
    .describe("Filter habits by name using case-insensitive pattern matching"),
});

export const functionGetHabits = async (props: {
  profileId: string;
  categoryId?: string | null;
  name?: string;
}) => {
  console.log(`[Tool: ${Tool.GetHabits}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  const query: Record<string, unknown> = { character_id: props.profileId };
  if (props.categoryId) {
    query.category_id = props.categoryId;
  }
  if (props.name) {
    query.name = { $regex: props.name, $options: "i" };
  }

  const habits = await HabitModel.find(query);

  return habits;
};

export const toolGetHabits = zodFunction({
  name: Tool.GetHabits,
  description:
    "Retrieve habits for a user with optional filtering by category and name. Always use this tool before creating new habits to check if similar habits already exist.",
  parameters: getHabitsParams,
});

const createHabitParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: HabitInput,
});

export const functionCreateHabit = async (props: {
  firebaseUID: string;
  input: z.infer<typeof HabitInput>;
}) => {
  console.log(`[Tool: ${Tool.CreateHabit}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertHabit(
      {
        categoryId: props.input.categoryId || undefined,
        name: props.input.name,
        value: props.input.value,
        unit: props.input.unit,
        completionType: props.input.completionType,
        rrule: props.input.rrule,
        resetDuration: props.input.resetDuration,
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

export const toolCreateHabit = zodFunction({
  name: Tool.CreateHabit,
  description: `Create a new habit. Before using this tool, always check for existing habits with GetHabits first to avoid creating duplicates. If a similar habit exists, use ${Tool.UpdateHabit} instead.`,
  parameters: createHabitParams,
});

const updateHabitParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: HabitInput,
});

export const functionUpdateHabit = async (props: {
  firebaseUID: string;
  input: z.infer<typeof HabitInput>;
}) => {
  console.log(`[Tool: ${Tool.UpdateHabit}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertHabit(
      {
        id: props.input.id || undefined,
        name: props.input.name,
        value: props.input.value,
        unit: props.input.unit,
        completionType: props.input.completionType,
        rrule: props.input.rrule,
        resetDuration: props.input.resetDuration,
        categoryId: props.input.categoryId || undefined,
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

export const toolUpdateHabit = zodFunction({
  name: Tool.UpdateHabit,
  description: "Update an existing habit",
  parameters: updateHabitParams,
});

const deleteHabitParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Habit ID"),
});

export const functionDeleteHabit = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteHabit}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteHabit(
      {
        id: props.id,
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

export const toolDeleteHabit = zodFunction({
  name: Tool.DeleteHabit,
  description: "Delete an existing habit",
  parameters: deleteHabitParams,
});
