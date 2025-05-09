import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { coreClient, createMetadata } from "../../grpc";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

const createCategoryParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  name: z.string().describe("Category name, e.g. 'Work'"),
  description: z.string().describe("Category description"),
  emoji: z.string().describe("Category emoji, e.g. '💼'"),
  color: z
    .string()
    .describe("Hex color, e.g. '#53D8CA'. The color should be inferred from the category name."),
});

export const functionCreateCategory = async (props: {
  firebaseUID: string;
  name: string;
  description: string;
  emoji: string;
  color: string;
}) => {
  console.log(`[Tool: ${Tool.CreateCategory}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertCategory(
      {
        name: props.name,
        description: props.description,
        style: {
          icon: props.emoji,
          color: props.color,
        },
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

export const toolCreateCategory = zodFunction({
  name: Tool.CreateCategory,
  description:
    "Create a new organizational category that can be applied across tasks, habits, and metrics. Categories help group related items and provide visual consistency with custom emoji icons and colors. Always check for similar existing categories before creating new ones.",
  parameters: createCategoryParams,
});

const updateCategoryParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Category ID"),
  name: z.string().describe("Category name, e.g. 'Work'"),
  description: z.string().describe("Category description"),
  emoji: z.string().describe("Category emoji, e.g. '💼'"),
  color: z
    .string()
    .describe("Hex color, e.g. '#53D8CA'. The color should be inferred from the category name."),
});

export const functionUpdateCategory = async (props: {
  firebaseUID: string;
  id: string;
  name: string;
  description: string;
  emoji: string;
  color: string;
}) => {
  console.log(`[Tool: ${Tool.UpdateCategory}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertCategory(
      {
        id: props.id,
        name: props.name,
        description: props.description,
        style: {
          icon: props.emoji,
          color: props.color,
        },
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

export const toolUpdateCategory = zodFunction({
  name: Tool.UpdateCategory,
  description:
    "Update an existing category's name, description, emoji icon, or color. Changes to categories will affect all associated tasks, habits, and metrics. Use this to maintain consistent organizational structure across the system.",
  parameters: updateCategoryParams,
});

const deleteCategoryParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Category ID"),
});

export const functionDeleteCategory = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteCategory}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteCategory(
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

export const toolDeleteCategory = zodFunction({
  name: Tool.DeleteCategory,
  description:
    "Delete an existing category. This will remove the category from the system but will not delete any associated tasks, habits, or metrics. Those items will become unassigned. Verify the category ID exists before attempting deletion.",
  parameters: deleteCategoryParams,
});
