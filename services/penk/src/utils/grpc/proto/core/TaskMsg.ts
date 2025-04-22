// Original file: ../../proto/core/task_message.proto

import type {
  Checkbox as _core_Checkbox,
  Checkbox__Output as _core_Checkbox__Output,
} from "./Checkbox";
import type { Long } from "@grpc/proto-loader";

export interface TaskMsg {
  id?: string;
  createdAt?: number | string | Long;
  updatedAt?: number | string | Long;
  characterId?: string;
  categoryId?: string;
  name?: string;
  priority?: number;
  completedTime?: number | string | Long;
  subtasks?: _core_Checkbox[];
  description?: string;
  deadline?: number | string | Long;
  _categoryId?: "categoryId";
  _completedTime?: "completedTime";
  _description?: "description";
  _deadline?: "deadline";
}

export interface TaskMsg__Output {
  id: string;
  createdAt: string;
  updatedAt: string;
  characterId: string;
  categoryId?: string;
  name: string;
  priority: number;
  completedTime?: string;
  subtasks: _core_Checkbox__Output[];
  description?: string;
  deadline?: string;
  _categoryId: "categoryId";
  _completedTime: "completedTime";
  _description: "description";
  _deadline: "deadline";
}
