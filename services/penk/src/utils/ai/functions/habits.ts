import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { HabitModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

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
  name: z.string().describe("Habit name"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  value: z
    .number()
    .describe(
      "Habit value. If completion type is Time, the value is in seconds, e.g. value = 300 means 5 minutes",
    ),
  unit: z.string().describe("Habit unit. If completion type is Time, unit should be empty."),
  completionType: z.enum(["Number", "Time"]).describe("Completion type"),
  rrule: z.string().describe("RRule for habit. Default is RRULE:FREQ=DAILY;INTERVAL=1"),
  reset: z.enum(["Daily", "Weekly", "Monthly"]).describe("Reset frequency"),
});

export const functionCreateHabit = async (props: {
  firebaseUID: string;
  name: string;
  categoryId?: string;
  value: number;
  unit: string;
  completionType: "Number" | "Time";
  rrule: string;
  reset: "Daily" | "Weekly" | "Monthly";
}) => {
  console.log(`[Tool: ${Tool.CreateHabit}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertHabit(
      {
        name: props.name,
        value: props.value,
        unit: props.unit,
        completionType: props.completionType,
        rrule: props.rrule,
        reset: props.reset,
        categoryId: props.categoryId,
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
  id: z.string().describe("Habit ID"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z.string().describe("Habit name"),
  value: z
    .number()
    .describe(
      "Habit value. If completion type is Time, the value is in seconds, e.g. value = 300 means 5 minutes",
    ),
  unit: z.string().describe("Habit unit. If completion type is Time, unit should be empty."),
  completionType: z.enum(["Number", "Time"]).describe("Completion type"),
  rrule: z.string().describe("RRule for habit. Default is RRULE:FREQ=DAILY;INTERVAL=1"),
  reset: z.enum(["Daily", "Weekly", "Monthly"]).describe("Reset frequency"),
});

export const functionUpdateHabit = async (props: {
  firebaseUID: string;
  id: string;
  categoryId?: string;
  name: string;
  value: number;
  unit: string;
  completionType: "Number" | "Time";
  rrule: string;
  reset: "Daily" | "Weekly" | "Monthly";
}) => {
  console.log(`[Tool: ${Tool.UpdateHabit}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertHabit(
      {
        id: props.id,
        name: props.name,
        value: props.value,
        unit: props.unit,
        completionType: props.completionType,
        rrule: props.rrule,
        reset: props.reset,
        categoryId: props.categoryId,
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
