import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { MetricModel } from "../../database/mongo";
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
