// Original file: ../../proto/core/character_message.proto

import type {
  Category as _core_Category,
  Category__Output as _core_Category__Output,
} from "./Category";
import type { Metric as _core_Metric, Metric__Output as _core_Metric__Output } from "./Metric";
import type { Long } from "@grpc/proto-loader";

export interface Character {
  id?: string;
  createdAt?: number | string | Long;
  updatedAt?: number | string | Long;
  profileId?: string;
  name?: string;
  gender?: boolean;
  tags?: string[];
  categories?: _core_Category[];
  metrics?: _core_Metric[];
}

export interface Character__Output {
  id: string;
  createdAt: string;
  updatedAt: string;
  profileId: string;
  name: string;
  gender: boolean;
  tags: string[];
  categories: _core_Category__Output[];
  metrics: _core_Metric__Output[];
}
