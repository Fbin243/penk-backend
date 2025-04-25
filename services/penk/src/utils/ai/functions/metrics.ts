import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { MetricModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

const getMetricsParams = z.object({
  profileId: z.string().describe(SharedDescription.profileId),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z
    .union([z.string(), z.null()])
    .describe(
      "Filter metrics by name using case-insensitive pattern matching (e.g., 'weight' would match 'Body Weight')",
    ),
});

export const functionGetMetrics = async (props: {
  profileId: string;
  categoryId?: string | null;
  name?: string;
}) => {
  console.log(`[Tool: ${Tool.GetMetrics}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  const query: Record<string, unknown> = { character_id: props.profileId };
  if (props.categoryId) {
    query.category_id = props.categoryId;
  }
  if (props.name) {
    query.name = { $regex: props.name, $options: "i" };
  }

  const metrics = await MetricModel.find(query);

  return metrics;
};

export const toolGetMetrics = zodFunction({
  name: Tool.GetMetrics,
  description:
    "Retrieve metrics for a user with optional filtering by category and name. Call this tool if user want to see stats or metrics.",
  parameters: getMetricsParams,
});

const createMetricParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z.string().describe("Metric name"),
  value: z.number().describe("Metric value"),
  unit: z.string().describe("Metric unit"),
});

export const functionCreateMetric = async (props: {
  firebaseUID: string;
  categoryId?: string;
  name: string;
  value: number;
  unit: string;
}) => {
  console.log(`[Tool: ${Tool.CreateMetric}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertMetric(
      {
        name: props.name,
        value: props.value,
        unit: props.unit,
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

export const toolCreateMetric = zodFunction({
  name: Tool.CreateMetric,
  description: `Create a new metric. Before using this tool, always check for existing metrics with GetMetrics first to avoid creating duplicates. If a similar metric exists, use ${Tool.UpdateMetric} instead.`,
  parameters: createMetricParams,
});

const updateMetricParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Metric ID"),
  name: z.string().describe("Metric name"),
  value: z.number().describe("Metric value"),
  unit: z.string().describe("Metric unit"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
});

export const functionUpdateMetric = async (props: {
  firebaseUID: string;
  id: string;
  name: string;
  value: number;
  unit: string;
  categoryId?: string;
}) => {
  console.log(`[Tool: ${Tool.UpdateMetric}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertMetric(
      {
        id: props.id,
        name: props.name,
        value: props.value,
        unit: props.unit,
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

export const toolUpdateMetric = zodFunction({
  name: Tool.UpdateMetric,
  description: "Update an existing metric",
  parameters: updateMetricParams,
});

const deleteMetricParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Metric ID"),
});

export const functionDeleteMetric = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteMetric}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteMetric(
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

export const toolDeleteMetric = zodFunction({
  name: Tool.DeleteMetric,
  description: "Delete an existing metric",
  parameters: deleteMetricParams,
});
