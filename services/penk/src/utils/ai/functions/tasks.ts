import { Metadata } from "@grpc/grpc-js";
import { zodFunction } from "openai/helpers/zod";
import { z } from "zod";

import { coreClient } from "../../grpc";
import { Tool } from "../../types";
import { SharedDescription } from "./shared";

const createTaskParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  name: z.string().describe("Task name"),
  categoryId: z.string().nullable().describe(SharedDescription.assignedCategoryId),
  priority: z.number().describe(SharedDescription.eisenHowerMatrix),
  deadline: z.string().nullable().describe(SharedDescription.datetime),
  subtasks: z
    .array(
      z.object({
        name: z.string().describe("Subtask name"),
        value: z.boolean().describe("Subtask completion status"),
      }),
    )
    .nullable()
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

  const metadata = new Metadata();
  metadata.add("service-name", "penk");
  metadata.add("x-user-id", props.firebaseUID);

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        name: props.name,
        categoryId: props.categoryId !== "" ? props.categoryId : undefined,
        priority: props.priority,
        deadline: props.deadline
          ? Math.floor(new Date(props.deadline).getTime() / 1000)
          : undefined,
        subtasks: props.subtasks,
      },
      metadata,
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
  description: "Create a new task",
  parameters: createTaskParams,
});

const updateTaskParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Task ID"),
  name: z.string().describe("Task name"),
  categoryId: z.string().nullable().describe(SharedDescription.assignedCategoryId),
  priority: z.number().describe(SharedDescription.eisenHowerMatrix),
  completedTime: z.string().nullable().describe(SharedDescription.datetime),
  deadline: z.string().nullable().describe(SharedDescription.datetime),
  subtasks: z
    .array(
      z.object({
        id: z.string().nullable().describe("Subtask ID"),
        name: z.string().describe("Subtask name"),
        value: z.boolean().describe("Subtask completion status"),
      }),
    )
    .nullable()
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

  const metadata = new Metadata();
  metadata.add("service-name", "penk");
  metadata.add("x-user-id", props.firebaseUID);

  return new Promise((resolve, reject) => {
    coreClient.UpsertTask(
      {
        id: props.id,
        name: props.name,
        categoryId: props.categoryId !== "" ? props.categoryId : undefined,
        priority: props.priority,
        completedTime: props.completedTime
          ? Math.floor(new Date(props.completedTime).getTime() / 1000)
          : undefined,
        deadline: props.deadline
          ? Math.floor(new Date(props.deadline).getTime() / 1000)
          : undefined,
        subtasks: props.subtasks.map((st) => ({
          id: st.id !== "" ? st.id : undefined,
          name: st.name,
          value: st.value,
        })),
      },
      metadata,
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
  description: "Update an existing task",
  parameters: updateTaskParams,
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

  const metadata = new Metadata();
  metadata.add("service-name", "penk");
  metadata.add("x-user-id", props.firebaseUID);

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        taskId: props.taskId,
        startTime: Math.floor(new Date(props.startTime).getTime() / 1000),
        endTime: Math.floor(new Date(props.endTime).getTime() / 1000),
      },
      metadata,
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
  description: "Create a new task session",
  parameters: createTaskSessionParams,
});

const updateTaskSessionParams = z.object({
  firebaseUID: z.string().describe(SharedDescription.firebaseUID),
  id: z.string().describe("Task session ID"),
  taskId: z.string().describe("Task ID"),
  startTime: z.string().describe(SharedDescription.datetime),
  endTime: z.string().describe(SharedDescription.datetime),
  completedTime: z.string().nullable().describe(SharedDescription.datetime),
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

  const metadata = new Metadata();
  metadata.add("service-name", "penk");
  metadata.add("x-user-id", props.firebaseUID);

  return new Promise((resolve, reject) => {
    coreClient.UpsertTaskSession(
      {
        id: props.id,
        taskId: props.taskId,
        startTime: Math.floor(new Date(props.startTime).getTime() / 1000),
        endTime: Math.floor(new Date(props.endTime).getTime() / 1000),
        completedTime: props.completedTime
          ? Math.floor(new Date(props.completedTime).getTime() / 1000)
          : undefined,
      },
      metadata,
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
  description: "Update an existing task session",
  parameters: updateTaskSessionParams,
});
