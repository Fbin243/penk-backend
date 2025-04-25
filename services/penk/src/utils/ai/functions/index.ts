import { ChatCompletionTool } from "openai/resources/index.mjs";

import { Tool } from "../../types";
import { functionGetCalendarEvents, toolGetCalendarEvents } from "./calendar";
import {
  functionCreateCategory,
  functionDeleteCategory,
  functionUpdateCategory,
  toolCreateCategory,
  toolDeleteCategory,
  toolUpdateCategory,
} from "./categories";
import { functionGetMails, toolGetMails } from "./gmail";
import {
  functionCreateGoal,
  functionDeleteGoal,
  functionGetGoals,
  functionUpdateGoal,
  toolCreateGoal,
  toolDeleteGoal,
  toolGetGoals,
  toolUpdateGoal,
} from "./goals";
import {
  functionCreateHabit,
  functionDeleteHabit,
  functionGetHabits,
  functionUpdateHabit,
  toolCreateHabit,
  toolDeleteHabit,
  toolGetHabits,
  toolUpdateHabit,
} from "./habits";
import {
  functionCreateMetric,
  functionDeleteMetric,
  functionGetMetrics,
  functionUpdateMetric,
  toolCreateMetric,
  toolDeleteMetric,
  toolGetMetrics,
  toolUpdateMetric,
} from "./metrics";
import {
  functionCreateTask,
  functionCreateTaskSession,
  functionDeleteTask,
  functionDeleteTaskSession,
  functionGetTasks,
  functionPlanDay,
  functionUpdateTask,
  functionUpdateTaskSession,
  toolCreateTask,
  toolCreateTaskSession,
  toolDeleteTask,
  toolDeleteTaskSession,
  toolGetTasks,
  toolPlanDay,
  toolUpdateTask,
  toolUpdateTaskSession,
} from "./tasks";

export const tools: ChatCompletionTool[] = [
  toolGetCalendarEvents,
  toolGetMails,
  toolGetTasks,
  toolCreateCategory,
  toolUpdateCategory,
  toolDeleteCategory,
  toolCreateTask,
  toolUpdateTask,
  toolDeleteTask,
  toolCreateTaskSession,
  toolUpdateTaskSession,
  toolDeleteTaskSession,
  toolPlanDay,
  toolGetMetrics,
  toolCreateMetric,
  toolUpdateMetric,
  toolDeleteMetric,
  toolGetHabits,
  toolCreateHabit,
  toolUpdateHabit,
  toolDeleteHabit,
  toolGetGoals,
  toolCreateGoal,
  toolUpdateGoal,
  toolDeleteGoal,
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const callFunction = async (functionName: Tool, parameters: any): Promise<any> => {
  switch (functionName) {
    case Tool.GetCalendarEvents:
      return functionGetCalendarEvents(parameters);
    case Tool.GetEmails:
      return functionGetMails(parameters);
    case Tool.CreateCategory:
      return functionCreateCategory(parameters);
    case Tool.UpdateCategory:
      return functionUpdateCategory(parameters);
    case Tool.DeleteCategory:
      return functionDeleteCategory(parameters);
    case Tool.GetTasks:
      return functionGetTasks(parameters);
    case Tool.CreateTask:
      return functionCreateTask(parameters);
    case Tool.UpdateTask:
      return functionUpdateTask(parameters);
    case Tool.DeleteTask:
      return functionDeleteTask(parameters);
    case Tool.CreateTaskSession:
      return functionCreateTaskSession(parameters);
    case Tool.UpdateTaskSession:
      return functionUpdateTaskSession(parameters);
    case Tool.DeleteTaskSession:
      return functionDeleteTaskSession(parameters);
    case Tool.PlanDay:
      return functionPlanDay(parameters);
    case Tool.GetMetrics:
      return functionGetMetrics(parameters);
    case Tool.CreateMetric:
      return functionCreateMetric(parameters);
    case Tool.UpdateMetric:
      return functionUpdateMetric(parameters);
    case Tool.DeleteMetric:
      return functionDeleteMetric(parameters);
    case Tool.GetHabits:
      return functionGetHabits(parameters);
    case Tool.CreateHabit:
      return functionCreateHabit(parameters);
    case Tool.UpdateHabit:
      return functionUpdateHabit(parameters);
    case Tool.DeleteHabit:
      return functionDeleteHabit(parameters);
    case Tool.GetGoals:
      return functionGetGoals(parameters);
    case Tool.CreateGoal:
      return functionCreateGoal(parameters);
    case Tool.UpdateGoal:
      return functionUpdateGoal(parameters);
    case Tool.DeleteGoal:
      return functionDeleteGoal(parameters);
    default:
      throw new Error(`Function ${functionName} not found`);
  }
};
