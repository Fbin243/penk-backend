// Original file: ../../proto/core/character_message.proto

import type {
  CategoryStyleInput as _core_CategoryStyleInput,
  CategoryStyleInput__Output as _core_CategoryStyleInput__Output,
} from "./CategoryStyleInput";
import type {
  MetricInput as _core_MetricInput,
  MetricInput__Output as _core_MetricInput__Output,
} from "./MetricInput";

export interface CategoryInput {
  id?: string;
  name?: string;
  description?: string;
  style?: _core_CategoryStyleInput | null;
  metrics?: _core_MetricInput[];
  _id?: "id";
  _description?: "description";
}

export interface CategoryInput__Output {
  id?: string;
  name: string;
  description?: string;
  style: _core_CategoryStyleInput__Output | null;
  metrics: _core_MetricInput__Output[];
  _id: "id";
  _description: "description";
}
