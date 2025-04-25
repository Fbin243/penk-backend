import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { TaskModel } from "../../database/mongo";
import { coreClient, createMetadata } from "../../grpc";
import { isoStringToUnixSeconds } from "../../time";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

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
  name: z.string().describe("Task name"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  priority: z.number().describe(SharedDescription.eisenHowerMatrix),
  deadline: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
  subtasks: z
    .array(
      z.object({
        name: z.string().describe("Subtask name"),
        value: z.boolean().describe("Subtask completion status"),
      }),
    )
    .describe("Subtasks"),
});

export const functionCreateTask = async (props: {
  firebaseUID: string;
  name: string;
  priority: number;
  categoryId?: string;
  deadline: string;
  subtasks: { name: string; value: boolean }[];
}) => {
  console.log(`[Tool: ${Tool.CreateTask}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        name: props.name,
        categoryId: props.categoryId,
        priority: props.priority,
        deadline: isoStringToUnixSeconds(props.deadline),
        subtasks: props.subtasks,
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
  id: z.string().describe("Task ID"),
  name: z.string().describe("Task name"),
  categoryId: z.union([z.string(), z.null()]).describe(SharedDescription.assignedCategoryId),
  priority: z.number().describe(SharedDescription.eisenHowerMatrix),
  completedTime: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
  deadline: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
  subtasks: z
    .array(
      z.object({
        id: z.string().describe("Subtask ID"),
        name: z.string().describe("Subtask name"),
        value: z.boolean().describe("Subtask completion status"),
      }),
    )
    .describe("Subtasks"),
});

export const functionUpdateTask = async (props: {
  firebaseUID: string;
  id: string;
  name: string;
  priority: number;
  categoryId?: string;
  completedTime?: string;
  deadline?: string;
  subtasks: { id?: string; name: string; value: boolean }[];
}) => {
  console.log(`[Tool: ${Tool.UpdateTask}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        id: props.id,
        name: props.name,
        categoryId: props.categoryId !== "" ? props.categoryId : undefined,
        priority: props.priority,
        completedTime: isoStringToUnixSeconds(props.completedTime),
        deadline: isoStringToUnixSeconds(props.deadline),
        subtasks: props.subtasks.map((st) => ({
          id: st.id !== "" ? st.id : undefined,
          name: st.name,
          value: st.value,
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

const createTaskSessionParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  taskId: z.string().describe("Task ID"),
  startTime: z.string().describe(SharedDescription.datetime),
  endTime: z.string().describe(SharedDescription.datetime),
});

export const functionCreateTaskSession = async (props: {
  firebaseUID: string;
  taskId: string;
  startTime: string;
  endTime: string;
}) => {
  console.log(`[Tool: ${Tool.CreateTaskSession}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        taskId: props.taskId,
        startTime: isoStringToUnixSeconds(props.startTime),
        endTime: isoStringToUnixSeconds(props.endTime),
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
  id: z.string().describe("Task session ID"),
  taskId: z.string().describe("Task ID"),
  startTime: z.string().describe(SharedDescription.datetime),
  endTime: z.string().describe(SharedDescription.datetime),
  completedTime: z.union([z.string(), z.null()]).describe(SharedDescription.datetime),
});

export const functionUpdateTaskSession = async (props: {
  firebaseUID: string;
  id: string;
  taskId: string;
  startTime: string;
  endTime: string;
  completedTime?: string;
}) => {
  console.log(`[Tool: ${Tool.UpdateTaskSession}]`);
  console.dir(props, { depth: null, colors: true });
  console.log();

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        id: props.id,
        taskId: props.taskId,
        startTime: isoStringToUnixSeconds(props.startTime),
        endTime: isoStringToUnixSeconds(props.endTime),
        completedTime: isoStringToUnixSeconds(props.completedTime),
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
  sessions: z.array(
    z.object({
      taskId: z.string().describe("Task ID"),
      startTime: z.string().describe(SharedDescription.datetime),
      endTime: z.string().describe(SharedDescription.datetime),
    }),
  ),
});

export const functionPlanDay = async (props: {
  firebaseUID: string;
  sessions: { taskId: string; startTime: string; endTime: string }[];
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
