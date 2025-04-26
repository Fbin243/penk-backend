import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { TaskModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { isoStringToUnixSeconds } from "../../time";
import { Tool } from "../../types";
import { CheckboxInput, SharedDescription } from "./shared";

const TaskInput = z.object({
  id: z
    .union([z.string(), z.null()])
    .describe("Null when creating a new task, otherwise the task ID"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  name: z.string().describe("Task name"),
  priority: z.number().describe(SharedDescription.eisenHowerMatrix),
  completedTime: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
  deadline: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
  subtasks: z.array(CheckboxInput).describe("Subtasks"),
});

const getTasksParams = z.object({
  profileId: z.string().describe(SharedDescription.profileId),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  priority: z.union([z.number(), z.null()]).describe(SharedDescription.eisenHowerMatrix),
  completed: z
    .union([z.boolean(), z.null()])
    .describe(
      "Filter for task completion status. When true, returns completed tasks; when false, returns incomplete tasks; when null, returns all tasks. Default is false.",
    ),
});

export const functionGetTasks = async (props: {
  profileId: string;
  categoryId?: string | null;
  priority?: number | null;
  completed?: boolean | null;
}) => {
  console.log(`[Tool: ${Tool.GetTasks}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  const query: Record<string, unknown> = { character_id: props.profileId };
  if (props.categoryId) {
    query.category_id = props.categoryId;
  }
  if (props.priority) {
    query.priority = props.priority;
  }
  if (props.completed !== null) {
    query.completed_time = props.completed ? { $ne: null } : null;
  }

  const tasks = await TaskModel.find(query);

  return tasks;
};

export const toolGetTasks = zodFunction({
  name: Tool.GetTasks,
  description:
    "Retrieves tasks for a user with optional filtering by category, priority level, and completion status. By default, returns only incomplete tasks. Always use this tool before creating new tasks to check if similar tasks already exist.",
  parameters: getTasksParams,
});

const createTaskParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: TaskInput,
});

export const functionCreateTask = async (props: {
  firebaseUID: string;
  input: z.infer<typeof TaskInput>;
}) => {
  console.log(`[Tool: ${Tool.CreateTask}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        name: props.input.name,
        categoryId: props.input.categoryId || undefined,
        priority: props.input.priority,
        completedTime: isoStringToUnixSeconds(props.input.completedTime),
        deadline: isoStringToUnixSeconds(props.input.deadline),
        subtasks: props.input.subtasks.map((st) => ({
          name: st.name,
          completed: st.completed,
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

export const toolCreateTask = zodFunction({
  name: Tool.CreateTask,
  description: `Create a new task. Before using this tool, always check for existing tasks with GetTasks first to avoid creating duplicates. If a similar task exists, use ${Tool.UpdateTask} instead. Tasks require a name, priority level, and can optionally have a deadline and subtasks.`,
  parameters: createTaskParams,
});

const updateTaskParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: TaskInput,
});

export const functionUpdateTask = async (props: {
  firebaseUID: string;
  input: z.infer<typeof TaskInput>;
}) => {
  console.log(`[Tool: ${Tool.UpdateTask}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        id: props.input.id || undefined,
        categoryId: props.input.categoryId || undefined,
        priority: props.input.priority,
        completedTime: isoStringToUnixSeconds(props.input.completedTime),
        deadline: isoStringToUnixSeconds(props.input.deadline),
        subtasks: props.input.subtasks.map((st) => ({
          id: st.id || undefined,
          name: st.name,
          completed: st.completed,
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

export const toolUpdateTask = zodFunction({
  name: Tool.UpdateTask,
  description:
    "Update an existing task. Use this tool to modify task properties or mark tasks as complete by setting the completedTime. Required parameters include task ID, name, and priority.",
  parameters: updateTaskParams,
});

const deleteTaskParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Task ID"),
});

export const functionDeleteTask = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteTask}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteTask(
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

export const toolDeleteTask = zodFunction({
  name: Tool.DeleteTask,
  description: "Delete an existing task",
  parameters: deleteTaskParams,
});

const TaskSessionInput = z.object({
  id: z
    .union([z.string(), z.null()])
    .describe("Null when creating a new task session, otherwise the task session ID"),
  taskId: z.string().describe("Task ID"),
  startTime: z.string().describe(SharedDescription.datetime),
  endTime: z.string().describe(SharedDescription.datetime),
  completedTime: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
});

const createTaskSessionParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: TaskSessionInput,
});

export const functionCreateTaskSession = async (props: {
  firebaseUID: string;
  input: z.infer<typeof TaskSessionInput>;
}) => {
  console.log(`[Tool: ${Tool.CreateTaskSession}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        taskId: props.input.taskId,
        startTime: isoStringToUnixSeconds(props.input.startTime),
        endTime: isoStringToUnixSeconds(props.input.endTime),
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

export const toolCreateTaskSession = zodFunction({
  name: Tool.CreateTaskSession,
  description: `Create a new task session block in the user's daily timeline visualization. Task sessions represent scheduled focused work periods with defined start/end times, appearing as visual blocks in the day planner. Always verify the task ID exists first using ${Tool.GetTasks}.`,
  parameters: createTaskSessionParams,
});

const updateTaskSessionParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  input: TaskSessionInput,
});

export const functionUpdateTaskSession = async (props: {
  firebaseUID: string;
  input: z.infer<typeof TaskSessionInput>;
}) => {
  console.log(`[Tool: ${Tool.UpdateTaskSession}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        id: props.input.id || undefined,
        taskId: props.input.taskId,
        startTime: isoStringToUnixSeconds(props.input.startTime),
        endTime: isoStringToUnixSeconds(props.input.endTime),
        completedTime: isoStringToUnixSeconds(props.input.completedTime),
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

export const toolUpdateTaskSession = zodFunction({
  name: Tool.UpdateTaskSession,
  description:
    "Update an existing task session block in the user's daily timeline. Use this to modify scheduled work periods by adjusting timing or marking sessions as completed. Changes will update the visual timeline representation.",
  parameters: updateTaskSessionParams,
});

const deleteTaskSessionParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Task session ID"),
});

export const functionDeleteTaskSession = async (props: { firebaseUID: string; id: string }) => {
  console.log(`[Tool: ${Tool.DeleteTaskSession}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.DeleteTaskSession(
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

export const toolDeleteTaskSession = zodFunction({
  name: Tool.DeleteTaskSession,
  description: "Delete an existing task session",
  parameters: deleteTaskSessionParams,
});

const planDayParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  sessions: z.array(TaskSessionInput),
});

export const functionPlanDay = async (props: {
  firebaseUID: string;
  sessions: z.infer<typeof TaskSessionInput>[];
}) => {
  console.log(`[Tool: ${Tool.PlanDay}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSessions(
      {
        taskSessionInputs: props.sessions.map((session) => ({
          taskId: session.taskId,
          startTime: isoStringToUnixSeconds(session.startTime),
          endTime: isoStringToUnixSeconds(session.endTime),
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

export const toolPlanDay = zodFunction({
  name: Tool.PlanDay,
  description: `Create a comprehensive visual day plan by scheduling multiple task sessions at once. This efficiently builds the user's daily timeline visualization with focused work blocks based on task priorities and available time. Always verify task IDs exist first using ${Tool.GetTasks}. Use this for creating an optimized daily schedule that balances priorities and time constraints.`,
  parameters: planDayParams,
});
