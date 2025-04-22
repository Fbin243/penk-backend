// Original file: ../../proto/core/task_message.proto

import type {
  CheckboxInput as _core_CheckboxInput,
  CheckboxInput__Output as _core_CheckboxInput__Output,
} from "./CheckboxInput";
import type { Long } from "@grpc/proto-loader";

export interface TaskInput {
  id?: string;
  categoryId?: string;
  name?: string;
  priority?: number;
  completedTime?: number | string | Long;
  subtasks?: _core_CheckboxInput[];
  description?: string;
  deadline?: number | string | Long;
  _id?: "id";
  _categoryId?: "categoryId";
  _completedTime?: "completedTime";
  _description?: "description";
  _deadline?: "deadline";
}

export interface TaskInput__Output {
  id?: string;
  categoryId?: string;
  name: string;
  priority: number;
  completedTime?: string;
  subtasks: _core_CheckboxInput__Output[];
  description?: string;
  deadline?: string;
  _id: "id";
  _categoryId: "categoryId";
  _completedTime: "completedTime";
  _description: "description";
  _deadline: "deadline";
}
