import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { HabitModel } from "../../database/mongo";
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
  description: "Retrieve habits for a user with optional filtering by category and name.",
  parameters: getHabitsParams,
});
